// The build tag makes sure the stub is not built in the final build.
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/4everland/ipfs-servers/app/blockstore/biz"
	"github.com/4everland/ipfs-servers/app/blockstore/conf"
	"github.com/4everland/ipfs-servers/app/blockstore/server"
	"github.com/4everland/ipfs-servers/app/blockstore/services"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init task receiver server application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, services.ProviderSet, biz.ProviderSet, newApp))
}
