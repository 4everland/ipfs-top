package dag

import (
	"github.com/ipfs/go-blockservice"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
)

func newBlockService(store blockstore.Blockstore, ex exchange.Interface) blockservice.BlockService {
	return blockservice.New(store, ex)
}

func NewGrpcBlockService(endpoint, exchangeEndpoint string) (blockservice.BlockService, error) {
	store, err := NewBlockStore(endpoint)
	if err != nil {
		return nil, err
	}
	return newBlockService(store, nil), nil
}
