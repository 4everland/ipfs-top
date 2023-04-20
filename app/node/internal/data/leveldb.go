package data

import (
	"github.com/4everland/ipfs-servers/app/node/internal/conf"
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/provider/simple"
	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
)

func NewLevelDbDatastore(server *conf.Server) (ds.Batching, error) {
	return leveldb.NewDatastore(server.Node.LeveldbPath, &leveldb.Options{
		Filter: filter.NewBloomFilter(10),
	})
}

func NewKeyChanFunc(data blockstore.Blockstore) simple.KeyChanFunc {
	return simple.NewBlockstoreProvider(data)
}
