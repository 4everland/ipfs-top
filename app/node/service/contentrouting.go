package service

import (
	"context"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/routing"
	"sync"

	pb "github.com/4everland/ipfs-servers/api/contentrouting"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	contentRoutingServiceWatchOnce = sync.Once{}
	contentRoutingServiceOnce      = sync.Once{}

	contentRoutingService *ContentRoutingService
)

type ContentRoutingService struct {
	pb.UnimplementedContentRoutingServer
	rt routing.Routing
}

func NewContentRoutingService() *ContentRoutingService {
	contentRoutingServiceOnce.Do(func() {
		contentRoutingService = &ContentRoutingService{}
	})
	return contentRoutingService
}

func (s *ContentRoutingService) Provide(ctx context.Context, req *pb.ProvideReq) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Cid.Str)
	if err != nil {
		return nil, err
	}
	err = s.rt.Provide(ctx, c, req.Provide)
	return &emptypb.Empty{}, err
}
func (s *ContentRoutingService) FindProvidersAsync(req *pb.GetProvidersReq, conn pb.ContentRouting_FindProvidersAsyncServer) error {
	//todo... reprovider do not need this api
	for {
		err := conn.Send(&pb.AddrInfo{})
		if err != nil {
			return err
		}
	}
}

func (s *ContentRoutingService) Watch(_ context.Context, node NodeInterface) {
	contentRoutingServiceWatchOnce.Do(func() {
		s.rt = node.GetContentRouting()
	})
}
