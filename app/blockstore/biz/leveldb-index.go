package biz

import (
	"context"
	"errors"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbIndexStore struct {
	db *leveldb.DB
}

func wrapperError(err error) error {
	if errors.Is(err, leveldb.ErrNotFound) {
		err = ipld.ErrNotFound{}
	}
	return err
}

func NewLevelDb(path string) BlockIndex {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDbIndexStore{
		db: db,
	}
}

func (lis *LevelDbIndexStore) Put(ctx context.Context, cid string, v IndexValue) error {
	return lis.db.Put([]byte(cid), v.Encode(), nil)
}

func (lis *LevelDbIndexStore) Has(ctx context.Context, cid string) (bool, error) {
	exists, err := lis.db.Has([]byte(cid), nil)
	return exists, wrapperError(err)
}

func (lis *LevelDbIndexStore) Delete(ctx context.Context, cid string) error {
	err := lis.db.Delete([]byte(cid), nil)
	return wrapperError(err)
}

func (lis *LevelDbIndexStore) Get(ctx context.Context, cid string) (*IndexValue, error) {
	v, err := lis.db.Get([]byte(cid), nil)

	if err != nil {
		return nil, wrapperError(err)
	}
	var iv = &IndexValue{}
	err = iv.Decode(v)
	return iv, wrapperError(err)

}

func (lis *LevelDbIndexStore) List(ctx context.Context) <-chan string {
	iter := lis.db.NewIterator(nil, nil)
	ch := make(chan string)
	go func() {
		defer iter.Release()
		defer close(ch)
		for iter.Next() {
			ch <- string(iter.Key())
		}
	}()
	return ch

}
