package service

import (
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	cmds "github.com/ipfs/go-ipfs-cmds"
	http2 "github.com/ipfs/go-ipfs-cmds/http"
	"github.com/ipfs/go-unixfs"
	unixfs_pb "github.com/ipfs/go-unixfs/pb"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"sort"
)

type LsService struct {
	unixfs *coreunix.UnixFsServer
}

func NewLsService(unixfs *coreunix.UnixFsServer) *LsService {
	return &LsService{unixfs}
}

func (s *LsService) RegisterRoute(route *httpctx.Router) {
	route.POST("/ls", s.Ls)
}

// LsLink contains printable data for a single ipld link in ls output
type LsLink struct {
	Name, Hash string
	Size       uint64
	Type       unixfs_pb.Data_DataType
	Target     string
}

// LsObject is an element of LsOutput
// It can represent all or part of a directory
type LsObject struct {
	Hash  string
	Links []LsLink
}

// LsOutput is a set of printable data for directories,
// it can be complete or partial
type LsOutput struct {
	Objects []LsObject
}

type LsRequest struct {
	Arg         []string `json:"arg,omitempty"`
	ResolveType bool     `json:"resolve-type,omitempty"`
	Size        bool     `json:"size,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
}

func (s *LsService) Ls(ctx httpctx.Context) (err error) {
	var req LsRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}

	if ctx.Request().URL.Query().Get("resolve-type") == "" {
		req.ResolveType = true
	}

	if ctx.Request().URL.Query().Get("size") == "" {
		req.Size = true
	}

	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})
	if err != nil {
		return err
	}

	var processLink func(path string, link LsLink) error
	var dirDone func(i int)

	processDir := func() (func(path string, link LsLink) error, func(i int)) {
		return func(path string, link LsLink) error {
			output := []LsObject{{
				Hash:  path,
				Links: []LsLink{link},
			}}
			return res.Emit(&LsOutput{output})
		}, func(i int) {}
	}
	done := func() error { return nil }

	if !req.Stream {
		output := make([]LsObject, len(req.Arg))

		processDir = func() (func(path string, link LsLink) error, func(i int)) {
			// for each dir
			outputLinks := make([]LsLink, 0)
			return func(path string, link LsLink) error {
					// for each link
					outputLinks = append(outputLinks, link)
					return nil
				}, func(i int) {
					// after each dir
					sort.Slice(outputLinks, func(i, j int) bool {
						return outputLinks[i].Name < outputLinks[j].Name
					})

					output[i] = LsObject{
						Hash:  req.Arg[i],
						Links: outputLinks,
					}
				}
		}

		done = func() error {
			return cmds.EmitOnce(res, &LsOutput{output})
		}
	}

	for i, fpath := range req.Arg {
		results, err := s.unixfs.Ls(ctx, path.New(fpath),
			options.Unixfs.ResolveChildren(req.Size || req.ResolveType))
		if err != nil {
			return err
		}

		processLink, dirDone = processDir()
		for link := range results {
			if link.Err != nil {
				return link.Err
			}
			var ftype unixfs_pb.Data_DataType
			switch link.Type {
			case iface.TFile:
				ftype = unixfs.TFile
			case iface.TDirectory:
				ftype = unixfs.TDirectory
			case iface.TSymlink:
				ftype = unixfs.TSymlink
			}
			lsLink := LsLink{
				Name: link.Name,
				Hash: link.Cid.String(),

				Size:   link.Size,
				Type:   ftype,
				Target: link.Target,
			}
			if err := processLink(req.Arg[i], lsLink); err != nil {
				return err
			}
		}
		dirDone(i)
	}
	return done()
}
