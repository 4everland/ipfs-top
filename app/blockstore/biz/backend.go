package biz

import (
	"context"
	"github.com/4everland/ipfs-servers/app/blockstore/conf"
	"github.com/4everland/ipfs-servers/third_party/s3client"
	"io"
)

type BlockStore interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key string, cid string, data []byte) error
	Delete(ctx context.Context, key string) error
	GetSize(ctx context.Context, key string) (int, error)
}

func NewBackendStorage(data *conf.Data) BlockStore {
	return s3client.NewS3Client(data.Storage)
}
