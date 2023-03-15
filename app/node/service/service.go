package service

import (
	"github.com/4everland/ipfs-servers/app/node/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/google/wire"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}

func NewNodeServices(bts *BitSwapService) []NodeService {
	return []NodeService{
		bts,
		NewContentRoutingService(),
	}
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewBlockStore, NewContentRoutingService, NewBitSwapService, NewNodeServices)
