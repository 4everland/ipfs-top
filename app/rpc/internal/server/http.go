package server

import (
	"github.com/4everland/ipfs-top/app/rpc/internal/conf"
	service2 "github.com/4everland/ipfs-top/app/rpc/internal/service"
	"github.com/4everland/ipfs-top/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewApiHttpServer(
	c *conf.Server,
	s *service2.AdderService,
	ps *service2.PinService,
	ls *service2.LsService,
	fs *service2.FilesService,
	cs *service2.CatService,
	vs *service2.VersionService,
	ds *service2.DagService,
	bs *service2.BlocksService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(),
			metadata.Server(metadata.WithConstants(md.Metadata{
				enum.MetadataServerKind: enum.ServerKindHTTP,
			})),
		),
		http.ErrorEncoder(errorEncoder),
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

	r := srv.Route("/api/v0")
	s.RegisterRoute(r)
	//ps.RegisterRoute(r)
	ls.RegisterRoute(r)
	fs.RegisterRoute(r)
	cs.RegisterRoute(r)
	vs.RegisterRoute(r)
	ds.RegisterRoute(r)
	bs.RegisterRoute(r)

	return srv
}
