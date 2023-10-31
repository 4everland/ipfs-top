// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/4everland/ipfs-top/app/node/internal/biz"
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/app/node/internal/data"
	"github.com/4everland/ipfs-top/app/node/internal/server"
	"github.com/4everland/ipfs-top/app/node/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init task receiver server application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	batching, err := data.NewLevelDbDatastore(confServer)
	if err != nil {
		return nil, nil, err
	}
	blockstore := data.NewBlockStore(confData)
	bitSwapService := service.NewBitSwapService(blockstore)
	v := biz.ProviderSystem(batching, blockstore)
	routingService := service.NewRoutingService(bitSwapService, v)
	v2 := service.NewNodeServices(bitSwapService, routingService)
	nodeServer, err := server.NewNodeServer(confServer, logger, batching, v2...)
	if err != nil {
		return nil, nil, err
	}
	grpcServer := server.NewContentRoutingGRPCServer(confServer, routingService, logger)
	httpServer := server.NewHttpServer(confServer, nodeServer, logger)
	app := newApp(logger, nodeServer, grpcServer, httpServer)
	return app, func() {
	}, nil
}
