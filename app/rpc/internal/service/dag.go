package service

import (
	"errors"
	"fmt"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	cmds "github.com/ipfs/go-ipfs-cmds"
	http2 "github.com/ipfs/go-ipfs-cmds/http"
	ipld "github.com/ipfs/go-ipld-format"
	ipldlegacy "github.com/ipfs/go-ipld-legacy"
	gocarv2 "github.com/ipld/go-car/v2"
	ipldp "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/traversal"
	mc "github.com/multiformats/go-multicodec"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

type DagService struct {
	dag         ipld.DAGService
	dagResolver coreunix.DagResolve
	decoder     *ipldlegacy.Decoder
}

func NewDagService(dag ipld.DAGService, resolver coreunix.DagResolve) *DagService {
	return &DagService{
		dag:         dag,
		dagResolver: resolver,
		decoder:     ipldlegacy.NewDecoder(),
	}
}

func (s *DagService) RegisterRoute(route *httpctx.Router) {
	route.POST("/dag/import", s.DagImport)
	route.POST("/dag/get", s.DagGet)
}

type DagImportRequest struct {
	PinRoots      bool `json:"pin-roots,omitempty"`
	Silent        bool `json:"silent,omitempty"`
	Stats         bool `json:"stats,omitempty"`
	AllowBigBlock bool `json:"allow-big-block,omitempty"`
}

func (s *DagService) DagImport(ctx httpctx.Context) (err error) {
	w := ctx.Response()
	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})
	r := ctx.Request()
	contentType := r.Header.Get(contentTypeHeader)
	mediatype, _, _ := mime.ParseMediaType(contentType)

	var req DagImportRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}
	if r.URL.Query().Get("pin-roots") == "" {
		req.PinRoots = true
	}

	// todo dag pinset
	batch := ipld.NewBatch(ctx, s.dag)

	roots := cid.NewSet()
	var blockCount, blockBytesCount uint64

	var f files.Directory
	if mediatype == "multipart/form-data" {
		var reader *multipart.Reader
		if reader, err = r.MultipartReader(); err != nil {
			return
		}

		if f, err = files.NewFileFromPartReader(reader, mediatype); err != nil {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Chunked-Output", "1")
	w.WriteHeader(http.StatusOK)
	it := f.Entries()
	for it.Next() {
		file := files.FileFromEntry(it)
		if file == nil {
			return errors.New("expected a file handle")
		}

		// import blocks
		err = func() error {
			defer file.Close()

			car, err := gocarv2.NewBlockReader(file)
			if err != nil {
				return err
			}

			for _, c := range car.Roots {
				roots.Add(c)
			}

			for {
				block, err := car.Next()
				if err != nil && err != io.EOF {
					return err
				} else if block == nil {
					break
				}

				if !req.AllowBigBlock && len(block.RawData()) > 1024*1024 {
					return fmt.Errorf("produced block is over 1MiB: big blocks can't be exchanged with other peers. consider using UnixFS for automatic chunking of bigger files, or pass --allow-big-block to override")
				}

				nd, err := s.decoder.DecodeNode(ctx, block)
				if err != nil {
					return err
				}

				if err := batch.Add(ctx, nd); err != nil {
					return err
				}
				blockCount++
				blockBytesCount += uint64(len(block.RawData()))
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	if err := batch.Commit(); err != nil {
		return err
	}

	if req.PinRoots {
		err = roots.ForEach(func(c cid.Cid) error {
			ret := RootMeta{Cid: c}
			if block, err := s.dag.Get(ctx, c); err != nil {
				ret.PinErrorMsg = err.Error()
			} else if _, err := s.decoder.DecodeNode(ctx, block); err != nil {
				ret.PinErrorMsg = err.Error()
			}
			//todo pinset

			return res.Emit(&CarImportOutput{Root: &ret})
		})
		if err != nil {
			return err
		}
	}

	if req.Stats {
		err = res.Emit(&CarImportOutput{
			Stats: &CarImportStats{
				BlockCount:      blockCount,
				BlockBytesCount: blockBytesCount,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type RootMeta struct {
	Cid         cid.Cid
	PinErrorMsg string
}

type CarImportStats struct {
	BlockCount      uint64
	BlockBytesCount uint64
}

type CarImportOutput struct {
	Root  *RootMeta       `json:",omitempty"`
	Stats *CarImportStats `json:",omitempty"`
}

type DagGetRequest struct {
	Arg         string `json:"arg,omitempty"`
	OutputCodec string `json:"output-codec,omitempty"`
}

func (s *DagService) DagGet(ctx httpctx.Context) (err error) {
	w := ctx.Response()
	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})

	var req DagGetRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}
	if req.Arg == "" {
		return errors.New("argument \"arg\" is required")
	}

	if req.OutputCodec == "" {
		req.OutputCodec = "dag-json"
	}
	var codec mc.Code
	if err := codec.Set(req.OutputCodec); err != nil {
		return err
	}

	rp, err := s.dagResolver.ResolvePath(ctx, path.New(req.Arg))
	if err != nil {
		return err
	}

	obj, err := s.dag.Get(ctx, rp.Cid())
	if err != nil {
		return err
	}

	universal, ok := obj.(ipldlegacy.UniversalNode)
	if !ok {
		return fmt.Errorf("%T is not a valid IPLD node", obj)
	}

	finalNode := universal.(ipldp.Node)

	if len(rp.Remainder()) > 0 {
		remainderPath := ipldp.ParsePath(rp.Remainder())

		finalNode, err = traversal.Get(finalNode, remainderPath)
		if err != nil {
			return err
		}
	}

	encoder, err := multicodec.LookupEncoder(uint64(codec))
	if err != nil {
		return fmt.Errorf("invalid encoding: %s - %s", codec, err)
	}

	r, wr := io.Pipe()
	go func() {
		defer wr.Close()
		if err := encoder(finalNode, w); err != nil {
			_ = wr.CloseWithError(err)
		}
	}()

	return res.Emit(r)
}
