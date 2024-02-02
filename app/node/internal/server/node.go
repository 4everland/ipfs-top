package server

import (
	"context"
	"encoding/base64"
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/app/node/internal/service"
	"github.com/4everland/ipfs-top/app/node/internal/types"
	rcmgr2 "github.com/4everland/ipfs-top/third_party/rcmgr"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/bootstrap"
	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/peering"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
	"io"
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
	logger     *log.Helper

	services []service.NodeService

	Bootstrapper io.Closer
}

type RoutingOptionArgs struct {
	Ctx                           context.Context
	Host                          host.Host
	Datastore                     datastore.Batching
	Validator                     record.Validator
	BootstrapPeers                []peer.AddrInfo
	OptimisticProvide             bool
	OptimisticProvideJobsPoolSize int
}

func NewNodeServer(serverConf *conf.Server, logger log.Logger, ds datastore.Batching, svcs ...service.NodeService) (*NodeServer, error) {
	connManger, err := connmgr.NewConnManager(
		int(serverConf.Node.LowWater), int(serverConf.Node.HighWater),
		connmgr.WithGracePeriod(time.Duration(serverConf.Node.GracePeriod)*time.Second))
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
		dhtDs:       ds,
		leveldbpath: serverConf.Node.LeveldbPath,
		addrs:       serverConf.Node.MultiAddr,
		priKey:      priKey,
		connManger:  connManger,
		peers:       addrs,
		logger:      log.NewHelper(logger),
		services:    svcs,
	}, nil
}

func RoutingOption(mode dht.ModeOpt, args RoutingOptionArgs) (routing.Routing, error) {
	dhtOpts := []dht.Option{
		dht.Concurrency(10),
		dht.Mode(mode),
		dht.Datastore(args.Datastore),
		dht.Validator(args.Validator),
	}
	if args.OptimisticProvide {
		dhtOpts = append(dhtOpts, dht.EnableOptimisticProvide())
	}
	if args.OptimisticProvideJobsPoolSize != 0 {
		dhtOpts = append(dhtOpts, dht.OptimisticProvideJobsPoolSize(args.OptimisticProvideJobsPoolSize))
	}
	return dual.New(
		args.Ctx, args.Host,
		dual.DHTOption(dhtOpts...),
		dual.WanDHTOption(dht.BootstrapPeers(args.BootstrapPeers...)),
	)
}

func (server *NodeServer) Start(ctx context.Context) (err error) {
	bwc := metrics.NewBandwidthCounter()
	limiter := rcmgr.NewFixedLimiter(rcmgr2.MakeResourceManagerConfig(2*1024*1024*1024, 1024, server.connManger.GetInfo().HighWater))
	mgr, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(server.addrs...),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(server.connManger), //todo
		libp2p.Identity(server.priKey),
		libp2p.BandwidthReporter(bwc),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.ResourceManager(mgr),
	}

	opts = append(opts, libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		args := RoutingOptionArgs{
			Ctx:       ctx,
			Datastore: server.dhtDs,
			Validator: record.NamespacedValidator{
				"pk":   record.PublicKeyValidator{},
				"ipns": ipns.Validator{KeyBook: h.Peerstore()},
			},
			BootstrapPeers:                server.peers,
			OptimisticProvide:             true,
			OptimisticProvideJobsPoolSize: 100,
		}
		args.Host = h
		r, err := RoutingOption(dht.ModeServer, args)
		server.rt = r
		return r, err

	}))

	server.h, err = libp2p.New(opts...)
	if err != nil {
		return err
	}

	err = server.Bootstrap(bootstrap.DefaultBootstrapConfig)
	if err != nil {
		return err
	}

	p := peering.NewPeeringService(server.h)
	for _, info := range server.peers {
		p.AddPeer(info)
	}
	if err = p.Start(); err != nil {
		server.logger.Errorf("Start ppering service error :%s", err)
	}
	server.logger.Infof("DHT node started.")

	for _, addr := range server.h.Addrs() {
		server.logger.Infof("addr: %s/p2p/%s", addr.String(), server.h.ID())
	}

	for _, s := range server.services {
		s.Watch(ctx, server)
	}
	return nil
}

func (server *NodeServer) Bootstrap(cfg bootstrap.BootstrapConfig) (err error) {
	peerID, err := peer.IDFromPublicKey(server.priKey.GetPublic())
	if err != nil {
		return err
	}
	//n.Identity, n.PeerHost, n.Routing, cfg
	if server.rt == nil {
		return nil
	}

	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []peer.AddrInfo {
			return server.peers
		}
	}
	if server.Bootstrapper != nil {
		server.Bootstrapper.Close() // stop previous bootstrap process.
	}
	server.Bootstrapper, err = bootstrap.Bootstrap(peerID, server.h, server.rt, cfg)
	return err
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
			Id:        c.RemotePeer().String(),
			Local:     c.LocalMultiaddr().String(),
			Addr:      c.RemoteMultiaddr().String(),
			Opened:    c.Stat().Opened,
			Direction: c.Stat().Direction.String(),
			Transient: c.Stat().Transient,
		}
		out = append(out, ci)
	}

	return out
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
