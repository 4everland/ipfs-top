package service

import (
	"context"
	pb "github.com/4everland/ipfs-servers/api/provide"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/routing"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"

	provider "github.com/ipfs/go-ipfs-provider"
)

var (
	simpleProvideOnce      = sync.Once{}
	simpleProvideWatchOnce = sync.Once{}

	simpleProvideService *ProvideService
)

type ProvideService struct {
	pb.UnimplementedProvideServer

	bprovide   func(rt routing.ContentRouting) provider.Provider
	breprobide func(rt routing.ContentRouting) provider.Reprovider

	p provider.Provider
	r provider.Reprovider
}

func NewSimpleProvideService(bprovide func(rt routing.ContentRouting) provider.Provider, breprobide func(rt routing.ContentRouting) provider.Reprovider) *ProvideService {
	simpleProvideOnce.Do(func() {
		simpleProvideService = &ProvideService{
			bprovide:   bprovide,
			breprobide: breprobide,
		}
	})
	return simpleProvideService
}

func (s *ProvideService) Provide(_ context.Context, req *pb.Cid) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Buf)
	if err != nil {
		return nil, err
	}
	err = s.p.Provide(c)
	return &emptypb.Empty{}, err
}

func (s *ProvideService) Reprovide(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.r.Trigger(ctx)
	return &emptypb.Empty{}, err
}

func (s *ProvideService) Watch(_ context.Context, node NodeInterface) {
	simpleProvideWatchOnce.Do(func() {
		rt := node.GetContentRouting()
		s.p = s.bprovide(rt)
		s.r = s.breprobide(rt)
		go s.p.Run()
		go s.r.Run()
	})

}
