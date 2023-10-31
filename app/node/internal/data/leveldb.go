package data

import (
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
)

func NewLevelDbDatastore(server *conf.Server) (ds.Batching, error) {
	return leveldb.NewDatastore(server.Node.LeveldbPath, &leveldb.Options{
		Filter: filter.NewBloomFilter(10),
	})
}
