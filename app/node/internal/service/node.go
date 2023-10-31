package service

import (
	"context"
	"github.com/4everland/ipfs-top/app/node/internal/types"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

type NodeInterface interface {
	GetContentRouting() routing.Routing
	GetHost() host.Host
	GetConnMgr() connmgr.CMInfo
	Peers() []types.ConnectPeer
}

type NodeService interface {
	Watch(context.Context, NodeInterface)
}
