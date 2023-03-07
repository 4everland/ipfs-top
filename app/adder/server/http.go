package server

import (
	"github.com/4everland/ipfs-servers/app/adder/conf"
	"github.com/4everland/ipfs-servers/app/adder/service"
	"github.com/4everland/ipfs-servers/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewApiHttpServer(c *conf.Server, s *service.AdderService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			metadata.Server(metadata.WithConstants(md.Metadata{
				enum.MetadataServerKind: enum.ServerKindHTTP,
			})),
		),
		//http.RequestDecoder(middleware.TransformAdderRequest),
		//http.ErrorEncoder(middleware.TransformAdderErrorResponse),
	}

	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	srv.Route("/").GET("ping", func(ctx http.Context) error {
		return ctx.String(200, "pong")
	})

	srv.Route("/").POST("/api/v0/add", s.Add)
	return srv
}
