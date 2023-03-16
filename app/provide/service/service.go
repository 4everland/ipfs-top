package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(ProviderQueue, SimpleProvider, SimpleReprovider, NewProvideService)
