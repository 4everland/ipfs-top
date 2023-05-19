package service

import (
	"github.com/4everland/ipfs-servers/app/gateway/internal/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/gateway"
)

func NewBlocksGateway(blockStore blockstore.Blockstore) (*gateway.BlocksGateway, error) {
	return gateway.NewBlocksGateway(blockservice.New(blockStore, nil))
}

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri, config.BlockstoreCert)
	if err != nil {
		panic(err)
	}
	return s
}
