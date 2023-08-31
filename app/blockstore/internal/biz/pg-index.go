package biz

import (
	"context"
	"fmt"
	ipld "github.com/ipfs/go-ipld-format"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const numberOfShards = 64

type PgIndexValue struct {
	Cid  string `gorm:"primarykey;column:id"`
	Size uint32
}

func (PgIndexValue) TableName() string {
	return "index"
}

type PgIndexStore struct {
	db *gorm.DB
}

func NewPg(db *gorm.DB) (BlockIndex, error) {
	return &PgIndexStore{
		db: db,
	}, nil
}

func (pg *PgIndexStore) Put(ctx context.Context, cid string, v IndexValue) error {
	return pg.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&PgIndexValue{
		Cid:  cid,
		Size: v.size,
	}).Error
}

func (pg *PgIndexStore) Has(ctx context.Context, cid string) (bool, error) {
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
