package service

import (
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdderService, coreunix.NewUnixFsServerOffline, NewBlockStore)
