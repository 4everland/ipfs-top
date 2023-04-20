package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/coreiface/path"
	dag "github.com/ipfs/boxo/ipld/merkledag"
	ft "github.com/ipfs/boxo/ipld/unixfs"
	cmds "github.com/ipfs/go-ipfs-cmds"
	ipld "github.com/ipfs/go-ipld-format"
	"net/http"
	gopath "path"
	"strings"
)

type FilesService struct {
	dagService        ipld.DAGService
	offlineDagService ipld.DAGService
	dagResolver       coreunix.DagResolve
}

func NewFilesService(dagService ipld.DAGService, offlineBs OfflineBlockService, resolver coreunix.DagResolve) *FilesService {
	return &FilesService{
		dagService:        dagService,
		offlineDagService: dag.NewDAGService(offlineBs),
		dagResolver:       resolver,
	}
}

func (s *FilesService) RegisterRoute(route *httpctx.Router) {
	route.POST("/files/stat", s.Stat)
	route.POST("/object/stat", s.ObjectStat)
}

type StatRequest struct {
	Arg       string `json:"arg,omitempty"`
	Format    string `json:"format,omitempty"`
	Size      bool   `json:"size,omitempty"`
	Hash      bool   `json:"hash,omitempty"`
	WithLocal bool   `json:"with-local,omitempty"`
}

func (s *FilesService) Stat(ctx httpctx.Context) (err error) {
	w := ctx.Response()

	var req StatRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}

	if ctx.Request().URL.Query().Get("format") == "" {
		req.Format = defaultStatFormat
	}

	if _, err = statGetFormatOptions(req.Hash, req.Size, req.Format); err != nil {
		return cmds.Errorf(cmds.ErrClient, err.Error())
	}

	p, err := checkPath(req.Arg)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(p, "/ipfs/") {
		return errors.New("path not support")
	}

	var dagserv = s.dagService
	if req.WithLocal {
		dagserv = s.offlineDagService
	}

	nd, err := s.dagResolver.ResolveNode(ctx, path.New(p))
	if err != nil {
		return err
	}

	o, err := statNode(nd)
	if err != nil {
		return err
	}

	if !req.WithLocal {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(o)
		w.Write(b)
		return nil
	}

	local, sizeLocal, err := walkBlock(ctx, dagserv, nd)
	if err != nil {
		return err
	}

	o.WithLocality = true
	o.Local = local
	o.SizeLocal = sizeLocal
	b, _ := json.Marshal(o)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

	return nil
}

func checkPath(p string) (string, error) {
	if len(p) == 0 {
		return "", fmt.Errorf("paths must not be empty")
	}

	if p[0] != '/' {
		return "", fmt.Errorf("paths must start with a leading slash")
	}

	cleaned := gopath.Clean(p)
	if p[len(p)-1] == '/' && p != "/" {
		cleaned += "/"
	}
	return cleaned, nil
}

const (
	defaultStatFormat = `<hash>
Size: <size>
CumulativeSize: <cumulsize>
ChildBlocks: <childs>
Type: <type>`
)

func statGetFormatOptions(hash, size bool, format string) (string, error) {
	f := format != defaultStatFormat
	if hash && size || size && f || hash && f {
		return "", errors.New("format was set by multiple options. Only one format option is allowed")
	}

	if hash {
		return "<hash>", nil
	} else if size {
		return "<cumulsize>", nil
	} else {
		return format, nil
	}
}

type statOutput struct {
	Hash           string
	Size           uint64
	CumulativeSize uint64
	Blocks         int
	Type           string
	WithLocality   bool   `json:",omitempty"`
	Local          bool   `json:",omitempty"`
	SizeLocal      uint64 `json:",omitempty"`
}

func statNode(nd ipld.Node) (*statOutput, error) {
	c := nd.Cid()

	cumulsize, err := nd.Size()
	if err != nil {
		return nil, err
	}

	switch n := nd.(type) {
	case *dag.ProtoNode:
		d, err := ft.FSNodeFromBytes(n.Data())
		if err != nil {
			return nil, err
		}

		var ndtype string
		switch d.Type() {
		case ft.TDirectory, ft.THAMTShard:
			ndtype = "directory"
		case ft.TFile, ft.TMetadata, ft.TRaw:
			ndtype = "file"
		default:
			return nil, fmt.Errorf("unrecognized node type: %s", d.Type())
		}

		return &statOutput{
			Hash:           c.String(),
			Blocks:         len(nd.Links()),
			Size:           d.FileSize(),
			CumulativeSize: cumulsize,
			Type:           ndtype,
		}, nil
	case *dag.RawNode:
		return &statOutput{
			Hash:           c.String(),
			Blocks:         0,
			Size:           cumulsize,
			CumulativeSize: cumulsize,
			Type:           "file",
		}, nil
	default:
		return nil, fmt.Errorf("not unixfs node (proto or raw)")
	}
}

func walkBlock(ctx context.Context, dagserv ipld.DAGService, nd ipld.Node) (bool, uint64, error) {
	// Start with the block data size
	sizeLocal := uint64(len(nd.RawData()))

	local := true

	for _, link := range nd.Links() {
		child, err := dagserv.Get(ctx, link.Cid)

		if ipld.IsNotFound(err) {
			local = false
			continue
		}

		if err != nil {
			return local, sizeLocal, err
		}

		childLocal, childLocalSize, err := walkBlock(ctx, dagserv, child)

		if err != nil {
			return local, sizeLocal, err
		}

		// Recursively add the child size
		local = local && childLocal
		sizeLocal += childLocalSize
	}

	return local, sizeLocal, nil
}

func (s *FilesService) ObjectStat(ctx httpctx.Context) (err error) {
	w := ctx.Response()

	p := ctx.Query().Get("arg")
	nd, err := s.dagResolver.ResolveNode(ctx, path.New(p))
	if err != nil {
		return err
	}

	stat, err := nd.Stat()
	if err != nil {
		return err
	}

	b, _ := json.Marshal(stat)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

	return nil
}
