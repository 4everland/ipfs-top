package service

import (
	"context"
	"github.com/4everland/ipfs-servers/app/provide/conf"
	"github.com/ipfs/go-datastore"
	provider "github.com/ipfs/go-ipfs-provider"
	q "github.com/ipfs/go-ipfs-provider/queue"
	"github.com/ipfs/go-ipfs-provider/simple"
	"github.com/libp2p/go-libp2p/core/routing"
	"time"
)

// ProviderQueue creates new datastore backed provider queue
func ProviderQueue(ds datastore.Datastore) (*q.Queue, error) {
	return q.NewQueue(context.Background(), "provider-v1", ds)
}

// SimpleProvider creates new record provider
func SimpleProvider(queue *q.Queue, rt routing.ContentRouting) provider.Provider {
	return simple.NewProvider(context.Background(), queue, rt)
}

func SimpleReprovider(data *conf.Data, rt routing.ContentRouting, keyProvider simple.KeyChanFunc) (provider.Reprovider, error) {
	reproviderInterval := time.Duration(data.ReproviderInterval) * time.Second
	return simple.NewReprovider(context.Background(), reproviderInterval, rt, keyProvider), nil
}
