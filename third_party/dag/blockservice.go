package dag

import (
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	exchange "github.com/ipfs/boxo/exchange"
)

func newBlockService(store blockstore.Blockstore, ex exchange.Interface) blockservice.BlockService {
	return blockservice.New(store, ex)
}

func NewGrpcBlockService(endpoint, exchangeEndpoint string) (blockservice.BlockService, error) {
	store, err := NewBlockStore(endpoint)
	if err != nil {
		return nil, err
	}
	ex, err := NewGrpcRouting(exchangeEndpoint)
	if err != nil {
		return nil, err
	}
	return newBlockService(store, ex), nil
}
