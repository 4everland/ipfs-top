package biz

import (
	"bytes"
	"context"
	"errors"
	"github.com/4everland/ipfs-top/app/blockstore/internal/conf"
	"github.com/4everland/ipfs-top/third_party/prom"
	"github.com/4everland/ipfs-top/third_party/s3client"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"os"
)

type BlockStore interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key string, cid string, data []byte) error
	Delete(ctx context.Context, key string) error
	GetSize(ctx context.Context, key string) (int, error)
}

type blockStore struct {
	s3Client *s3client.S3Storage
	c        *diskv.Diskv
	log      *log.Helper
	metrics  *prom.BlockStoreMetrics
}

func NewBackendStorage(data *conf.Data, logger log.Logger) BlockStore {
	if data.Cache == nil {
		return s3client.NewS3Client(data.Storage)
	}
	return &blockStore{
		s3Client: s3client.NewS3Client(data.Storage),
		c: diskv.New(diskv.Options{
			BasePath:     data.Cache.BasePath,
			LruSize:      int(data.Cache.LruSize),
			LruIndexPath: data.Cache.IndexPath,
			Transform:    diskv.BlockTransform(3, 3, true),
		}),
		log:     log.NewHelper(logger),
		metrics: prom.NewBlockStoreMetrics(),
	}
}

func (bs *blockStore) Get(ctx context.Context, key string) (r io.ReadCloser, err error) {
	if r, err = bs.c.ReadStream(key, true); err == nil {
		bs.metrics.IncrCacheHits()
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		bs.log.Error("get disk cache error:", err)
	}

	if r, err = bs.s3Client.Get(ctx, key); err != nil {
		return
	}

	defer r.Close()
	var bf bytes.Buffer
	if err = bs.c.WriteStream(key, io.TeeReader(r, &bf), false); err != nil {
		bs.log.Error("write disk cache error:", err)
	}

	bs.metrics.IncrBlockStoreHits()

	return io.NopCloser(&bf), nil
}

func (bs *blockStore) Put(ctx context.Context, key string, cid string, data []byte) error {
	return bs.s3Client.Put(ctx, key, cid, data)
}

func (bs *blockStore) Delete(ctx context.Context, key string) error {
	if err := bs.c.Erase(key); err != nil {
		bs.log.Error("delete disk cache error:", err)
	}
	return bs.s3Client.Delete(ctx, key)
}

func (bs *blockStore) GetSize(ctx context.Context, key string) (int, error) {
	if size := bs.c.Size(key); size >= 0 {
		return int(size), nil
	}

	return bs.s3Client.GetSize(ctx, key)
}
