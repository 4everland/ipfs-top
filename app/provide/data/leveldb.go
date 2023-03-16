package data

import (
	"context"
	"github.com/4everland/ipfs-servers/api/contentrouting"
	"github.com/4everland/ipfs-servers/app/provide/conf"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/ipfs/go-ipfs-provider/simple"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewLevelDbDatastore(data *conf.Data) ds.Datastore {
	datastore, err := leveldb.NewDatastore(data.LeveldbPath, &leveldb.Options{})
	if err != nil {
		panic(err)
	}
	return datastore
}

func NewContentRouting(data *conf.Data) routing.ContentRouting {
	tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(data.NodeUri, tlsOption)
	if err != nil {
		panic(err)
	}
	client := contentrouting.NewContentRoutingClient(conn)
	return &GrpcContentRouting{
		client: client,
	}
}

func NewKeyChanFunc(data *conf.Data) simple.KeyChanFunc {
	s, err := dag.NewBlockStore(data.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return simple.NewBlockstoreProvider(s)
}

type GrpcContentRouting struct {
	client contentrouting.ContentRoutingClient
}

// Provide adds the given cid to the content routing system. If 'true' is
// passed, it also announces it, otherwise it is just kept in the local
// accounting of which objects are being provided.
func (grt *GrpcContentRouting) Provide(ctx context.Context, c cid.Cid, provide bool) error {
	_, err := grt.client.Provide(ctx, &contentrouting.ProvideReq{
		Cid:     &contentrouting.Cid{Str: c.Bytes()},
		Provide: provide,
	})
	return err
}

// Search for peers who are able to provide a given key
//
// When count is 0, this method will return an unbounded number of
// results.
func (grt *GrpcContentRouting) FindProvidersAsync(_ context.Context, _ cid.Cid, _ int) <-chan peer.AddrInfo {
	return nil
}
