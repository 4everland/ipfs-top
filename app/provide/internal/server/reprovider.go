package server

import (
	"context"
	"github.com/4everland/ipfs-top/app/provide/internal/biz"
)

type ReproviderServer struct {
	r *biz.ReProviderBiz
}

func (server *ReproviderServer) Start(ctx context.Context) error {
	return server.r.Run(ctx)
}

func (server *ReproviderServer) Stop(context.Context) error {
	return nil
}

func NewReproviderServer(
	r *biz.ReProviderBiz,
) *ReproviderServer {

	return &ReproviderServer{
		r: r,
	}
}
