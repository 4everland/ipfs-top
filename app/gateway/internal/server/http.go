package server

import (
	"github.com/4everland/ipfs-servers/app/gateway/internal/conf"
	"github.com/4everland/ipfs-servers/enum"
	"github.com/go-kratos/kratos/v2/log"
	md "github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/boxo/gateway"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewGatewayServer(c *conf.Server, gw *gateway.BlocksBackend, logger log.Logger) *http.Server {
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
	headers := map[string][]string{}
	gateway.AddAccessControlHeaders(headers)

	gwConf := gateway.Config{
		Headers:               headers,
		DeserializedResponses: true,
		PublicGateways: map[string]*gateway.PublicGateway{
			// Support public requests with Host: CID.ipfs.example.net and ID.ipns.example.net
			"example.net": {
				Paths:                 []string{"/ipfs"}, //[]string{"/ipfs", "/ipns"},
				NoDNSLink:             true,
				UseSubdomains:         true,
				DeserializedResponses: true,
			},
			// Support local requests
			"localhost": {
				Paths:         []string{"/ipfs"}, //[]string{"/ipfs", "/ipns"},
				NoDNSLink:     true,
				UseSubdomains: true,
			},
		},
	}

	gwHandler := gateway.NewHandler(gwConf, gw)
	handler := gateway.NewHostnameHandler(gwConf, gw, gwHandler)
	srv.HandlePrefix("/ipfs/", handler)
	srv.HandlePrefix("/", handler)
	srv.Handle("/metrics", promhttp.Handler())

	//srv.HandlePrefix("/ipns/", gwHandler)
	srv.Route("/ping").GET("/", func(ctx http.Context) error {
		//Hello from IPFS Gateway Checker
		//bafybeifx7yeb55armcsxwwitkymga5xf53dxiarykms3ygqic223w5sk3m
		return ctx.String(200, "pong")
	})

	return srv
}
