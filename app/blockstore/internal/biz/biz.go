package biz

import (
	"github.com/4everland/ipfs-servers/app/blockstore/internal/conf"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/sharding"
	"log"
	"os"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewBackendStorage, NewIndexStore)

func NewIndexStore(data *conf.Data) (BlockIndex, error) {
	switch data.GetDb().GetType() {
	case conf.Data_TiKV:
		return NewTiKv(data.GetDb().GetTikv().GetAddrs()...)
	case conf.Data_PG:
		dsn := data.GetDb().GetPg().GetSourcesDsn()
		sourceDials := make([]gorm.Dialector, len(dsn))
		for i, dsn := range data.GetDb().GetPg().GetSourcesDsn() {
			sourceDials[i] = postgres.Open(dsn)
		}

		d, err := gorm.Open(sourceDials[0], &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
			}),
			SkipDefaultTransaction: true,
		})
		if err != nil {
			return nil, err
		}

		middleware := sharding.Register(sharding.Config{
			ShardingKey:         "id",
			NumberOfShards:      numberOfShards,
			PrimaryKeyGenerator: sharding.PKSnowflake,
		}, PgIndexValue{})
		if err := d.Use(middleware); err != nil {
			return nil, err
		}

		dsn = data.GetDb().GetPg().GetReplicasDsn()
		if len(dsn) > 0 {
			replicasDials := make([]gorm.Dialector, len(dsn))
			for i, s := range dsn {
				replicasDials[i] = postgres.Open(s)
			}
			if err = d.Use(dbresolver.Register(dbresolver.Config{
				Sources:           sourceDials,
				Replicas:          replicasDials,
				TraceResolverMode: true,
			})); err != nil {
				return nil, err
			}
		}

		return NewPg(d)
	default:
		return NewLevelDb(data.GetDb().GetLeveldb().GetPath())
	}
}
