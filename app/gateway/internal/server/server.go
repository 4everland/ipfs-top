package server

import (
	"github.com/4everland/ipfs-top/app/gateway/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"
	"github.com/ipfs/boxo/gateway"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

func NewServer(c *conf.Server, gw *gateway.BlocksBackend, logger log.Logger) []transport.Server {
	return []transport.Server{NewGatewayServer(c, gw, logger), NewMetricsServer(c, logger)}
}
