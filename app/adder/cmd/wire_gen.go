// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/4everland/ipfs-servers/app/adder/conf"
	"github.com/4everland/ipfs-servers/app/adder/server"
	"github.com/4everland/ipfs-servers/app/adder/service"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init task receiver server application.
func wireApp(confServer *conf.Server, data *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	blockstore := service.NewBlockStore(data)
	exchangeInterface := service.NewExchange(data)
	unixFsServer := coreunix.NewUnixFsServer(blockstore, exchangeInterface)
	adderService := service.NewAdderService(unixFsServer)
	httpServer := server.NewApiHttpServer(confServer, adderService, logger)
	app := newApp(logger, httpServer)
	return app, func() {
	}, nil
}
