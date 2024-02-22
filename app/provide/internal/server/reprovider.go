package server

import (
	"context"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/provide/internal/data"
	"github.com/4everland/ipfs-top/app/provide/internal/helpers"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/ipld/merkledag"
	ipld "github.com/ipfs/go-ipld-format"
	"time"
)

type ReproviderServer struct {
	nodes []routing.RoutingClient
	bs    blockstore.Blockstore
	ng    ipld.NodeGetter

	pin    *data.PinSetRepo
	logger *log.Helper
}

func (server *ReproviderServer) Start(ctx context.Context) error {
	initial := make(chan struct{}, 1)
	ticker := time.NewTicker(time.Hour * 22)
	initial <- struct{}{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-initial:
			start := time.Now()
			server.logger.WithContext(ctx).Infof("initial start reprovide loop at: %s", start.Format("2006-01-02 15:04:05"))
			if err := server.reProvider(ctx); err != nil {
				server.logger.WithContext(ctx).Errorf("initial reprovide loop error: %v", err)
			}
			server.logger.WithContext(ctx).Infof("initial finish reprovide loop at: %s", time.Now().Format("2006-01-02 15:04:05"))
			initial = nil
		case <-ticker.C:
			start := time.Now()
			server.logger.WithContext(ctx).Infof("start reprovide loop at: %s", start.Format("2006-01-02 15:04:05"))
			if err := server.reProvider(ctx); err != nil {
				server.logger.WithContext(ctx).Errorf("reprovide loop error: %v", err)
			}
			server.logger.WithContext(ctx).Infof("finish reprovide loop at: %s", time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}

func (server *ReproviderServer) reProvider(ctx context.Context) error {
	ch, errCh := server.pin.AllKeys(ctx, time.Now())

	for {
		select {
		case cid, ok := <-ch:
			if !ok {
				return nil
			}
			n, err := server.ng.Get(ctx, cid)
			if err != nil {
				server.logger.WithContext(ctx).Errorf("reprovide %s error: %v", cid.String(), err)
				continue
			}

			iter := helpers.NewDagNodeIter(ctx, []ipld.Node{n}, server.ng)
			for {
				nn, err := iter.Next()
				if err != nil {
					server.logger.WithContext(ctx).Errorf("reprovide %s error: %v", nn.String(), err)
					break
				}
				if nn == nil {
					break
				}
				for _, node := range server.nodes {
					if _, err = node.Provide(ctx, &routing.ProvideReq{
						Cid:     &routing.Cid{Str: nn.Cid().Bytes()},
						Provide: true,
					}); err != nil {
						server.logger.WithContext(ctx).Errorf("reprovide %s error: %v", cid.String(), err)
					}
				}
			}
		case err := <-errCh:
			return err
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
	pin *data.PinSetRepo,
	logger log.Logger,
) *ReproviderServer {
	blockService := blockservice.New(bs, nil)
	dagService := merkledag.NewDAGService(blockService)
	return &ReproviderServer{
		nodes:  nodes,
		bs:     bs,
		ng:     dagService,
		pin:    pin,
		logger: log.NewHelper(logger),
	}
}
