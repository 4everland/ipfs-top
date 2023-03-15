package service

import (
	"context"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-libipfs/bitswap"
	"github.com/ipfs/go-libipfs/bitswap/network"
	"sync"
)

var (
	bitSwapServiceOnce = sync.Once{}
)

type BitSwapService struct {
	bs blockstore.Blockstore

	bitswapimpl *bitswap.Bitswap
}

func NewBitSwapService(bstore blockstore.Blockstore) *BitSwapService {
	return &BitSwapService{
		bs: bstore,
	}
}

func (bss *BitSwapService) Watch(ctx context.Context, node NodeInterface) {
	bitSwapServiceOnce.Do(func() {
		net := network.NewFromIpfsHost(node.GetHost(), node.GetContentRouting())
		bss.bitswapimpl = bitswap.New(ctx, net, bss.bs)
	})
}
