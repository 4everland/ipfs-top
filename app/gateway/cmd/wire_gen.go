// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/4everland/ipfs-top/app/gateway/internal/conf"
	"github.com/4everland/ipfs-top/app/gateway/internal/server"
	"github.com/4everland/ipfs-top/app/gateway/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init task receiver server application.
func wireApp(confServer *conf.Server, data *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	blockstore := service.NewBlockStore(data, logger)
	blocksBackend, err := service.NewBlocksGateway(blockstore)
	if err != nil {
		return nil, nil, err
	}
	v := server.NewServer(confServer, blocksBackend, logger)
	app := newApp(logger, v)
	return app, func() {
	}, nil
}
