package biz

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	ipld "github.com/ipfs/go-ipld-format"
	red "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hash/crc32"
	"time"
)

const (
	numberOfShards  = 64
	BloomFilterKey  = "blockstore:cid-bloom-filter"
	BloomErrorRatio = 0.01
	BloomSize       = 5000000
	IndexCacheKey   = "blockstore:index-cache"

	defaultPgIndexLRUSize = 1024
)

type PgIndexValue struct {
	Cid  string `gorm:"primarykey;column:id"`
	Size uint32
}

func (PgIndexValue) TableName() string {
	return "index"
}

type PgIndexStore struct {
	db *gorm.DB
	rd *red.Client

	cache *lru.Cache

	enableBloomQuery bool
}

func NewPg(db *gorm.DB, rd *red.Client, enableBloomQuery bool, lruSize int) (BlockIndex, error) {
	if lruSize <= 0 {
		lruSize = defaultPgIndexLRUSize
	}

	cache, err := lru.New(lruSize)
	if err != nil {
		return nil, err
	}

	for i := 0; i < numberOfShards; i++ {
		exists, err := rd.Exists(context.Background(), bloomFilterKey(i)).Result()
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			err = rd.BFReserve(context.Background(), bloomFilterKey(i), BloomErrorRatio, BloomSize).Err()
			if err != nil {
				return nil, err
			}
		}
	}

	return &PgIndexStore{
		db:    db,
		rd:    rd,
		cache: cache,

		enableBloomQuery: enableBloomQuery,
	}, nil
}

func (pg *PgIndexStore) Put(ctx context.Context, cid string, v IndexValue) error {
	err := pg.rd.BFAdd(ctx, bloomFilterKey(cid2TableIndex(cid)), cid).Err()
	if err != nil {
		return err
	}

	if err = pg.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&PgIndexValue{
		Cid:  cid,
		Size: v.size,
	}).Error; err != nil {
		return err
	}

	if err = pg.rd.Set(ctx, indexCacheKey(cid), v.size, time.Second*10).Err(); err != nil {
		return err
	}
	pg.setMemoryCache(cid, v.size)

	return nil
}

func (pg *PgIndexStore) Has(ctx context.Context, cid string) (bool, error) {
	if pg.enableBloomQuery {
		exists, err := pg.rd.BFExists(ctx, bloomFilterKey(cid2TableIndex(cid)), cid).Result()
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}

	if exists := pg.rd.Exists(ctx, indexCacheKey(cid)).Val(); exists == 1 {
		return true, nil
	}

	var v PgIndexValue
	if err := pg.db.WithContext(ctx).Select([]string{"id", "size"}).Take(&v, "id = ?", cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	pg.setMemoryCache(v.Cid, v.Size)
	return true, nil
}

func (pg *PgIndexStore) Delete(ctx context.Context, cid string) error {
	return pg.db.WithContext(ctx).Delete(&PgIndexValue{}, "id = ?", cid).Error
}

func (pg *PgIndexStore) Get(ctx context.Context, cid string) (*IndexValue, error) {
	if pg.enableBloomQuery {
		exists, err := pg.rd.BFExists(ctx, bloomFilterKey(cid2TableIndex(cid)), cid).Result()
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ipld.ErrNotFound{}
		}
	}

	size, err := pg.rd.Get(ctx, indexCacheKey(cid)).Uint64()
	if err == nil {
		v := newPgCachedIndexValue(cid, uint32(size))
		pg.cache.Add(cid, v)
		return &v, nil
	}

	if v, ok := pg.getMemoryCache(cid); ok {
		return v, nil
	}

	var v PgIndexValue
	if err := pg.db.WithContext(ctx).Select([]string{"id", "size"}).Take(&v, "id = ?", cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ipld.ErrNotFound{}
		}
		return nil, err
	}

	value := newPgCachedIndexValue(v.Cid, v.Size)
	pg.cache.Add(cid, value)

	return &value, nil
}

func (pg *PgIndexStore) List(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var (
			startKey  string
			tableName string
			data      []PgIndexValue
		)

		for i := 0; i < numberOfShards; i++ {
			tableName = fmt.Sprintf("%s_%02d", PgIndexValue{}.TableName(), i)
			startKey = ""
			for {
				if err := pg.db.WithContext(ctx).Table(tableName).Where("id > ?", startKey).Select("id").
					Limit(1000).Order("id ASC").Find(&data).Error; err != nil || len(data) == 0 {
					break
				}

				for _, row := range data {
					ch <- row.Cid
				}

				startKey = data[len(data)-1].Cid
			}
		}

	}()
	return ch
}

func bloomFilterKey(i int) string {
	return fmt.Sprintf("%s:%d", BloomFilterKey, i)
}

func cid2TableIndex(cid string) int {
	return int(crc32.ChecksumIEEE([]byte(cid))) % 64
}

func indexCacheKey(cid string) string {
	return fmt.Sprintf("%s:%s", IndexCacheKey, cid)
}

func newPgCachedIndexValue(cid string, size uint32) IndexValue {
	return IndexValue{
		size:     size,
		storeKey: cid,
	}
}

func (pg *PgIndexStore) getMemoryCache(cid string) (*IndexValue, bool) {
	if pg.cache == nil {
		return nil, false
	}

	v, ok := pg.cache.Get(cid)
	if !ok {
		return nil, false
	}

	value := v.(IndexValue)
	return &value, true
}

func (pg *PgIndexStore) setMemoryCache(cid string, size uint32) {
	if pg.cache == nil {
		return
	}

	pg.cache.Add(cid, newPgCachedIndexValue(cid, size))
}
