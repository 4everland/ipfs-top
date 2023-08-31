package biz

import (
	"context"
	"fmt"
	ipld "github.com/ipfs/go-ipld-format"
	red "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hash/crc32"
)

const (
	numberOfShards  = 64
	BloomFilterKey  = "blockstore:cid-bloom-filter"
	BloomErrorRatio = 0.01
	BloomSize       = 5000000
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

	enableBloomQuery bool
}

func NewPg(db *gorm.DB, rd *red.Client, enableBloomQuery bool) (BlockIndex, error) {

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
		db: db,
		rd: rd,

		enableBloomQuery: enableBloomQuery,
	}, nil
}

func (pg *PgIndexStore) Put(ctx context.Context, cid string, v IndexValue) error {
	err := pg.rd.BFAdd(ctx, bloomFilterKey(cid2TableIndex(cid)), cid).Err()
	if err != nil {
		return err
	}
	return pg.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&PgIndexValue{
		Cid:  cid,
		Size: v.size,
	}).Error
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
	if err := pg.db.WithContext(ctx).Select("id").Take(&PgIndexValue{}, "id = ?", cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (pg *PgIndexStore) Delete(ctx context.Context, cid string) error {
	return pg.db.WithContext(ctx).Delete(&PgIndexValue{}, "id = ?", cid).Error
}

func (pg *PgIndexStore) Get(ctx context.Context, cid string) (*IndexValue, error) {
	var v PgIndexValue
	if err := pg.db.WithContext(ctx).Select([]string{"id", "size"}).Take(&v, "id = ?", cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ipld.ErrNotFound{}
		}
		return nil, err
	}

	return &IndexValue{
		size:     v.Size,
		storeKey: v.Cid,
	}, nil
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
				if err := pg.db.WithContext(ctx).Table(tableName).Where("id > ? COLLATE \"C\"", startKey).Select("id").
					Limit(1000).Order("id COLLATE \"C\" ASC").Find(&data).Error; err != nil || len(data) == 0 {
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
