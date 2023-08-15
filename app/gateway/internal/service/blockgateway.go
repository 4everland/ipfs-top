package service

import (
	"github.com/4everland/ipfs-servers/app/gateway/internal/biz"
	"github.com/4everland/ipfs-servers/app/gateway/internal/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/gateway"
)

func NewBlocksGateway(blockStore blockstore.Blockstore) (*gateway.BlocksBackend, error) {
	return gateway.NewBlocksBackend(blockservice.New(blockStore, nil))
}

func NewBlockStore(config *conf.Data, logger log.Logger) blockstore.Blockstore {
	if config.GetRo() != nil {
		return biz.NewS3readOnlyS3blockStore(config, logger)
	}
	s, err := dag.NewBlockStore(config.GetRw().GetUri(), config.GetRw().GetUri())
	if err != nil {
		panic(err)
	}
	return s
}
