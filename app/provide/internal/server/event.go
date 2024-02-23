package server

import (
	"context"
	"errors"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/4everland/ipfs-top/app/provide/internal/data"
	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/go-cid"
	"runtime"
	"strings"
	"time"
)

type EventServer struct {
	conf     *conf.Kafka
	nodes    []*data.NamedRoutingClient
	consumer sarama.ConsumerGroup
	//metrics  *prom.ProvideMetrics
}

type Consumer struct {
	ready  chan bool
	logger *log.Helper

	event chan cid.Cid
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				consumer.logger.Warnf("message channel was closed")
				return nil
			}
			//consumer.logger.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
			if message == nil || message.Value == nil {
				continue
			}

			c, err := cid.Cast(message.Value)
			//c, err := cid.Decode(string(message.Value))

			if err != nil {
				consumer.logger.Warnf("Message claimed: value = %s, timestamp = %v, topic = %s, decode err: %s", string(message.Value), message.Timestamp, message.Topic, err)
				continue
			}

			if consumer.event != nil {
				consumer.event <- c
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func (server *EventServer) Start(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 64<<10) //nolint:gomnd
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			log.NewHelper(log.DefaultLogger).WithContext(ctx).Errorf("event server panic %v:", err)
		}
	}()
	logger := log.NewHelper(log.DefaultLogger)
	consumer := Consumer{
		ready:  make(chan bool),
		logger: logger,
		event:  make(chan cid.Cid),
	}
	product := make(chan cid.Cid, len(server.nodes))
	for _, node := range server.nodes {
		go func(ctx context.Context, node *data.NamedRoutingClient) {
			errCount := 0
			for c := range product {
				_, err := node.Client.Provide(ctx, &routing.ProvideReq{
					Cid:     &routing.Cid{Str: c.Bytes()},
					Provide: true,
				})
				if err != nil {
					errCount++
					log.NewHelper(log.DefaultLogger).WithContext(ctx).Errorf("provide %s error: %v", c.String(), err)
					product <- c
				} else {
					//server.metrics.Provide(node.Name)
					continue
				}
				if errCount >= 10 {
					log.NewHelper(log.DefaultLogger).Warnf("provide %s error count: %d", c.String(), errCount)
					time.Sleep(time.Second * 20)
				}
				errCount = 0
			}
		}(ctx, node)
	}

	go func() {
		for c := range consumer.event {
			product <- c
			//var wg sync.WaitGroup
			//for _, node := range server.nodes {
			//	wg.Add(1)
			//	go func(r routing.RoutingClient) {
			//		defer wg.Done()
			//		if _, err := r.Provide(ctx, &routing.ProvideReq{
			//			Cid:     &routing.Cid{Str: c.Bytes()},
			//			Provide: true,
			//		}); err != nil {
			//			log.NewHelper(log.DefaultLogger).WithContext(ctx).Errorf("provide %s error: %v:", c.String(), err)
			//		}
			//	}(node)
			//}
			//wg.Wait()
		}
	}()

	go func() {
		for {
			if err := server.consumer.Consume(ctx, strings.Split(server.conf.Topic, ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				logger.Warnf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	return nil
}

func (server *EventServer) Stop(context.Context) error {
	return server.consumer.Close()
}

func NewEventServer(
	conf *conf.Data,
	consumer sarama.ConsumerGroup,
	nodes []*data.NamedRoutingClient,
) *EventServer {
	if !conf.EnableProvider {
		return nil
	}
	return &EventServer{
		conf:     conf.Kafka,
		consumer: consumer,
		nodes:    nodes,
	}
}
