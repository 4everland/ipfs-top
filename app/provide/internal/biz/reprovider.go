package biz

import (
	"context"
	"fmt"
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/provide/internal/data"
	"github.com/4everland/ipfs-top/app/provide/internal/helpers"
	"github.com/4everland/ipfs-top/third_party/prom"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/panjf2000/ants/v2"
	"time"
)

type ReProviderBiz struct {
	nodes []*data.NamedRoutingClient
	bs    blockstore.Blockstore
	ng    ipld.NodeGetter

	pin    *data.PinSetRepo
	logger *log.Helper

	metrics *prom.ProvideMetrics

	current string
}

func (server *ReProviderBiz) Current() string {
	return server.current
}

func (server *ReProviderBiz) Run(ctx context.Context) error {
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

func (server *ReProviderBiz) reProvider(ctx context.Context) error {
	ch, errCh := server.pin.AllKeys(ctx, time.Now())
	server.metrics.Clean()
	defer func() {
		server.metrics.Clean()
	}()

	if len(server.nodes) == 0 {
		return fmt.Errorf("at least one node is required to provide data")
	}

	product := make(chan cid.Cid, len(server.nodes))
	defer close(product)
	for _, node := range server.nodes {
		go func(ctx context.Context, node *data.NamedRoutingClient) {
			errCount := 0
			for c := range product {
				_, err := node.Client.Provide(ctx, &routing.ProvideReq{
					Cid:     &routing.Cid{Str: c.Bytes()},
					Provide: true,
				})
				if err != nil {
					errCount++
					server.logger.WithContext(ctx).Errorf("provide %s error: %v", c.String(), err)
					product <- c
				} else {
					server.metrics.Provide(node.Name)
					continue
				}
				if errCount >= 10 {
					time.Sleep(time.Second * 20)
				}
				errCount = 0
			}
		}(ctx, node)
	}

	p, err := ants.NewPoolWithFunc(len(server.nodes), func(i interface{}) {
		cc, ok := i.(cid.Cid)
		if !ok {
			return
		}
		n, err := server.ng.Get(ctx, cc)
		if err != nil {
			server.logger.WithContext(ctx).Errorf("reprovide %s error: %v", cc.String(), err)
			return
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
			product <- nn.Cid()
		}
		server.metrics.Inc()
	}, ants.WithPanicHandler(func(i interface{}) {
		server.logger.Errorf("reprovide panic: %v", i)
	}))

	if err != nil {
		return err
	}

	defer p.Release()
	for {
		select {
		case cc, ok := <-ch:
			if !ok {
				return nil
			}
			server.current = cc.String()
			if cc.Prefix().Codec == cid.Raw {
				product <- cc
				server.metrics.Inc()
				continue
			}
			p.Invoke(cc)
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func NewReProviderBiz(
	bs blockstore.Blockstore,
	nodes []*data.NamedRoutingClient,
	pin *data.PinSetRepo,
	logger log.Logger,
) *ReProviderBiz {
	blockService := blockservice.New(bs, nil)
	dagService := merkledag.NewDAGService(blockService)
	return &ReProviderBiz{
		nodes:   nodes,
		bs:      bs,
		ng:      dagService,
		pin:     pin,
		logger:  log.NewHelper(logger),
		metrics: prom.NewProvideMetrics("ip"),
	}
}
