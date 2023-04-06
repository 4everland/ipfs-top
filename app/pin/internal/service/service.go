package service

import (
	"context"
	"github.com/4everland/ipfs-servers/app/pin/internal/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/4everland/ipfs-servers/third_party/datastore"
	"github.com/google/wire"
	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/tikv/client-go/v2/rawkv"
)

var ProviderSet = wire.NewSet(
	NewBlockStore,
	NewExchange,
	NewDatastore,
	NewPinService,
)

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}

func NewExchange(config *conf.Data) exchange.Interface {
	s, err := dag.NewGrpcRouting(config.ExchangeEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func NewDatastore(config *conf.Data) ds.Datastore {
	if config.GetDatastore().GetType() == conf.Data_TiKV {
		client, err := NewTiKv(config.GetDatastore().GetTikv().GetAddrs()...)
		if err != nil {
			panic(err)
		}

		return datastore.NewRawKVDatastore(client)
	}

	d, err := leveldb.NewDatastore(config.GetDatastore().GetLeveldb().GetPath(), &leveldb.Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil {
		panic(err)
	}

	return d
}

func NewTiKv(pdAddr ...string) (*rawkv.Client, error) {
	client, err := rawkv.NewClientWithOpts(context.Background(), pdAddr)
	if err != nil {
		return nil, err
	}
	return client, nil
}
