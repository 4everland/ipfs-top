package service

import (
	"context"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-libipfs/bitswap"
	"github.com/ipfs/go-libipfs/bitswap/network"
	"sync"
)

var (
	bitSwapServiceOnce      = sync.Once{}
	bitSwapServiceWatchOnce = sync.Once{}

	bitSwapService *BitSwapService
)

type BitSwapService struct {
	bs blockstore.Blockstore

	bitswapimpl *bitswap.Bitswap
}

func NewBitSwapService(bstore blockstore.Blockstore) *BitSwapService {
	bitSwapServiceOnce.Do(func() {
		bitSwapService = &BitSwapService{
			bs: bstore,
		}
	})
	return bitSwapService
}

func (bss *BitSwapService) BitSwap() *bitswap.Bitswap {
	return bss.bitswapimpl
}

func (bss *BitSwapService) Watch(ctx context.Context, node NodeInterface) {
	bitSwapServiceWatchOnce.Do(func() {
		net := network.NewFromIpfsHost(node.GetHost(), node.GetContentRouting())
		bss.bitswapimpl = bitswap.New(ctx, net, bss.bs)
	})
}
