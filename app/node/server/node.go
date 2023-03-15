package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/4everland/ipfs-servers/app/node/conf"
	"github.com/4everland/ipfs-servers/app/node/data"
	"github.com/4everland/ipfs-servers/app/node/service"
	"github.com/4everland/ipfs-servers/app/node/types"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/go-datastore"
	dsleveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"time"
)

type NodeServer struct {
	leveldbpath string
	addrs       []string

	priKey crypto.PrivKey

	ps    peerstore.Peerstore
	dhtDs datastore.Batching
	pmDs  datastore.Batching
	peers []peer.AddrInfo

	connManger *connmgr.BasicConnMgr
	h          host.Host
	rt         routing.Routing
	nat        autonat.AutoNAT

	logger *log.Helper

	services []service.NodeService
}

func NewNodeServer(serverConf *conf.Server, logger log.Logger, svcs ...service.NodeService) (*NodeServer, error) {
	connManger, err := connmgr.NewConnManager(
		int(serverConf.Node.LowWater), int(serverConf.Node.HighWater),
		connmgr.WithGracePeriod(time.Duration(serverConf.Node.GracePeriod)*time.Second))
	if err != nil {
		return nil, err
	}

	ps, err := pstoremem.NewPeerstore()
	if err != nil {
		return nil, err
	}

	pmDs, err := data.NewLruMapDatastore()
	if err != nil {
		return nil, err
	}

	b, err := base64.StdEncoding.DecodeString(serverConf.Node.PrivateKey)
	if err != nil {
		return nil, err
	}

	priKey, err := crypto.UnmarshalPrivateKey(b)
	if err != nil {
		return nil, err
	}

	peers, err := getPeerAddrs(serverConf.Node.GetPeers())
	if err != nil {
		return nil, err
	}

	addrs, err := peer.AddrInfosFromP2pAddrs(peers...)
	if err != nil {
		return nil, err
	}

	return &NodeServer{
		leveldbpath: serverConf.Node.LeveldbPath,
		addrs:       serverConf.Node.MultiAddr,
		priKey:      priKey,
		ps:          ps,
		pmDs:        pmDs,
		connManger:  connManger,
		peers:       addrs,
		logger:      log.NewHelper(logger),
		services:    svcs,
	}, nil
}

func (server *NodeServer) Start(ctx context.Context) (err error) {
	server.dhtDs, err = dsleveldb.NewDatastore(server.leveldbpath, &dsleveldb.Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if server.h, err = libp2p.New(
		libp2p.ListenAddrStrings(server.addrs...),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			pms, er := providers.NewProviderManager(ctx, h.ID(), server.ps, server.pmDs)
			if er != nil {
				return nil, er
			}
			nodeDht, er := dual.New(
				ctx, h,
				dual.DHTOption(dht.ProtocolPrefix(dht.DefaultPrefix)),
				dual.DHTOption(
					dht.RoutingTableLatencyTolerance(10*time.Second),
					dht.Concurrency(3),
					dht.Mode(dht.ModeAutoServer),
					dht.Datastore(server.dhtDs),
					dht.ProviderStore(pms),
					dht.Validator(record.NamespacedValidator{
						"pk":   record.PublicKeyValidator{},
						"ipns": ipns.Validator{KeyBook: server.ps},
					}),
				),
				dual.WanDHTOption(dht.BootstrapPeers(server.peers...)),
			)
			if er != nil {
				return nil, er
			}
			server.rt = nodeDht
			return nodeDht, nil
		}),
		libp2p.AutoNATServiceRateLimit(600, 5, time.Minute),
		libp2p.EnableRelayService(),
		libp2p.EnableHolePunching(),
		libp2p.Identity(server.priKey),
		libp2p.EnableNATService(),
		libp2p.Peerstore(server.ps),
		libp2p.ConnectionManager(server.connManger),
		libp2p.MultiaddrResolver(madns.DefaultResolver),
	); err != nil {
		return err
	}
	dialback, err := libp2p.New(libp2p.NoListenAddrs)
	n, err := autonat.New(server.h, autonat.EnableService(dialback.Network()))
	if err != nil {
		return err
	}
	server.nat = n

	server.logger.Infof("DHT node started.")

	for _, addr := range server.h.Addrs() {
		server.logger.Infof("addr: %s/p2p/%s", addr.String(), server.h.ID())
	}

	for _, s := range server.services {
		s.Watch(ctx, server)
	}
	return nil
}

func (server *NodeServer) Stop(ctx context.Context) (err error) {
	if server.h != nil {
		err = server.h.Close()
	}
	return err
}

func (server *NodeServer) GetConnMgr() connmgr.CMInfo {
	if server.connManger == nil {
		return connmgr.CMInfo{}
	}
	return server.connManger.GetInfo()
}

func (server *NodeServer) Peers() []types.ConnectPeer {
	if server.h == nil {
		return nil
	}
	conns := server.h.Network().Conns()

	out := make([]types.ConnectPeer, 0, len(conns))
	for _, c := range conns {

		ci := types.ConnectPeer{
			Id:      c.RemotePeer().String(),
			Local:   c.LocalMultiaddr().String(),
			Addr:    c.RemoteMultiaddr().String(),
			Network: c.Stat().Stats,
		}
		out = append(out, ci)
	}

	return out
}

func (server *NodeServer) NatState() autonat.AutoNAT {
	if server.h == nil {
		return nil
	}
	return server.nat

}

func (server *NodeServer) GetContentRouting() routing.Routing {
	return server.rt
}

func (server *NodeServer) GetHost() host.Host {
	return server.h
}

func (server *NodeServer) ConnectCount() connmgr.CMInfo {
	if server.connManger == nil {
		return connmgr.CMInfo{}
	}
	return server.connManger.GetInfo()
}

func (server *NodeServer) PrintNode() {
	server.logger.Infof("conn count: %d peer count:%d",
		server.ConnectCount().ConnCount, len(server.Peers()))
}

func getPeerAddrs(addrs []string) ([]ma.Multiaddr, error) {
	var maddrs []ma.Multiaddr
	for _, s := range addrs {
		a, err := ma.NewMultiaddr(s)
		if err != nil {
			return nil, err
		}
		maddrs = append(maddrs, a)
	}

	return maddrs, nil
}
