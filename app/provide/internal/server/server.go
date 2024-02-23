package server

import (
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/4everland/ipfs-top/third_party/dag"
	"github.com/IBM/sarama"

	"github.com/google/wire"
	"github.com/ipfs/boxo/blockstore"

	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewBlockStore,
	NewConsumer,
	NewEventServer,
	NewReproviderServer,
	NewMetricsServer,
)

func NewBlockStore(conf *conf.Data) blockstore.Blockstore {
	bs, err := dag.NewBlockStore(conf.BlockstoreUri, "")
	if err != nil {
		panic(err)
	}

	return bs
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
