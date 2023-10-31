package tasks

import (
	"context"
	"github.com/4everland/ipfs-top/third_party/arweave"
	"github.com/4everland/ipfs-top/third_party/s3client"
	arseedSchema "github.com/everFinance/arseeding/schema"
	"github.com/everFinance/arseeding/sdk/schema"
	"github.com/everFinance/goar/types"
	"io"
)

type (
	Bridge interface {
		SendtoAR(ctx context.Context, uuid, key string, tags []types.Tag) (*arseedSchema.RespOrder, error)
	}
	S3Bridge struct {
		Ar      *arweave.SDK
		storage *s3client.S3Storage
	}
)

func NewS3Bridge(arseeding, mnemonic string, config s3client.S3Config) (Bridge, error) {
	storage := s3client.NewS3Client(config, s3client.DisableTransformKeyOptions)
	sdk, err := arweave.NewSDK(arseeding, mnemonic)
	if err != nil {
		return nil, err
	}
	return &S3Bridge{
		Ar:      sdk,
		storage: storage,
	}, nil
}

func (b *S3Bridge) SendtoAR(ctx context.Context, uuid, obj string, tags []types.Tag) (*arseedSchema.RespOrder, error) {
	body, err := b.storage.Get(ctx, obj)
	if err != nil {
		return nil, err
	}

	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	transaction, _, err := b.Ar.SendData(uuid, data, &schema.OptionItem{
		Tags: tags,
	})
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
