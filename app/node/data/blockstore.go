package data

import (
	"github.com/4everland/ipfs-servers/app/node/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}
