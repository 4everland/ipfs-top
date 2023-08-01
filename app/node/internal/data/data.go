package data

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewBlockStore, NewLevelDbDatastore, NewLruMapDatastore)
