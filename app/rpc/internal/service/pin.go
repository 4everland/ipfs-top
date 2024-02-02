package service

import (
	"context"
	"encoding/json"
	"github.com/4everland/ipfs-top/third_party/coreunix"
	"github.com/4everland/ipfs-top/third_party/coreunix/options"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/path"
	"net/http"
)

type PinService struct {
	pinning  coreunix.PinAPI
	resolver coreunix.DagResolve
}

func NewPinService(pinning coreunix.PinAPI, resolver coreunix.DagResolve) *PinService {
	return &PinService{pinning: pinning, resolver: resolver}
}

func (s *PinService) RegisterRoute(route *httpctx.Router) {
	route.POST("/pin/add", s.Add)
}

type PinAddRequest struct {
	Arg       []string `json:"arg,omitempty"`
	Progress  bool     `json:"progress,omitempty"`
	Recursive bool     `json:"recursive,omitempty"`
}

func (s *PinService) Add(ctx httpctx.Context) (err error) {
	w := ctx.Response()

	var req PinAddRequest
	if err = ctx.BindQuery(&req); err != nil {
		return
	}

	if ctx.Request().URL.Query().Get("recursive") == "" {
		req.Recursive = true
	}

	// todo req.Progress
	added, err := s.addMany(ctx, req.Arg, req.Recursive)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(PinAddResponse{added})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

	return nil
}

type PinAddResponse struct {
	Pins []string
}

func (s *PinService) addMany(ctx context.Context, paths []string, recursive bool) ([]string, error) {
	added := make([]string, len(paths))
	for i, b := range paths {
		p, err := path.NewPath(b)
		if err != nil {
			return nil, err
		}
		rp, err := s.resolver.ResolvePath(ctx, p)

		if err := s.pinning.Add(ctx, rp, options.Pin.Recursive(recursive)); err != nil {
			return nil, err
		}
		added[i] = rp.RootCid().String()
	}

	return added, nil
}
