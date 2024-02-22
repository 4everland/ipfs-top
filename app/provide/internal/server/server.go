package server

import (
	"context"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/4everland/ipfs-top/third_party/dag"
	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/ipfs/boxo/blockstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewBlockStore,
	NewNodes,
	NewConsumer,
	NewEventServer,
	NewReproviderServer,
)

func NewBlockStore(conf *conf.Data) blockstore.Blockstore {
	bs, err := dag.NewBlockStore(conf.BlockstoreUri, "")
	if err != nil {
		panic(err)
	}

	return bs
}

func NewNodes(config *conf.Data) []routing.RoutingClient {
	nodes := make([]routing.RoutingClient, 0)
	for _, endpoint := range config.Target {
		if endpoint.ServerType != conf.ServerType_GRPC {
			continue
		}
		tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc2.Dial(
			context.Background(),
			grpc2.WithEndpoint(endpoint.Endpoint),
			grpc2.WithMiddleware(
				tracing.Client(),
				recovery.Recovery(),
			),
			grpc2.WithTimeout(time.Minute),
			grpc2.WithOptions(tlsOption),
		)
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, routing.NewRoutingClient(conn))
	}

	return nodes
}

func NewConsumer(conf *conf.Data) sarama.ConsumerGroup {
	if !conf.EnableProvider {
		return nil
	}
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Fetch.Min = 8192
	config.Consumer.MaxWaitTime = time.Millisecond * 500
	config.Version = sarama.V2_6_2_0
	//config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	//config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}

	//consumer, err := sarama.NewConsumer(conf.Kafka.Addr, config)
	client, err := sarama.NewConsumerGroup(conf.Kafka.Addr, conf.Kafka.Group, config)

	if err != nil {
		panic(err)
	}

	return client
}
