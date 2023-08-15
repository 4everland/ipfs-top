package biz

import (
	"github.com/4everland/ipfs-servers/app/blockstore/internal/conf"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewBackendStorage, NewIndexStore)

func NewIndexStore(data *conf.Data) (BlockIndex, error) {
	switch data.GetDb().GetType() {
	case conf.Data_TiKV:
		return NewTiKv(data.GetDb().GetTikv().GetAddrs()...)
	case conf.Data_PG:
		d, err := gorm.Open(postgres.Open(data.GetDb().GetPg().GetDsn()))
		if err != nil {
			return nil, err
		}
		return NewPg(d)
	default:
		return NewLevelDb(data.GetDb().GetLeveldb().GetPath())
	}
}
