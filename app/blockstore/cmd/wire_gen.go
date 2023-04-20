// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	biz2 "github.com/4everland/ipfs-servers/app/blockstore/internal/biz"
	"github.com/4everland/ipfs-servers/app/blockstore/internal/conf"
	"github.com/4everland/ipfs-servers/app/blockstore/internal/server"
	"github.com/4everland/ipfs-servers/app/blockstore/internal/services"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init task receiver server application.
func wireApp(confServer *conf.Server, data *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	blockStore := biz2.NewBackendStorage(data, logger)
	blockIndex := biz2.NewIndexStore(data)
	blockstoreService := services.NewBlockstoreService(blockStore, blockIndex)
	grpcServer := server.NewGRPCServer(confServer, blockstoreService, logger)
	app := newApp(logger, grpcServer)
	return app, func() {
	}, nil
}
