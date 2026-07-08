package coreunix

import (
	"context"
	"fmt"

	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	exchange "github.com/ipfs/boxo/exchange"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/fetcher"
	bsfetcher "github.com/ipfs/boxo/fetcher/impl/blockservice"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	dagtest "github.com/ipfs/boxo/ipld/merkledag/test"
	ft "github.com/ipfs/boxo/ipld/unixfs"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	uio "github.com/ipfs/boxo/ipld/unixfs/io"
	"github.com/ipfs/boxo/mfs"
	ipfspath "github.com/ipfs/boxo/path"
	ipfspathresolver "github.com/ipfs/boxo/path/resolver"
	pin "github.com/ipfs/boxo/pinning/pinner"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfsnode"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/coreiface/options"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
	gopath "path"
)

type UnixFsServer struct {
	blockstore blockstore.GCBlockstore
	baseBlocks blockstore.Blockstore
	pinning    pin.Pinner
	exchange   exchange.Interface
	provider   providerNotify

	dag format.DAGService
	bs  blockservice.BlockService
}

type DagResolve interface {
	ResolvePath(ctx context.Context, p ipfspath.Path) (ResolvedPath, error)
	ResolveNode(ctx context.Context, p ipfspath.Path) (format.Node, error)
}

type ResolvedPath interface {
	ipfspath.Path
	Cid() cid.Cid
	Root() cid.Cid
	Remainder() string
}

type resolvedPath struct {
	ipfspath.Path
	cid       cid.Cid
	root      cid.Cid
	remainder string
}

func (p resolvedPath) Cid() cid.Cid {
	return p.cid
}

func (p resolvedPath) Root() cid.Cid {
	return p.root
}

func (p resolvedPath) Remainder() string {
	return p.remainder
}

type providerNotify interface {
	Provide(cid.Cid) error
}

type emptyGCBlockstore struct {
	blockstore.GCLocker
	blockstore.Blockstore
}

func NewUnixFsServer(baseBlocks blockstore.Blockstore, exchange exchange.Interface) *UnixFsServer {
	bs := blockservice.New(baseBlocks, exchange)
	dagService := merkledag.NewDAGService(bs)

	return &UnixFsServer{
		provider: emptyProviderNotify{},
		exchange: exchange,
		// todo... pinset
		pinning:    nil,
		baseBlocks: baseBlocks,
		blockstore: &emptyGCBlockstore{
			blockstore.NewGCLocker(),
			baseBlocks,
		},
		dag: dagService,
		bs:  bs,
	}

}

func NewUnixFsServerOffline(baseBlocks blockstore.Blockstore) *UnixFsServer {
	bs := blockservice.New(baseBlocks, nil)
	return &UnixFsServer{
		provider:   emptyProviderNotify{},
		exchange:   offline.Exchange(baseBlocks),
		pinning:    nil,
		baseBlocks: baseBlocks,
		blockstore: &emptyGCBlockstore{
			blockstore.NewGCLocker(),
			baseBlocks,
		},
		dag: merkledag.NewDAGService(bs),
		bs:  bs,
	}
}

func newEmptyGCBlockstore() blockstore.GCBlockstore {
	return &emptyGCBlockstore{blockstore.NewGCLocker(), blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))}
}

func (api *UnixFsServer) Add(ctx context.Context, files files.Node, opts ...options.UnixfsAddOption) (ipfspath.ImmutablePath, error) {
	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return ipfspath.ImmutablePath{}, err
	}

	addblockstore := api.blockstore

	exch := api.exchange
	pinning := api.pinning

	if settings.OnlyHash {
		addblockstore = newEmptyGCBlockstore()
		exch = offline.Exchange(api.baseBlocks)
		pinning = nil
	}

	bserv := blockservice.New(addblockstore, exch) // hash security 001
	dserv := merkledag.NewDAGService(bserv)

	fileAdder, err := NewAdder(ctx, pinning, addblockstore, dserv)
	if err != nil {
		return ipfspath.ImmutablePath{}, err
	}

	fileAdder.Chunker = settings.Chunker
	if settings.Events != nil {
		fileAdder.Out = settings.Events
		fileAdder.Progress = settings.Progress
	}
	fileAdder.Pin = settings.Pin && !settings.OnlyHash
	fileAdder.Silent = settings.Silent
	fileAdder.RawLeaves = settings.RawLeaves
	fileAdder.CidBuilder = prefix

	switch settings.Layout {
	case options.BalancedLayout:
		// Default
	case options.TrickleLayout:
		fileAdder.Trickle = true
	default:
		return ipfspath.ImmutablePath{}, fmt.Errorf("unknown layout: %d", settings.Layout)
	}

	if settings.OnlyHash {
		md := dagtest.Mock()
		emptyDirNode := ft.EmptyDirNode()
		// Use the same prefix for the "empty" MFS root as for the file rpc.
		err = emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
		if err != nil {
			return ipfspath.ImmutablePath{}, err
		}
		mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil, nil, mfs.WithCidBuilder(fileAdder.CidBuilder))
		if err != nil {
			return ipfspath.ImmutablePath{}, err
		}

		fileAdder.SetMfsRoot(mr)
	}

	nd, err := fileAdder.AddAllAndPin(ctx, files)
	if err != nil {
		return ipfspath.ImmutablePath{}, err
	}

	if !settings.OnlyHash {
		if err = api.provider.Provide(nd.Cid()); err != nil {
			return ipfspath.ImmutablePath{}, err
		}
	}

	return ipfspath.FromCid(nd.Cid()), nil
}

func (api *UnixFsServer) Get(ctx context.Context, p ipfspath.Path) (files.Node, error) {
	ses := newDagResolver(ctx, api.dag, api.bs)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	return unixfile.NewUnixfsFile(ctx, ses.dag, nd)
}

func (api *UnixFsServer) Ls(ctx context.Context, p ipfspath.Path, out chan<- coreiface.DirEntry, opts ...options.UnixfsLsOption) error {
	defer close(out)

	settings, err := options.UnixfsLsOptions(opts...)
	if err != nil {
		return err
	}

	ses := newDagResolver(ctx, api.dag, api.bs)

	dagnode, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return err
	}

	dir, err := uio.NewDirectoryFromNode(ses.dag, dagnode)
	if err == uio.ErrNotADir {
		return lsFromLinks(ctx, api.dag, dagnode.Links(), settings, out)
	}
	if err != nil {
		return err
	}

	return lsFromLinksAsync(ctx, api.dag, dir, settings, out)
}

type dagResolver struct {
	dag                  format.DAGService
	unixFSFetcherFactory fetcher.Factory
	ipldPathResolver     ipfspathresolver.Resolver
	unixFSPathResolver   ipfspathresolver.Resolver
}

func NewDagResolver(ctx context.Context, d format.NodeGetter, b blockservice.BlockService) *dagResolver {
	return newDagResolver(ctx, d, b)
}
func newDagResolver(ctx context.Context, d format.NodeGetter, b blockservice.BlockService) *dagResolver {
	fetcherConfig := bsfetcher.NewFetcherConfig(b)
	fetcherConfig.PrototypeChooser = dagpb.AddSupportToChooser(func(lnk ipld.Link, lnkCtx ipld.LinkContext) (ipld.NodePrototype, error) {
		if tlnkNd, ok := lnkCtx.LinkNode.(schema.TypedLinkNode); ok {
			return tlnkNd.LinkTargetNodePrototype(), nil
		}
		return basicnode.Prototype.Any, nil
	})
	fsFetcher := fetcherConfig.WithReifier(unixfsnode.Reify)
	return &dagResolver{
		dag:                  merkledag.NewReadOnlyDagService(merkledag.NewSession(ctx, d)),
		unixFSFetcherFactory: fsFetcher,
		ipldPathResolver:     ipfspathresolver.NewBasicResolver(fetcherConfig),
		unixFSPathResolver:   ipfspathresolver.NewBasicResolver(fsFetcher),
	}
}

func (dr *dagResolver) ResolvePath(ctx context.Context, p ipfspath.Path) (ResolvedPath, error) {
	if rp, ok := p.(ResolvedPath); ok {
		return rp, nil
	}

	ipath, err := ipfspath.NewPath(p.String())
	if err != nil {
		return nil, err
	}

	if ipath.Segments()[0] != "ipfs" && ipath.Segments()[0] != "ipld" {
		return nil, fmt.Errorf("unsupported path namespace: %s", p.Namespace())
	}

	immutablePath, err := ipfspath.NewImmutablePath(ipath)
	if err != nil {
		return nil, err
	}

	var resolver ipfspathresolver.Resolver
	if ipath.Segments()[0] == "ipld" {
		resolver = dr.ipldPathResolver
	} else {
		resolver = dr.unixFSPathResolver
	}

	node, rest, err := resolver.ResolveToLastNode(ctx, immutablePath)
	if err != nil {
		return nil, err
	}

	root, err := cid.Parse(ipath.Segments()[1])
	if err != nil {
		return nil, err
	}

	return resolvedPath{
		Path:      ipath,
		cid:       node,
		root:      root,
		remainder: gopath.Join(rest...),
	}, nil
}

func (dr *dagResolver) ResolveNode(ctx context.Context, p ipfspath.Path) (format.Node, error) {
	rp, err := dr.ResolvePath(ctx, p)
	if err != nil {
		return nil, err
	}

	node, err := dr.dag.Get(ctx, rp.Cid())
	if err != nil {
		return nil, err
	}
	return node, nil
}

type emptyProviderNotify struct {
}

func (emptyProviderNotify) Provide(cid.Cid) error {
	return nil
}

func processLink(ctx context.Context, dag format.NodeGetter, linkres ft.LinkResult, settings *options.UnixfsLsSettings) (coreiface.DirEntry, error) {

	if linkres.Link != nil {
		//span.SetAttributes(attribute.String("linkname", linkres.Link.Name), attribute.String("cid", linkres.Link.Cid.String()))
	}

	if linkres.Err != nil {
		return coreiface.DirEntry{}, linkres.Err
	}

	lnk := coreiface.DirEntry{
		Name: linkres.Link.Name,
		Cid:  linkres.Link.Cid,
	}

	switch lnk.Cid.Type() {
	case cid.Raw:
		// No need to check with raw leaves
		lnk.Type = coreiface.TFile
		lnk.Size = linkres.Link.Size
	case cid.DagProtobuf:
		if settings.ResolveChildren {
			linkNode, err := linkres.Link.GetNode(ctx, dag)
			if err != nil {
				return coreiface.DirEntry{}, err
			}

			if pn, ok := linkNode.(*merkledag.ProtoNode); ok {
				d, err := ft.FSNodeFromBytes(pn.Data())
				if err != nil {
					return coreiface.DirEntry{}, err
				}
				switch d.Type() {
				case ft.TFile, ft.TRaw:
					lnk.Type = coreiface.TFile
				case ft.THAMTShard, ft.TDirectory, ft.TMetadata:
					lnk.Type = coreiface.TDirectory
				case ft.TSymlink:
					lnk.Type = coreiface.TSymlink
					lnk.Target = string(d.Data())
				}
				if !settings.UseCumulativeSize {
					lnk.Size = d.FileSize()
				}
			}
		}

		if settings.UseCumulativeSize {
			lnk.Size = linkres.Link.Size
		}
	}

	return lnk, nil
}

func lsFromLinksAsync(ctx context.Context, dag format.NodeGetter, dir uio.Directory, settings *options.UnixfsLsSettings, out chan<- coreiface.DirEntry) error {
	for l := range dir.EnumLinksAsync(ctx) {
		entry, err := processLink(ctx, dag, l, settings)
		if err != nil {
			return err
		}
		select {
		case out <- entry: // TODO: perf: processing can be done in background and in parallel
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func lsFromLinks(ctx context.Context, dag format.NodeGetter, ndlinks []*format.Link, settings *options.UnixfsLsSettings, out chan<- coreiface.DirEntry) error {
	for _, l := range ndlinks {
		lr := ft.LinkResult{Link: &format.Link{Name: l.Name, Size: l.Size, Cid: l.Cid}}

		entry, err := processLink(ctx, dag, lr, settings)
		if err != nil {
			return err
		}
		select {
		case out <- entry: // TODO: can be parallel if settings.Async
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
