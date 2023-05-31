package service

import (
	"context"
	"github.com/4everland/ipfs-servers/app/rpc/internal/conf"
	"github.com/4everland/ipfs-servers/third_party/coreapi"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/google/wire"
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	iface "github.com/ipfs/boxo/coreiface"
	exchange "github.com/ipfs/boxo/exchange"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	format "github.com/ipfs/go-ipld-format"
	"runtime"
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

	NewVersionInfo,

	NewAdderService,
	NewPinService,
	NewLsService,
	NewFilesService,
	NewCatService,
	NewVersionService,
	NewDagService,
	NewBlocksService,
)

type OfflineBlockService interface {
	blockservice.BlockService
}

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri, "")
	if err != nil {
		panic(err)
	}
	return s
}

func NewExchange(config *conf.Data) exchange.Interface {
	if config.GetExchangeEndpoint() == "" {
		return nil
	}
	s, err := dag.NewGrpcRouting(config.ExchangeEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func NewPinAPI(config *conf.Data) iface.PinAPI {
	if config.GetPinEndpoint() == "" {
		return nil
	}
	s, err := coreapi.NewPinAPI(config.PinEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func NewBlockService(blockStore blockstore.Blockstore, exchange exchange.Interface) blockservice.BlockService {
	return blockservice.New(blockStore, exchange)
}

func NewOfflineBlockService(blockStore blockstore.Blockstore) OfflineBlockService {
	return blockservice.New(blockStore, offline.Exchange(blockStore))
}

func NewDAGService(bs blockservice.BlockService) format.DAGService {
	return merkledag.NewDAGService(bs)
}

func NewDagResolve(dagService format.DAGService, bs blockservice.BlockService) coreunix.DagResolve {
	return coreunix.NewDagResolver(context.Background(), dagService, bs)
}

func NewVersionInfo(v *conf.Version) *VersionInfo {
	return &VersionInfo{
		Version: v.Version,
		Commit:  v.Commit,
		Repo:    v.Repo,
		System:  runtime.GOARCH + "/" + runtime.GOOS,
		Golang:  runtime.Version(),
	}
}
