package service

import (
	"context"
	"github.com/ipfs/boxo/bitswap"
	"github.com/ipfs/boxo/bitswap/network"
	blockstore "github.com/ipfs/boxo/blockstore"
	metri "github.com/ipfs/go-metrics-interface"
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
		bsctx := metri.CtxScope(ctx, "node")

		bss.bitswapimpl = bitswap.New(bsctx, net, bss.bs)
		net.Start(bss.bitswapimpl)
	})
}
