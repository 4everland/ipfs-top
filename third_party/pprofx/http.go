package pprofx

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
)

func Server(servers ...transport.Server) kratos.Option {
	server := http.NewServer(http.Address(":0"))
	server.HandlePrefix("/", pprof.NewHandler())
	servers = append(servers, server)
	return kratos.Server(servers...)
}
