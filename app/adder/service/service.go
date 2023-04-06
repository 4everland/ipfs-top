package service

import (
	"github.com/4everland/ipfs-servers/app/adder/conf"
	"github.com/4everland/ipfs-servers/third_party/coreapi"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/4everland/ipfs-servers/third_party/dag"
	"github.com/google/wire"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewBlockStore,
	NewExchange,
	NewPinAPI,

	coreunix.NewUnixFsServer,

	NewAdderService,
	NewPinService,
)

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
