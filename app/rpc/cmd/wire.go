// The build tag makes sure the stub is not built in the final build.
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/4everland/ipfs-servers/app/rpc/internal/conf"
	"github.com/4everland/ipfs-servers/app/rpc/internal/server"
	"github.com/4everland/ipfs-servers/app/rpc/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init task receiver server application.
func wireApp(*conf.Server, *conf.Data, *conf.Version, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, newApp))
}
