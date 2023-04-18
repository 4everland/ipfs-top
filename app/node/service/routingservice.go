package service

import (
	"context"
	provider "github.com/ipfs/boxo/provider"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"sync"

	pb "github.com/4everland/ipfs-servers/api/routing"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	routingServiceWatchOnce = sync.Once{}
	routingServiceOnce      = sync.Once{}

	routingService *RoutingService
)

type RoutingService struct {
	pb.UnimplementedRoutingServer
	rt routing.Routing

	bprovide   func(rt routing.ContentRouting) provider.Provider
	breprobide func(rt routing.ContentRouting) provider.Reprovider

	p provider.Provider
	r provider.Reprovider

	bitSwapService *BitSwapService
}

func NewRoutingService(bitSwapService *BitSwapService, bprovide func(rt routing.ContentRouting) provider.Provider, breprobide func(rt routing.ContentRouting) provider.Reprovider) *RoutingService {
	routingServiceOnce.Do(func() {
		routingService = &RoutingService{
			bprovide:   bprovide,
			breprobide: breprobide,

			bitSwapService: bitSwapService,
		}
	})
	return routingService
}

func (s *RoutingService) Provide(_ context.Context, req *pb.ProvideReq) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Cid.Str)
	if err != nil {
		return nil, err
	}
	err = s.p.Provide(c)
	return &emptypb.Empty{}, err
}

func (s *RoutingService) FindProvidersAsync(req *pb.GetProvidersReq, conn pb.Routing_FindProvidersAsyncServer) error {
	c, err := cid.Cast(req.Cid.Str)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	result := s.rt.FindProvidersAsync(ctx, c, int(req.Count))
	if err != nil {
		cancel()
		return err
	}
	go func() {
		defer cancel()
		for addrInfo := range result {
			if er := conn.Send(convertAddrInfo2Pb(addrInfo)); er != nil {
				break
			}
		}
	}()
	return nil
}

func (s *RoutingService) PutValue(ctx context.Context, req *pb.PutReq) (*emptypb.Empty, error) {
	options := make([]routing.Option, 0)
	if req.Option.Expired {
		options = append(options, routing.Expired)
	}
	if req.Option.Offline {
		options = append(options, routing.Offline)
	}
	err := s.rt.PutValue(ctx, req.Key, req.Body, options...)
	return &emptypb.Empty{}, err
}

func (s *RoutingService) GetValue(ctx context.Context, req *pb.GetReq) (*pb.GetReply, error) {
	options := make([]routing.Option, 0)
	if req.Option.Expired {
		options = append(options, routing.Expired)
	}
	if req.Option.Offline {
		options = append(options, routing.Offline)
	}
	body, err := s.rt.GetValue(ctx, req.Key, options...)
	if err != nil {
		return nil, err
	}
	return &pb.GetReply{
		Data: body,
	}, nil
}

func (s *RoutingService) SearchValue(req *pb.SearchReq, resp pb.Routing_SearchValueServer) error {
	options := make([]routing.Option, 0)
	if req.Option.Expired {
		options = append(options, routing.Expired)
	}
	if req.Option.Offline {
		options = append(options, routing.Offline)
	}
	ctx, cancel := context.WithCancel(context.Background())
	result, err := s.rt.SearchValue(ctx, req.Key, options...)
	if err != nil {
		cancel()
		return err
	}
	go func() {
		defer cancel()
		for body := range result {
			if er := resp.Send(&pb.SearchReply{Data: body}); er != nil {
				break
			}
		}
	}()
	return nil
}

func (s *RoutingService) FindPeer(ctx context.Context, p *pb.Peer) (*pb.AddrInfo, error) {
	addr, err := s.rt.FindPeer(ctx, peer.ID(p.Buf))
	if err != nil {
		return nil, err
	}
	return convertAddrInfo2Pb(addr), nil
}

func (s *RoutingService) GetBlock(ctx context.Context, req *pb.Cid) (*pb.Block, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, err
	}
	block, err := s.bitSwapService.BitSwap().GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	return &pb.Block{
		Cid:  req,
		Data: block.RawData(),
	}, nil
}

func (s *RoutingService) GetBlocks(cids *pb.Cids, resp pb.Routing_GetBlocksServer) error {
	keys := make([]cid.Cid, 0, len(cids.Cid))
	for _, v := range cids.Cid {
		c, err := cid.Cast(v.Str)
		if err != nil {
			keys = append(keys, c)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	result, err := s.bitSwapService.BitSwap().GetBlocks(ctx, keys)
	if err != nil {
		cancel()
		return err
	}
	go func() {
		defer cancel()
		for body := range result {
			if er := resp.Send(&pb.Block{Cid: &pb.Cid{Str: body.Cid().Bytes()}, Data: body.RawData()}); er != nil {
				break
			}
		}
	}()
	return nil

}

func (s *RoutingService) Reprovide(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.r.Trigger(ctx)
	return &emptypb.Empty{}, err
}

func (s *RoutingService) Watch(_ context.Context, node NodeInterface) {
	routingServiceWatchOnce.Do(func() {
		s.rt = node.GetContentRouting()

		s.p = s.bprovide(s.rt)
		s.r = s.breprobide(s.rt)
		go s.p.Run()
		go s.r.Run()
	})
}

func convertAddrInfo2Pb(addr peer.AddrInfo) *pb.AddrInfo {
	multiaddrs := make([][]byte, len(addr.Addrs))
	for i, d := range addr.Addrs {
		multiaddrs[i] = d.Bytes()
	}
	return &pb.AddrInfo{
		Id: string(addr.ID),

		Multiaddr: multiaddrs,
	}
}
