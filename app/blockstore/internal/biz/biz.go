package biz

import (
	"github.com/4everland/ipfs-servers/app/blockstore/internal/conf"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
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

func NewRedisClient(data *conf.Data) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         data.Redis.Addr,
		Password:     data.Redis.Password,
		DB:           0,
		PoolSize:     32,
		ReadTimeout:  data.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: data.Redis.WriteTimeout.AsDuration(),
	})
	return client, nil
}

func NewIndexStore(data *conf.Data) (BlockIndex, error) {
	rd, err := NewRedisClient(data)
	if err != nil {
		return nil, err
	}
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
				SlowThreshold:             300 * time.Millisecond,
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

		return NewPg(d, rd, data.GetDb().GetPg().GetEnableBloom())
	default:
		return NewLevelDb(data.GetDb().GetLeveldb().GetPath())
	}
}
