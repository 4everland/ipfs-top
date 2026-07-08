package server

import (
	"github.com/4everland/ipfs-top/app/gateway/internal/conf"
	"github.com/4everland/ipfs-top/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsServer(c *conf.Server, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			metadata.Server(metadata.WithConstants(md.Metadata{
				enum.MetadataServerKind: []string{enum.ServerKindHTTP},
			})),
		),
	}

	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Metrics.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Metrics.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.Handle("/metrics", promhttp.Handler())
	srv.Route("/ping").GET("/", func(ctx http.Context) error {
		//Hello from IPFS Gateway Checker
		//bafybeifx7yeb55armcsxwwitkymga5xf53dxiarykms3ygqic223w5sk3m
		return ctx.String(200, "pong")
	})

	return srv
}
