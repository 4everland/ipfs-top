package service

import (
	"github.com/google/wire"
)

func NewNodeServices(bts *BitSwapService, ps *RoutingService) []NodeService {
	return []NodeService{
		bts,
		ps,
	}
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewBitSwapService, NewNodeServices, NewRoutingService)
