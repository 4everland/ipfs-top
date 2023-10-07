package biz

import (
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/provider"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/routing"
)

func ProviderSystem(ds datastore.Batching, bs blockstore.Blockstore) func(rt routing.ContentRouting) (provider.System, error) {
	return func(rt routing.ContentRouting) (provider.System, error) {
		//return provider.New(ds, provider.Online(rt), provider.KeyProvider(bs.AllKeysChan))
		return provider.New(ds, provider.Online(rt))
	}
}
