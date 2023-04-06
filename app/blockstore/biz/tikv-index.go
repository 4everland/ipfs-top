package biz

import (
	"context"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/tikv/client-go/v2/rawkv"
)

type TiKvIndexStore struct {
	client *rawkv.Client
}

func NewTiKv(pdAddr ...string) (BlockIndex, error) {
	client, err := rawkv.NewClientWithOpts(context.Background(), pdAddr)
	if err != nil {
		return nil, err
	}
	return &TiKvIndexStore{
		client: client,
	}, nil
}

func (tis *TiKvIndexStore) Put(ctx context.Context, cid string, v IndexValue) error {
	return tis.client.Put(ctx, []byte(cid), v.Encode())
}

func (tis *TiKvIndexStore) Has(ctx context.Context, cid string) (bool, error) {
	ttl, err := tis.client.GetKeyTTL(ctx, []byte(cid))
	return ttl != nil, err
}

func (tis *TiKvIndexStore) Delete(ctx context.Context, cid string) error {
	return tis.client.Delete(ctx, []byte(cid))
}

func (tis *TiKvIndexStore) Get(ctx context.Context, cid string) (*IndexValue, error) {
	v, err := tis.client.Get(ctx, []byte(cid))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ipld.ErrNotFound{}
	}

	var iv = &IndexValue{}
	err = iv.Decode(v)
	return iv, nil
}

func (tis *TiKvIndexStore) List(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var startKey = []byte("")
		for {
			keys, _, err := tis.client.Scan(ctx, startKey, []byte(""), 1000, rawkv.ScanKeyOnly())
			if err != nil || len(keys) == 0 {
				return
			}
			for _, k := range keys {
				ch <- string(k)
			}
			endKey := keys[len(keys)-1]
			startKey = make([]byte, len(endKey)+1)
			copy(startKey, endKey)

			startKey[len(endKey)] = '\000'
		}
	}()
	return ch
}
