package service

import (
	"context"
	"github.com/ipfs/go-cid"
	"os"

	pb "github.com/4everland/ipfs-servers/api/provide"
	"google.golang.org/protobuf/types/known/emptypb"

	provider "github.com/ipfs/go-ipfs-provider"
)

type ProvideService struct {
	pb.UnimplementedProvideServer

	p provider.Provider
	r provider.Reprovider
}

func NewProvideService(p provider.Provider, r provider.Reprovider) *ProvideService {
	ps := &ProvideService{
		p: p,
		r: r,
	}
	go p.Run()
	go r.Run()
	return ps
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

func (s *ProvideService) Close(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	s.p.Close()
	s.r.Close()
	os.Exit(0)
	return &emptypb.Empty{}, nil
}
