package biz

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/4everland/ipfs-servers/app/blockstore/conf"
)

type BlockIndex interface {
	Put(ctx context.Context, cid string, v IndexValue) error
	Has(ctx context.Context, cid string) (bool, error)
	Get(ctx context.Context, cid string) (*IndexValue, error)

	Delete(ctx context.Context, cid string) error
	List(ctx context.Context) <-chan string
}

type IndexValue struct {
	offsetStart uint32
	size        uint32

	cumulativeSize uint64

	storeKey string
}

func NewIndexValue(offsetStart, size uint32, cumulativeSize uint64, storeKey string) IndexValue {
	return IndexValue{
		offsetStart:    offsetStart,
		size:           size,
		cumulativeSize: cumulativeSize,
		storeKey:       storeKey,
	}
}

func (iv *IndexValue) Encode() []byte {
	encoded := make([]byte, len(iv.storeKey)+16)
	binary.LittleEndian.PutUint32(encoded[:4], iv.offsetStart)
	binary.LittleEndian.PutUint32(encoded[4:8], iv.size)
	binary.LittleEndian.PutUint64(encoded[8:16], iv.cumulativeSize)

	copy(encoded[16:], iv.storeKey)
	return encoded
}

func (iv *IndexValue) Decode(encoded []byte) error {
	if len(encoded) < 16 {
		return errors.New("encode length to short")
	}
	iv.offsetStart = binary.LittleEndian.Uint32(encoded[0:4])
	iv.size = binary.LittleEndian.Uint32(encoded[4:8])
	iv.cumulativeSize = binary.LittleEndian.Uint64(encoded[8:16])

	iv.storeKey = string(encoded[16:])
	return nil
}

func (iv IndexValue) Size() (size uint32, cumulativeSize uint64) {
	return iv.size, iv.cumulativeSize
}

func (iv IndexValue) OffSet() (start, end uint32) {
	return iv.offsetStart, iv.offsetStart + iv.size
}

func NewIndexStore(data *conf.Data) BlockIndex {
	if data.GetDb().GetType() == conf.Data_TiKV {
		db, err := NewTiKv(data.GetDb().GetTikv().GetAddrs()...)
		if err != nil {
			panic(err)
		}

		return db
	}

	return NewLevelDb(data.GetDb().GetLeveldb().GetPath())
}
