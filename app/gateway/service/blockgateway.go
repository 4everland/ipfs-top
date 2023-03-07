package service

import (
	"context"
	"errors"
	"github.com/4everland/ipfs-servers/app/gateway/conf"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-libipfs/blocks"
	"github.com/ipfs/go-libipfs/files"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-namesys"
	iface "github.com/ipfs/interface-go-ipfs-core"
	ifaceoption "github.com/ipfs/interface-go-ipfs-core/options"
	nsopts "github.com/ipfs/interface-go-ipfs-core/options/namesys"
	ifacepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	mc "github.com/multiformats/go-multicodec"
)

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}

func NewUnixfs(bs blockstore.Blockstore) iface.UnixfsAPI {
	return coreunix.NewUnixFsServerOffline(bs)
}

func NewRouting() routing.ValueStore {
	return nil
}

type BlocksGateway struct {
	blockStore blockstore.Blockstore
	unixfs     iface.UnixfsAPI
	resolver   coreunix.DagResolve

	// Optional routing system to handle /ipns addresses.
	namesys namesys.NameSystem
	routing routing.ValueStore
}

func NewBlocksGateway(blockStore blockstore.Blockstore, unixfs iface.UnixfsAPI, routing routing.ValueStore) (*BlocksGateway, error) {
	bs := blockservice.New(blockStore, nil)
	dagService := merkledag.NewDAGService(bs)
	resolver := coreunix.NewDagResolver(context.Background(), dagService, bs)
	// Setup a name system so that we are able to resolve /ipns links.
	var (
		ns  namesys.NameSystem
		err error
	)
	if routing != nil {
		ns, err = namesys.NewNameSystem(routing)
		if err != nil {
			return nil, err
		}
	}

	return &BlocksGateway{
		blockStore: blockStore,
		unixfs:     unixfs,
		resolver:   resolver,

		routing: routing,
		namesys: ns,
	}, nil
}

func (api *BlocksGateway) GetUnixFsNode(ctx context.Context, p ifacepath.Resolved) (files.Node, error) {
	return api.unixfs.Get(ctx, p)
}

func (api *BlocksGateway) LsUnixFsDir(ctx context.Context, p ifacepath.Resolved) (<-chan iface.DirEntry, error) {
	return api.unixfs.Ls(ctx, p, ifaceoption.Unixfs.UseCumulativeSize(true))
}

func (api *BlocksGateway) GetBlock(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	return api.blockStore.Get(ctx, c)
}

func (api *BlocksGateway) GetIPNSRecord(ctx context.Context, c cid.Cid) ([]byte, error) {
	if api.routing == nil {
		return nil, routing.ErrNotSupported
	}

	// Fails fast if the CID is not an encoded Libp2p Key, avoids wasteful
	// round trips to the remote routing provider.
	if mc.Code(c.Type()) != mc.Libp2pKey {
		return nil, errors.New("provided cid is not an encoded libp2p key")
	}

	// The value store expects the key itself to be encoded as a multihash.
	id, err := peer.FromCid(c)
	if err != nil {
		return nil, err
	}

	return api.routing.GetValue(ctx, "/ipns/"+string(id))
}

func (api *BlocksGateway) GetDNSLinkRecord(ctx context.Context, hostname string) (ifacepath.Path, error) {
	if api.namesys != nil {
		p, err := api.namesys.Resolve(ctx, "/ipns/"+hostname, nsopts.Depth(1))
		if err == namesys.ErrResolveRecursion {
			err = nil
		}
		return ifacepath.New(p.String()), err
	}

	return nil, errors.New("not implemented")
}

func (api *BlocksGateway) IsCached(ctx context.Context, p ifacepath.Path) bool {
	rp, err := api.ResolvePath(ctx, p)
	if err != nil {
		return false
	}

	has, _ := api.blockStore.Has(ctx, rp.Cid())
	return has
}

func (api *BlocksGateway) ResolvePath(ctx context.Context, p ifacepath.Path) (ifacepath.Resolved, error) {
	return api.resolver.ResolvePath(ctx, p)
}
