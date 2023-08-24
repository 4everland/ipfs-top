package biz

import (
	"bytes"
	"context"
	"errors"
	"github.com/4everland/diskv/v3"
	"github.com/4everland/ipfs-servers/app/gateway/internal/conf"
	"github.com/4everland/ipfs-servers/third_party/s3client"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"io"
	"os"
)

type readOnlyS3blockStore struct {
	s3Client *s3client.S3Storage
	c        *diskv.Diskv
	log      *log.Helper
}

func NewS3readOnlyS3blockStore(data *conf.Data, logger log.Logger) blockstore.Blockstore {
	if data.GetRo().GetCache() == nil {
		return &readOnlyS3blockStore{s3Client: s3client.NewS3Client(data.GetRo().GetStorage())}
	}
	return &readOnlyS3blockStore{
		s3Client: s3client.NewS3Client(data.GetRo().GetStorage()),
		c: diskv.New(diskv.Options{
			BasePath:     data.GetRo().GetCache().GetBasePath(),
			LruSize:      int(data.GetRo().GetCache().GetLruSize()),
			LruIndexPath: data.GetRo().GetCache().GetIndexPath(),
			Transform:    diskv.BlockTransform(3, 3, true),
		}),
		log: log.NewHelper(logger),
	}
}

func (bs *readOnlyS3blockStore) Get(ctx context.Context, c cid.Cid) (block blocks.Block, err error) {
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	key := c.String()
	var r io.ReadCloser
	if r, err = bs.c.ReadStream(key, true); err == nil {
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

	return blocks.NewBlockWithCid(bf.Bytes(), c)
}

func (bs *readOnlyS3blockStore) GetSize(ctx context.Context, c cid.Cid) (int, error) {
	key := c.String()
	if size := bs.c.Size(key); size >= 0 {
		return int(size), nil
	}

	return bs.s3Client.GetSize(ctx, key)
}

func (bs *readOnlyS3blockStore) DeleteBlock(_ context.Context, c cid.Cid) error {
	if err := bs.c.Erase(c.String()); err != nil {
		bs.log.Error("delete disk cache error:", err)
		return err
	}

	return nil
}

func (bs *readOnlyS3blockStore) Has(ctx context.Context, c cid.Cid) (bool, error) {
	return true, nil
}

func (bs *readOnlyS3blockStore) Put(context.Context, blocks.Block) error {
	return nil
}

func (bs *readOnlyS3blockStore) PutMany(context.Context, []blocks.Block) error {
	return nil
}

func (bs *readOnlyS3blockStore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, nil
}

func (bs *readOnlyS3blockStore) HashOnRead(enabled bool) {}
