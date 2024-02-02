package coreunix

import (
	"context"
	"fmt"
	"github.com/4everland/ipfs-top/third_party/coreunix/options"
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/boxo/path/resolver"
	//coreiface "github.com/ipfs/boxo/coreiface"
	//"github.com/ipfs/boxo/coreiface/options"
	//"github.com/ipfs/boxo/coreiface/path"
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
	//ipafspathresolver "github.com/ipfs/boxo/path/resolver"
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

// FileType is an enum of possible UnixFS file types.
type FileType int32

const (
	// TUnknown means the file type isn't known (e.g., it hasn't been
	// resolved).
	TUnknown FileType = iota
	// TFile is a regular file.
	TFile
	// TDirectory is a directory.
	TDirectory
	// TSymlink is a symlink.
	TSymlink
)

func (t FileType) String() string {
	switch t {
	case TUnknown:
		return "unknown"
	case TFile:
		return "file"
	case TDirectory:
		return "directory"
	case TSymlink:
		return "symlink"
	default:
		return "<unknown file type>"
	}
}

// DirEntry is a directory entry returned by `Ls`.
type DirEntry struct {
	Name string
	Cid  cid.Cid

	// Only filled when asked to resolve the directory entry.
	Size   uint64   // The size of the file in bytes (or the size of the symlink).
	Type   FileType // The type of the file.
	Target string   // The symlink target (if a symlink).

	Err error
}

// UnixfsAPI is the basic interface to immutable files in IPFS
// NOTE: This API is heavily WIP, things are guaranteed to break frequently
type UnixfsAPI interface {
	// Add imports the data from the reader into merkledag file
	//
	// TODO: a long useful comment on how to use this for many different scenarios
	Add(context.Context, files.Node, ...options.UnixfsAddOption) (path.ImmutablePath, error)

	// Get returns a read-only handle to a file tree referenced by a path
	//
	// Note that some implementations of this API may apply the specified context
	// to operations performed on the returned file
	Get(context.Context, path.Path) (files.Node, error)

	// Ls returns the list of links in a directory. Links aren't guaranteed to be
	// returned in order
	Ls(context.Context, path.Path, ...options.UnixfsLsOption) (<-chan DirEntry, error)
}

type DagResolve interface {
	ResolvePath(ctx context.Context, p path.Path) (path.ImmutablePath, error)
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

func (api *UnixFsServer) Add(ctx context.Context, files files.Node, opts ...options.UnixfsAddOption) (path.ImmutablePath, error) {
	settings, prefix, err := options.UnixfsAddOptions(opts...)
	if err != nil {
		return path.ImmutablePath{}, err
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
		return path.ImmutablePath{}, err
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
		return path.ImmutablePath{}, fmt.Errorf("unknown layout: %d", settings.Layout)
	}

	if settings.OnlyHash {
		md := dagtest.Mock()
		emptyDirNode := ft.EmptyDirNode()
		// Use the same prefix for the "empty" MFS root as for the file rpc.
		err = emptyDirNode.SetCidBuilder(fileAdder.CidBuilder)
		if err != nil {
			return path.ImmutablePath{}, err
		}
		mr, err := mfs.NewRoot(ctx, md, emptyDirNode, nil)
		if err != nil {
			return path.ImmutablePath{}, err
		}

		fileAdder.SetMfsRoot(mr)
	}

	nd, err := fileAdder.AddAllAndPin(ctx, files)
	if err != nil {
		return path.ImmutablePath{}, err
	}

	if !settings.OnlyHash {
		if err = api.provider.Provide(nd.Cid()); err != nil {
			return path.ImmutablePath{}, err
		}
	}

	return path.FromCid(nd.Cid()), nil
}

func (api *UnixFsServer) Get(ctx context.Context, p path.Path) (files.Node, error) {

	ses := newDagResolver(ctx, api.dag, api.bs)

	nd, err := ses.ResolveNode(ctx, p)
	if err != nil {
		return nil, err
	}

	return unixfile.NewUnixfsFile(ctx, ses.dag, nd)
}

func (api *UnixFsServer) Ls(ctx context.Context, p path.Path, opts ...options.UnixfsLsOption) (<-chan DirEntry, error) {
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
	ipldPathResolver     resolver.Resolver
	unixFSPathResolver   resolver.Resolver
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
		ipldPathResolver:     resolver.NewBasicResolver(fetcherConfig),
		unixFSPathResolver:   resolver.NewBasicResolver(fsFetcher),
	}
}

func (dr *dagResolver) ResolvePath(ctx context.Context, p path.Path) (path.ImmutablePath, error) {
	pp, err := path.NewImmutablePath(p)
	if err != nil {
		return path.ImmutablePath{}, err
	}
	if pp.Segments()[0] != "ipfs" && pp.Segments()[0] != "ipld" {
		return path.ImmutablePath{}, fmt.Errorf("unsupported path namespace: %s", p.Namespace())
	}

	var r resolver.Resolver
	if pp.Segments()[0] == "ipld" {
		r = dr.ipldPathResolver
	} else {
		r = dr.unixFSPathResolver
	}

	node, _, err := r.ResolveToLastNode(ctx, pp)
	if err != nil {
		return path.ImmutablePath{}, err
	}

	return path.FromCid(node), nil
}

func (dr *dagResolver) ResolveNode(ctx context.Context, p path.Path) (format.Node, error) {

	rp, err := dr.ResolvePath(ctx, p)
	if err != nil {
		return nil, err
	}

	node, err := dr.dag.Get(ctx, rp.RootCid())
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

func processLink(ctx context.Context, dag format.NodeGetter, linkres ft.LinkResult, settings *options.UnixfsLsSettings) DirEntry {

	if linkres.Link != nil {
		//span.SetAttributes(attribute.String("linkname", linkres.Link.Name), attribute.String("cid", linkres.Link.Cid.String()))
	}

	if linkres.Err != nil {
		return DirEntry{Err: linkres.Err}
	}

	lnk := DirEntry{
		Name: linkres.Link.Name,
		Cid:  linkres.Link.Cid,
	}

	switch lnk.Cid.Type() {
	case cid.Raw:
		// No need to check with raw leaves
		lnk.Type = TFile
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
					lnk.Type = TFile
				case ft.THAMTShard, ft.TDirectory, ft.TMetadata:
					lnk.Type = TDirectory
				case ft.TSymlink:
					lnk.Type = TSymlink
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

func lsFromLinksAsync(ctx context.Context, dag format.NodeGetter, dir uio.Directory, settings *options.UnixfsLsSettings) (<-chan DirEntry, error) {
	out := make(chan DirEntry, uio.DefaultShardWidth)

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

func lsFromLinks(ctx context.Context, dag format.NodeGetter, ndlinks []*format.Link, settings *options.UnixfsLsSettings) (<-chan DirEntry, error) {
	links := make(chan DirEntry, len(ndlinks))
	for _, l := range ndlinks {
		lr := ft.LinkResult{Link: &format.Link{Name: l.Name, Size: l.Size, Cid: l.Cid}}

		links <- processLink(ctx, dag, lr, settings) //TODO: can be parallel if settings.Async
	}
	close(links)
	return links, nil
}
