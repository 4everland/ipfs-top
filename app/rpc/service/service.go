package service

import (
	"context"
	"github.com/4everland/ipfs-servers/app/rpc/conf"
	"github.com/4everland/ipfs-servers/third_party/coreapi"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/google/wire"
	"github.com/ipfs/go-blockservice"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewBlockStore,
	NewExchange,
	NewBlockService,
	NewOfflineBlockService,
	NewDAGService,
	NewDagResolve,
	NewPinAPI,

	coreunix.NewUnixFsServer,

	NewAdderService,
	NewPinService,
	NewLsService,
	NewFilesService,
	NewCatService,
)

type offlineBlockService interface {
	blockservice.BlockService
}

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}

func NewExchange(config *conf.Data) exchange.Interface {
	s, err := dag.NewGrpcRouting(config.ExchangeEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func NewPinAPI(config *conf.Data) iface.PinAPI {
	s, err := coreapi.NewPinAPI(config.PinEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func NewBlockService(blockStore blockstore.Blockstore, exchange exchange.Interface) blockservice.BlockService {
	return blockservice.New(blockStore, exchange)
}

func NewOfflineBlockService(blockStore blockstore.Blockstore) offlineBlockService {
	return blockservice.New(blockStore, offline.Exchange(blockStore))
}

func NewDAGService(bs blockservice.BlockService) format.DAGService {
	return merkledag.NewDAGService(bs)
}

func NewDagResolve(dagService format.DAGService, bs blockservice.BlockService) coreunix.DagResolve {
	return coreunix.NewDagResolver(context.Background(), dagService, bs)
}
