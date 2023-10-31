// The build tag makes sure the stub is not built in the final build.
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/4everland/ipfs-top/app/node/internal/biz"
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/app/node/internal/data"
	"github.com/4everland/ipfs-top/app/node/internal/server"
	"github.com/4everland/ipfs-top/app/node/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init task receiver server application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, biz.ProviderSet, data.ProviderSet, newApp))
}
