package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/4everland/ipfs-top/app/node/internal/conf"
	"github.com/4everland/ipfs-top/app/node/internal/service"
	"github.com/4everland/ipfs-top/app/node/internal/types"
	rcmgr2 "github.com/4everland/ipfs-top/third_party/rcmgr"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
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
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	webrtc "github.com/libp2p/go-libp2p/p2p/transport/webrtc"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	webtransport "github.com/libp2p/go-libp2p/p2p/transport/webtransport"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	htp "net/http"
	"os"
	"time"
)

type NodeServer struct {
	leveldbpath string
	addrs       []string

	maxMemory uint64
	maxFd     int

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

	ids identify.IDService

	ping *ping.PingService
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

type AddrInfo struct {
	Addr string `json:"addr"`
}

type Peer struct {
	Id string `json:"id"`
}

func NewNodeServer(serverConf *conf.Server, logger log.Logger, ds datastore.Batching, svcs ...service.NodeService) (*NodeServer, error) {
	connManger, err := connmgr.NewConnManager(
		int(serverConf.Node.LowWater), int(serverConf.Node.HighWater),
		connmgr.WithGracePeriod(time.Duration(serverConf.Node.GracePeriod)*time.Second))
	if err != nil {
		return nil, err
	}

	l := log.NewHelper(logger)
	var sk crypto.PrivKey
	privKey := serverConf.Node.PrivateKey
	if privKey == "" {
		keyPath := serverConf.Node.PrivateKeyPath
		if keyPath != "" {
			//load private key from file
			skbytes, err := os.ReadFile(keyPath)
			if err != nil {
				return nil, err
			}
			privKey = string(skbytes)
		}
	}
	if privKey == "" {
		sk, _, err = crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}
		skbytes, err := crypto.MarshalPrivateKey(sk)
		if err != nil {
			return nil, err
		}
		//save private key
		privKey = base64.StdEncoding.EncodeToString(skbytes)
		l.Infof("init private key: %s", privKey)
		if serverConf.Node.PrivateKeyPath != "" {
			err := os.WriteFile(serverConf.Node.PrivateKeyPath, []byte(privKey), 0644)
			if err != nil {
				l.Errorf("save private key to file error: %s", err)
			}
		}
	} else {
		b, err := base64.StdEncoding.DecodeString(privKey)
		if err != nil {
			return nil, err
		}

		priKey, err := crypto.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
		sk = priKey
	}

	peers, err := getPeerAddrs(serverConf.Node.GetPeers())
	if err != nil {
		return nil, err
	}

	addrs, err := peer.AddrInfosFromP2pAddrs(peers...)
	if err != nil {
		return nil, err
	}

	maxMemory := serverConf.Node.MaxMemory
	maxFd := serverConf.Node.MaxFd
	if maxMemory == 0 {
		maxMemory = 2 * 1024 * 1024 * 1024 // 2GB
	}
	if maxFd == 0 {
		maxFd = 1024
	}
	return &NodeServer{
		dhtDs:       ds,
		leveldbpath: serverConf.Node.LeveldbPath,
		addrs:       serverConf.Node.MultiAddr,
		priKey:      sk,
		connManger:  connManger,
		peers:       addrs,
		logger:      l,
		services:    svcs,
		maxFd:       int(maxFd),
		maxMemory:   maxMemory,
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
	limiter := rcmgr.NewFixedLimiter(rcmgr2.MakeResourceManagerConfig(server.maxMemory, server.maxFd, server.connManger.GetInfo().HighWater))
	mgr, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(server.addrs...),
		//libp2p.NATPortMap(),
		libp2p.ConnectionManager(server.connManger), //todo
		libp2p.Identity(server.priKey),
		libp2p.BandwidthReporter(bwc),
		libp2p.Transport(tcp.NewTCPTransport, tcp.WithMetrics()),
		libp2p.Transport(websocket.New),
		libp2p.Transport(quic.NewTransport),
		libp2p.Transport(webtransport.New),
		libp2p.Transport(webrtc.New),
		libp2p.DefaultMuxers,
		libp2p.ResourceManager(mgr),
		//libp2p.UserAgent("kubo/0.27.0"),
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
		server.logger.Warnf("addr: %s/p2p/%s", addr.String(), server.h.ID())
	}

	//ping service
	server.ping = ping.NewPingService(server.h)
	//id service
	server.ids, err = identify.NewIDService(server.h)
	if err != nil {
		return err
	}
	server.ids.Start()

	for _, s := range server.services {
		s.Watch(ctx, server)
	}
	//protect bootstrap peers
	go func() {
		for _, bp := range server.peers {
			server.h.ConnManager().Protect(bp.ID, "bootstrap")
			connectErr := server.h.Connect(ctx, bp)
			if connectErr != nil {
				server.logger.Errorf("connect to bootstrap peer %s error: %s", bp.ID, connectErr)
				continue
			}
			server.h.Peerstore().AddAddrs(bp.ID, bp.Addrs, peerstore.PermanentAddrTTL)
		}
		server.logger.Warnf("connect to bootstrap peers success")
	}()

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
		_ = server.Bootstrapper.Close() // stop previous bootstrap process.
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

func (server *NodeServer) RegisterApi(route *http.Router) {
	route.GET("/peers", func(ctx http.Context) error {
		return ctx.JSON(htp.StatusOK, server.Peers())
	})
	route.GET("/conn", func(ctx http.Context) error {
		return ctx.JSON(htp.StatusOK, server.GetConnMgr())
	})
	route.GET("/addrs", func(ctx http.Context) error {
		ret := ""
		func() {
			defer func() {
				if err := recover(); err != nil {
					server.logger.Errorf("api query addrs panic: %v", err)
				}
			}()
			ret += fmt.Sprintf("addr count: %d\n", len(server.h.Addrs()))
			for _, addr := range server.h.Addrs() {
				ret += addr.String() + "/p2p/" + server.h.ID().String() + "\n"
			}
		}()

		return ctx.String(htp.StatusOK, ret)
	})

	route.POST("/connect", func(ctx http.Context) error {
		var info AddrInfo
		if err := ctx.Bind(&info); err != nil {
			return ctx.String(htp.StatusBadRequest, fmt.Sprintf("bad request error: %s", err))
		}

		addr, err := peer.AddrInfoFromString(info.Addr)
		if err != nil {
			return ctx.String(htp.StatusBadRequest, fmt.Sprintf("bad addr: %s", err))
		}
		err = server.h.Connect(ctx, *addr)
		if err != nil {
			return ctx.String(htp.StatusInternalServerError, fmt.Sprintf("connect failed: %s", err))
		}
		return ctx.String(htp.StatusOK, fmt.Sprintf("connect %s success", addr.ID.String()))
	})

	route.GET("/findpeer", func(ctx http.Context) error {
		var p Peer
		err := ctx.BindQuery(&p)
		if err != nil {
			return ctx.String(htp.StatusBadRequest, fmt.Sprintf("bad request error: %s", err))
		}
		info, err := server.rt.FindPeer(ctx, peer.ID(p.Id))
		if err != nil {
			return ctx.String(htp.StatusInternalServerError, fmt.Sprintf("find peer failed: %s", err))
		}
		return ctx.String(htp.StatusOK, fmt.Sprintf("find peer success: %s", info.String()))
	})
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
