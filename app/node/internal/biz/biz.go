package biz

import (
	"github.com/4everland/ipfs-servers/app/node/internal/biz/provide"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(provide.SimpleProvider, provide.SimpleReprovider, provide.ProviderQueue)
