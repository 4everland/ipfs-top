package data

import (
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/third_party/dag"
	"github.com/ipfs/boxo/blockstore"
)

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri, config.BlockstoreCert)
	if err != nil {
		panic(err)
	}
	return s
}
