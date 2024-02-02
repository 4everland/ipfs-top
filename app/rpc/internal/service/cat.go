package service

import (
	"context"
	"fmt"
	"github.com/4everland/ipfs-top/third_party/coreunix"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	cmds "github.com/ipfs/go-ipfs-cmds"
	http2 "github.com/ipfs/go-ipfs-cmds/http"
	"io"
)

type CatService struct {
	unixfs *coreunix.UnixFsServer
}

func NewCatService(unixfs *coreunix.UnixFsServer) *CatService {
	return &CatService{unixfs}
}

func (s *CatService) RegisterRoute(route *httpctx.Router) {
	route.POST("/cat", s.Cat)
}

type CatRequest struct {
	Arg      []string `json:"arg,omitempty"`
	Offset   int64    `json:"offset,omitempty"`
	Length   int64    `json:"length,omitempty"`
	Progress bool     `json:"progress,omitempty"`
}

func (s *CatService) Cat(ctx httpctx.Context) (err error) {
	res, err := http2.NewResponseEmitter(ctx.Response(), ctx.Request().Method, &cmds.Request{
		Options: cmds.OptMap{cmds.EncLong: cmds.JSON},
		Context: ctx,
	})
	var req CatRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}
	if req.Offset < 0 {
		return fmt.Errorf("cannot specify negative offset")
	}

	if req.Length < 0 {
		return fmt.Errorf("cannot specify negative length")
	}

	if ctx.Query().Get("length") == "" {
		req.Length = -1
	}

	readers, length, err := cat(ctx, s.unixfs, req.Arg, req.Offset, req.Length)
	if err != nil {
		return err
	}

	res.SetLength(length)
	reader := io.MultiReader(readers...)

	// Since the reader returns the error that a block is missing, and that error is
	// returned from io.Copy inside Emit, we need to take Emit errors and send
	// them to the client. Usually we don't do that because it means the connection
	// is broken or we supplied an illegal argument etc.
	return res.Emit(reader)
}

func cat(ctx context.Context, api coreunix.UnixfsAPI, paths []string, offset int64, max int64) ([]io.Reader, uint64, error) {
	readers := make([]io.Reader, 0, len(paths))
	length := uint64(0)
	if max == 0 {
		return nil, 0, nil
	}
	for _, p := range paths {
		pp, err := path.NewPath(p)
		if err != nil {
			return nil, 0, err
		}
		f, err := api.Get(ctx, pp)
		if err != nil {
			return nil, 0, err
		}

		var file files.File
		switch f := f.(type) {
		case files.File:
			file = f
		case files.Directory:
			return nil, 0, coreunix.ErrIsDir
		default:
			return nil, 0, coreunix.ErrNotSupported
		}

		fsize, err := file.Size()
		if err != nil {
			return nil, 0, err
		}

		if offset > fsize {
			offset = offset - fsize
			continue
		}

		count, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, 0, err
		}
		offset = 0

		fsize, err = file.Size()
		if err != nil {
			return nil, 0, err
		}

		size := uint64(fsize - count)
		length += size
		if max > 0 && length >= uint64(max) {
			var r io.Reader = file
			if overshoot := int64(length - uint64(max)); overshoot != 0 {
				r = io.LimitReader(file, int64(size)-overshoot)
				length = uint64(max)
			}
			readers = append(readers, r)
			break
		}
		readers = append(readers, file)
	}
	return readers, length, nil
}
