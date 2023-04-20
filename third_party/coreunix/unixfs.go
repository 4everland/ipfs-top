package coreunix

import (
	"context"
	"fmt"
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	coreiface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
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
	ResolvePath(ctx context.Context, p path.Path) (path.Resolved, error)
	ResolveNode(ctx context.Context, p path.Path) (format.Node, error)
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

func (api *UnixFsServer) Add(ctx context.Context, files files.Node, opts ...options.UnixfsAddOption) (path.Resolved, error) {
	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, fmt.Errorf("unknown layout: %d", settings.Layout)
	}

	if settings.OnlyHash {
		md := dagtest.Mock()
		emptyDirNode := ft.EmptyDirNode()
		// Use the same prefix for the "empty" MFS root as for the file rpc.
		err = emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
		if err != nil {
			return nil, err
		}
		mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil)
		if err != nil {
			return nil, err
		}

		fileAdder.SetMfsRoot(mr)
	}

	nd, err := fileAdder.AddAllAndPin(ctx, files)
	if err != nil {
		return nil, err
	}

	if !settings.OnlyHash {
		if err = api.provider.Provide(nd.Cid()); err != nil {
			return nil, err
		}
	}

	return path.IpfsPath(nd.Cid()), nil
}

func (api *UnixFsServer) Get(ctx context.Context, p path.Path) (files.Node, error) {
	ses := newDagResolver(ctx, api.dag, api.bs)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	return unixfile.NewUnixfsFile(ctx, ses.dag, nd)
}

func (api *UnixFsServer) Ls(ctx context.Context, p path.Path, opts ...options.UnixfsLsOption) (<-chan coreiface.DirEntry, error) {
	settings, err := options.UnixfsLsOptions(opts...)
	if err != nil {
		return nil, err
	}

	ses := newDagResolver(ctx, api.dag, api.bs)

	dagnode, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	dir, err := uio.NewDirectoryFromNode(ses.dag, dagnode)
	if err == uio.ErrNotADir {
		return lsFromLinks(ctx, api.dag, dagnode.Links(), settings)
	}
	if err != nil {
		return nil, err
	}

	return lsFromLinksAsync(ctx, api.dag, dir, settings)
}

type dagResolver struct {
	dag                  format.DAGService
	unixFSFetcherFactory fetcher.Factory
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
	//fetcher := fetcherConfig.WithReifier(unixfsnode.Reify)
	return &dagResolver{
		dag:                  merkledag.NewReadOnlyDagService(merkledag.NewSession(ctx, d)),
		unixFSFetcherFactory: fetcherConfig.WithReifier(unixfsnode.Reify),
	}
}

func (dr *dagResolver) ResolvePath(ctx context.Context, p path.Path) (path.Resolved, error) {
	if _, ok := p.(path.Resolved); ok {
		return p.(path.Resolved), nil
	}
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	ipath := ipfspath.Path(p.String())
	//if ipath.Segments()[0] != "ipfs" && ipath.Segments()[0] != "ipld" {
	if ipath.Segments()[0] != "ipfs" {
		return nil, fmt.Errorf("unsupported path namespace: %s", p.Namespace())
	}

	resolver := ipfspathresolver.NewBasicResolver(dr.unixFSFetcherFactory)

	node, rest, err := resolver.ResolveToLastNode(ctx, ipath)
	if err != nil {
		return nil, err
	}

	root, err := cid.Parse(ipath.Segments()[1])
	if err != nil {
		return nil, err
	}

	return path.NewResolvedPath(ipath, node, root, gopath.Join(rest...)), nil
}

func (dr *dagResolver) ResolveNode(ctx context.Context, p path.Path) (format.Node, error) {
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

func processLink(ctx context.Context, dag format.NodeGetter, linkres ft.LinkResult, settings *options.UnixfsLsSettings) coreiface.DirEntry {

	if linkres.Link != nil {
		//span.SetAttributes(attribute.String("linkname", linkres.Link.Name), attribute.String("cid", linkres.Link.Cid.String()))
	}

	if linkres.Err != nil {
		return coreiface.DirEntry{Err: linkres.Err}
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
				lnk.Err = err
				break
			}

			if pn, ok := linkNode.(*merkledag.ProtoNode); ok {
				d, err := ft.FSNodeFromBytes(pn.Data())
				if err != nil {
					lnk.Err = err
					break
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

	return lnk
}

func lsFromLinksAsync(ctx context.Context, dag format.NodeGetter, dir uio.Directory, settings *options.UnixfsLsSettings) (<-chan coreiface.DirEntry, error) {
	out := make(chan coreiface.DirEntry, uio.DefaultShardWidth)

	go func() {
		defer close(out)
		for l := range dir.EnumLinksAsync(ctx) {
			select {
			case out <- processLink(ctx, dag, l, settings): //TODO: perf: processing can be done in background and in parallel
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func lsFromLinks(ctx context.Context, dag format.NodeGetter, ndlinks []*format.Link, settings *options.UnixfsLsSettings) (<-chan coreiface.DirEntry, error) {
	links := make(chan coreiface.DirEntry, len(ndlinks))
	for _, l := range ndlinks {
		lr := ft.LinkResult{Link: &format.Link{Name: l.Name, Size: l.Size, Cid: l.Cid}}

		links <- processLink(ctx, dag, lr, settings) //TODO: can be parallel if settings.Async
	}
	close(links)
	return links, nil
}
