package service

import (
	"github.com/4everland/ipfs-top/app/gateway/internal/biz"
	"github.com/4everland/ipfs-top/app/gateway/internal/conf"
	"github.com/4everland/ipfs-top/third_party/dag"
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
		if config.GetRo().GetCache() == nil {
			panic("cache not config")
		}
		return biz.NewS3readOnlyS3blockStore(config, logger)
	}
	s, err := dag.NewBlockStore(config.GetRw().GetUri(), config.GetRw().GetCert())
	if err != nil {
		panic(err)
	}
	return s
}
