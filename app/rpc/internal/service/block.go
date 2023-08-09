package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	blocks "github.com/ipfs/go-block-format"
	cmds "github.com/ipfs/go-ipfs-cmds"
	http2 "github.com/ipfs/go-ipfs-cmds/http"
	mh "github.com/multiformats/go-multihash"
	"io"
	"mime"
	"mime/multipart"
)

type BlocksService struct {
	bs          blockservice.BlockService
	dagResolver coreunix.DagResolve
}

func NewBlocksService(bs blockservice.BlockService, dagResolver coreunix.DagResolve) *BlocksService {
	return &BlocksService{
		bs:          bs,
		dagResolver: dagResolver,
	}
}

func (s *BlocksService) RegisterRoute(route *httpctx.Router) {
	route.POST("/block/put", s.BlockPut)
	route.POST("/block/get", s.BlockGet)
}

type BlockPutRequest struct {
	CidCodec      string `json:"cid-codec,omitempty"`
	Mhtype        string `json:"mhtype,omitempty"`
	Mhlen         int    `json:"mhlen,omitempty"`
	Pin           bool   `json:"pin,omitempty"`
	AllowBigBlock bool   `json:"allow-big-Block,omitempty"`
	Format        string `json:"format,omitempty"`
}

func (s *BlocksService) BlockPut(ctx httpctx.Context) (err error) {
	w := ctx.Response()
	r := ctx.Request()
	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})
	contentType := r.Header.Get(contentTypeHeader)
	mediatype, _, _ := mime.ParseMediaType(contentType)

	var req BlockPutRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}

	mhtval, ok := mh.Names[req.Mhtype]
	if !ok {
		return fmt.Errorf("unrecognized multihash function: %s", req.Mhtype)
	}

	if r.URL.Query().Get("mhlen") == "" {
		return errors.New("missing option \"mhlen\"")
	}

	if req.Format != "" {
		if req.CidCodec != "" && req.CidCodec != "raw" {
			return errors.New("unable to use formart (deprecated) and a custom cid-codec at the same time")
		}
		req.CidCodec = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Chunked-Output", "1")
	//w.WriteHeader(http.StatusOK)
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

	it := f.Entries()
	for it.Next() {
		file := files.FileFromEntry(it)
		if file == nil {
			return errors.New("expected a file handle")
		}

		settings, err := options.BlockPutOptions(options.Block.Hash(mhtval, req.Mhlen),
			options.Block.CidCodec(req.CidCodec),
			options.Block.Format(req.Format),
			options.Block.Pin(req.Pin))
		if err != nil {
			return err
		}

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		bcid, err := settings.CidPrefix.Sum(data)
		if err != nil {
			return err
		}

		b, err := blocks.NewBlockWithCid(data, bcid)
		if err != nil {
			return err
		}

		err = s.bs.AddBlock(ctx, b)
		if err != nil {
			return err
		}

		// todo pin
		//if settings.Pin {
		//	if err = s.pinning.PinWithMode(ctx, b.Cid(), pin.Recursive); err != nil {
		//		return err
		//	}
		//	if err := s.pinning.Flush(ctx); err != nil {
		//		return err
		//	}
		//}

		if !req.AllowBigBlock && len(b.RawData()) > 1024*1024 {
			return fmt.Errorf("produced block is over 1MiB: big blocks can't be exchanged with other peers. consider using UnixFS for automatic chunking of bigger files, or pass --allow-big-block to override")
		}

		if err = res.Emit(&BlockStat{
			Key:  b.Cid().String(),
			Size: len(b.RawData()),
		}); err != nil {
			return err
		}
	}

	return it.Err()
}

type BlockStat struct {
	Key  string
	Size int
}

type BlockGetRequest struct {
	Arg string `json:"arg,omitempty"`
}

func (s *BlocksService) BlockGet(ctx httpctx.Context) (err error) {
	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})

	var req BlockGetRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}

	if req.Arg == "" {
		return errors.New("argument \"arg\" is required")
	}

	rp, err := s.dagResolver.ResolvePath(ctx, path.New(req.Arg))
	if err != nil {
		return err
	}

	b, err := s.bs.GetBlock(ctx, rp.Cid())
	if err != nil {
		return err
	}

	return res.Emit(bytes.NewReader(b.RawData()))
}
