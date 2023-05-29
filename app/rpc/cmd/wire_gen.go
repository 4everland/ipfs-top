// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/4everland/ipfs-servers/app/rpc/internal/conf"
	"github.com/4everland/ipfs-servers/app/rpc/internal/server"
	"github.com/4everland/ipfs-servers/app/rpc/internal/service"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init task receiver server application.
func wireApp(confServer *conf.Server, data *conf.Data, version *conf.Version, logger log.Logger) (*kratos.App, func(), error) {
	blockstore := service.NewBlockStore(data)
	exchangeInterface := service.NewExchange(data)
	unixFsServer := coreunix.NewUnixFsServer(blockstore, exchangeInterface)
	adderService := service.NewAdderService(unixFsServer)
	pinAPI := service.NewPinAPI(data)
	blockService := service.NewBlockService(blockstore, exchangeInterface)
	dagService := service.NewDAGService(blockService)
	dagResolve := service.NewDagResolve(dagService, blockService)
	pinService := service.NewPinService(pinAPI, dagResolve)
	lsService := service.NewLsService(unixFsServer)
	offlineBlockService := service.NewOfflineBlockService(blockstore)
	filesService := service.NewFilesService(dagService, offlineBlockService, dagResolve)
	catService := service.NewCatService(unixFsServer)
	versionInfo := service.NewVersionInfo(version)
	versionService := service.NewVersionService(versionInfo)
	serviceDagService := service.NewDagService(dagService)
	httpServer := server.NewApiHttpServer(confServer, adderService, pinService, lsService, filesService, catService, versionService, serviceDagService, logger)
	app := newApp(logger, httpServer)
	return app, func() {
	}, nil
}
