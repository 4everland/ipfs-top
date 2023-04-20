package data

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
	"strings"
)

const (
	defaultMaxSize        = 524280
	defaultMaxPeerPerItem = 1024
)

// LruMapDatastore uses a standard Go map for internal storage.
type LruMapDatastore struct {
	values   *lru.Cache
	maxSize  int
	maxPeers int
}

type Option func(*LruMapDatastore) error

func (d *LruMapDatastore) applyOptions(opts ...Option) error {
	for i, opt := range opts {
		if err := opt(d); err != nil {
			return fmt.Errorf("LruMapDatastore option %d failed: %s", i, err)
		}
	}
	return nil
}

func LruMaxSize(maxSize int) Option {
	return func(d *LruMapDatastore) error {
		d.maxSize = maxSize
		return nil
	}
}

func LruMaxPeer(maxPeer int) Option {
	return func(d *LruMapDatastore) error {
		d.maxPeers = maxPeer
		return nil
	}
}

var _ datastore.Datastore = (*LruMapDatastore)(nil)
var _ datastore.Batching = (*LruMapDatastore)(nil)

func NewLruMapDatastore() (d *LruMapDatastore, err error) {
	cache, err := lru.New(defaultMaxSize)
	if err != nil {
		return nil, err
	}
	return &LruMapDatastore{
		values:   cache,
		maxSize:  defaultMaxSize,
		maxPeers: defaultMaxPeerPerItem,
	}, nil
}

func NewLruMapDatastoreWithOpt(opts ...Option) (lmd *LruMapDatastore, err error) {
	lmd = new(LruMapDatastore)
	err = lmd.applyOptions(opts...)
	if err != nil {
		return nil, err
	}
	if lmd.maxSize == 0 {
		lmd.maxSize = defaultMaxSize
	}
	if lmd.maxPeers == 0 {
		lmd.maxPeers = defaultMaxPeerPerItem
	}
	lmd.values, err = lru.New(lmd.maxSize)
	if err != nil {
		return nil, err
	}
	return
}

// Put implements Datastore.Put
func (d *LruMapDatastore) Put(ctx context.Context, key datastore.Key, value []byte) (err error) {
	prefix, p := keySplit(key.String())
	set, ok := d.values.Get(prefix)
	if ok {
		peerMap := set.(*lru.Cache)
		peerMap.Add(p, value)
	} else {
		size := d.maxSize
		if size <= 0 {
			size = defaultMaxPeerPerItem
		}
		peerMap, _ := lru.New(size)
		peerMap.Add(p, value)
	}
	return nil
}

// Sync implements Datastore.Sync
func (d *LruMapDatastore) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

// Get implements Datastore.Get
func (d *LruMapDatastore) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	prefix, p := keySplit(key.String())
	set, ok := d.values.Get(prefix)

	if !ok {
		return nil, datastore.ErrNotFound
	}
	peerMap := set.(*lru.Cache)
	val, found := peerMap.Get(p)
	if !found {
		return nil, datastore.ErrNotFound
	}
	return val.([]byte), nil
}

// Has implements Datastore.Has
func (d *LruMapDatastore) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	prefix, p := keySplit(key.String())
	set, ok := d.values.Get(prefix)

	if !ok {
		return false, datastore.ErrNotFound
	}
	peerMap := set.(*lru.Cache)
	found := peerMap.Contains(p)
	return found, nil
}

// GetSize implements Datastore.GetSize
func (d *LruMapDatastore) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	prefix, p := keySplit(key.String())
	set, ok := d.values.Get(prefix)

	if !ok {
		return -1, datastore.ErrNotFound
	}
	peerMap := set.(*lru.Cache)
	val, found := peerMap.Get(p)
	if found {
		return len(val.([]byte)), nil
	}

	return -1, datastore.ErrNotFound
}

// Delete implements Datastore.Delete
func (d *LruMapDatastore) Delete(ctx context.Context, key datastore.Key) (err error) {
	prefix, p := keySplit(key.String())
	set, ok := d.values.Get(prefix)

	if !ok {
		return nil
	}
	peerMap := set.(*lru.Cache)
	peerMap.Remove(p)
	if peerMap.Len() == 0 {
		d.values.Remove(prefix)
	}
	return nil
}

// Query implements Datastore.Query
func (d *LruMapDatastore) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
	re := make([]dsq.Entry, 0, d.values.Len())
	if q.Prefix != "" {
		set, ok := d.values.Get(q.Prefix)
		if !ok {
			r := dsq.ResultsWithEntries(q, re)
			r = dsq.NaiveQueryApply(q, r)
			return r, nil
		}
		peerMap := set.(*lru.Cache)
		for _, ki := range peerMap.Keys() {
			kk := q.Prefix + ki.(string)
			k := datastore.NewKey(kk)
			vi, _ := d.values.Get(k)
			v := vi.([]byte)
			e := dsq.Entry{Key: k.String(), Size: len(v)}
			if !q.KeysOnly {
				e.Value = v
			}
			re = append(re, e)
		}
		r := dsq.ResultsWithEntries(q, re)
		r = dsq.NaiveQueryApply(q, r)
		return r, nil
	}
	for _, kPrefix := range d.values.Keys() {
		set, ok := d.values.Get(kPrefix)
		if !ok {
			continue
		}
		peerMap := set.(*lru.Cache)
		for _, ki := range peerMap.Keys() {
			kk := kPrefix.(string) + ki.(string)
			k := datastore.NewKey(kk)
			vi, _ := d.values.Get(k)
			v := vi.([]byte)
			e := dsq.Entry{Key: k.String(), Size: len(v)}
			if !q.KeysOnly {
				e.Value = v
			}
			re = append(re, e)
		}
	}
	r := dsq.ResultsWithEntries(q, re)
	r = dsq.NaiveQueryApply(q, r)
	return r, nil
}

func (d *LruMapDatastore) Batch(ctx context.Context) (datastore.Batch, error) {
	return datastore.NewBasicBatch(d), nil
}

func (d *LruMapDatastore) Close() error {
	return nil
}

func keySplit(key string) (keyName string, peer string) {
	// key format: /providers/data/peer
	i := strings.LastIndexByte(key, '/')
	if i != -1 {
		return key[:i], key[i:]
	}
	return "", ""
}
