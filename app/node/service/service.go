package service

import (
	"github.com/google/wire"
)

func NewNodeServices(bts *BitSwapService, ps *ProvideService) []NodeService {
	return []NodeService{
		bts,
		NewContentRoutingService(),
		ps,
	}
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewContentRoutingService, NewBitSwapService, NewNodeServices, NewSimpleProvideService)
