package server

import (
	"context"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/blockstore"
	"time"
)

type ReproviderServer struct {
	nodes []routing.RoutingClient
	bs    blockstore.Blockstore
}

func (server *ReproviderServer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.NewTicker(time.Hour * 22).C:
			if err := server.reProvider(ctx); err != nil {
				log.NewHelper(log.DefaultLogger).WithContext(ctx).Errorf("reprovide error: %v", err)
			}
		}
	}
}

func (server *ReproviderServer) reProvider(ctx context.Context) error {
	ch, err := server.bs.AllKeysChan(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case cid, ok := <-ch:
			if !ok {
				return nil
			}

			for _, node := range server.nodes {
				if _, err = node.Provide(ctx, &routing.ProvideReq{
					Cid:     &routing.Cid{Str: cid.Bytes()},
					Provide: true,
				}); err != nil {
					log.NewHelper(log.DefaultLogger).WithContext(ctx).Errorf("reprovide %s error: %v", cid.String(), err)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (server *ReproviderServer) Stop(context.Context) error {
	return nil
}

func NewReproviderServer(
	bs blockstore.Blockstore,
	nodes []routing.RoutingClient,
) *ReproviderServer {
	return &ReproviderServer{
		nodes: nodes,
		bs:    bs,
	}
}
