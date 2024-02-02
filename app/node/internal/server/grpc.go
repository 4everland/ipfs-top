package server

import (
	"github.com/4everland/ipfs-top/api/routing"
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/app/node/internal/service"
	"github.com/4everland/ipfs-top/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewContentRoutingGRPCServer new a gRPC server.
func NewContentRoutingGRPCServer(c *conf.Server, routingSvc *service.RoutingService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(metadata.WithConstants(md.Metadata{
				enum.MetadataServerKind: []string{enum.ServerKindGRPC},
			})),
			logging.Server(logger),
			tracing.Server(),
		),
	}

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	routing.RegisterRoutingServer(srv, routingSvc)
	return srv
}
