package datastore

import (
	"context"
	"errors"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
	"github.com/tikv/client-go/v2/rawkv"
	"time"
)

var _ ds.Datastore = (*RawKVDatastore)(nil)
var _ ds.TTLDatastore = (*RawKVDatastore)(nil)

//var _ ds.Batching = (*RawKVDatastore)(nil)

type RawKVDatastore struct {
	client *rawkv.Client
}

func NewRawKVDatastore(client *rawkv.Client) *RawKVDatastore {
	return &RawKVDatastore{
		client: client,
	}
}

func (raw *RawKVDatastore) Put(ctx context.Context, key ds.Key, value []byte) error {
	return raw.client.Put(ctx, key.Bytes(), value)
}

func (raw *RawKVDatastore) Get(ctx context.Context, key ds.Key) ([]byte, error) {
	b, err := raw.client.Get(ctx, key.Bytes())
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ds.ErrNotFound
	}

	return b, nil
}

func (raw *RawKVDatastore) Has(ctx context.Context, key ds.Key) (bool, error) {
	b, err := raw.client.Get(ctx, key.Bytes())
	if err != nil {
		return false, err
	}
	if b == nil {
		return false, ds.ErrNotFound
	}

	return true, nil
}

func (raw *RawKVDatastore) GetSize(ctx context.Context, key ds.Key) (int, error) {
	b, err := raw.client.Get(ctx, key.Bytes())
	if err != nil {
		return 0, err
	}
	if b == nil {
		return 0, ds.ErrNotFound
	}

	return len(b), nil
}

func (raw *RawKVDatastore) Delete(ctx context.Context, key ds.Key) error {
	return raw.client.Delete(ctx, key.Bytes())
}

func (raw *RawKVDatastore) Sync(context.Context, ds.Key) error {
	return nil
}

func (raw *RawKVDatastore) Close() error {
	return raw.client.Close()
}

func (raw *RawKVDatastore) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
	if q.Orders != nil || q.Filters != nil {
		return nil, fmt.Errorf("rawkv: orders or filters are not supported")
	}

	var opts []rawkv.RawOption
	if q.KeysOnly && !q.ReturnsSizes {
		opts = append(opts, rawkv.ScanKeyOnly())
	}

	keys, values, err := raw.client.Scan(ctx, []byte(q.Prefix), nil, q.Limit, opts...)
	if err != nil {
		return nil, err
	}

	var (
		index     int
		keysCount = len(keys)
	)
	return dsq.ResultsFromIterator(q, dsq.Iterator{
		Close: func() error {
			return nil
		},
		Next: func() (dsq.Result, bool) {
			if index >= keysCount {
				return dsq.Result{}, false
			}
			ent := dsq.Entry{
				Key: string(keys[index]),
			}
			if !q.KeysOnly {
				ent.Value = values[index]
			}
			if q.ReturnsSizes {
				ent.Size = len(ent.Value)
			}
			if q.ReturnExpirations {
				exp, err := raw.GetExpiration(ctx, ds.NewKey(ent.Key))
				if err != nil {
					return dsq.Result{Error: err}, true
				}
				ent.Expiration = exp
			}
			return dsq.Result{Entry: ent}, true
		},
	}), nil
}

func (raw *RawKVDatastore) GetExpiration(ctx context.Context, key ds.Key) (time.Time, error) {
	ttl, err := raw.client.GetKeyTTL(ctx, key.Bytes())
	if err != nil {
		if errors.Is(err, ds.ErrNotFound) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return time.Now().Add(time.Duration(*ttl)), nil
}

func (raw *RawKVDatastore) PutWithTTL(ctx context.Context, key ds.Key, value []byte, ttl time.Duration) error {
	return raw.client.PutWithTTL(ctx, key.Bytes(), value, uint64(ttl))
}

func (raw *RawKVDatastore) SetTTL(ctx context.Context, key ds.Key, ttl time.Duration) error {
	value, err := raw.Get(ctx, key)
	if err != nil {
		return err
	}
	return raw.PutWithTTL(ctx, key, value, ttl)
}
