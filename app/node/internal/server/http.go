package server

import (
	"github.com/4everland/ipfs-servers/app/node/internal/conf"
	"github.com/4everland/ipfs-servers/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHttpServer(c *conf.Server, s *NodeServer, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			metadata.Server(metadata.WithConstants(md.Metadata{
				enum.MetadataServerKind: enum.ServerKindHTTP,
			})),
		),
	}

	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.Route("/ping").GET("/", func(ctx http.Context) error {
		return ctx.String(200, "pong")
	})
	srv.Handle("/metrics", promhttp.Handler())
	stats := srv.Route("/stats")
	stats.GET("/peers", func(ctx http.Context) error {
		return ctx.JSON(200, s.Peers())
	})
	stats.GET("/conn", func(ctx http.Context) error {
		return ctx.JSON(200, s.GetConnMgr())
	})
	return srv
}
