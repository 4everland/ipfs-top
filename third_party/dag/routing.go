package dag

import (
	"context"
	pb "github.com/4everland/ipfs-servers/api/routing"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	exchange "github.com/ipfs/boxo/exchange"
	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	ma "github.com/multiformats/go-multiaddr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type GrpcNodeRouting interface {
	routing.Routing
	exchange.Interface
}

type grpcRouting struct {
	client pb.RoutingClient
}

func (g grpcRouting) Close() error {
	return nil
}

func (g grpcRouting) Provide(ctx context.Context, c cid.Cid, b bool) error {
	_, err := g.client.Provide(ctx, &pb.ProvideReq{
		Cid:     &pb.Cid{Str: c.Bytes()},
		Provide: b,
	})
	return err
}

func (g grpcRouting) FindProvidersAsync(ctx context.Context, c cid.Cid, i int) <-chan peer.AddrInfo {
	cc, err := g.client.FindProvidersAsync(ctx, &pb.GetProvidersReq{
		Cid:   &pb.Cid{Str: c.Bytes()},
		Count: int64(i),
	})
	peerOut := make(chan peer.AddrInfo)
	if err != nil {
		close(peerOut)
		return peerOut
	}

	go func() {
		defer cc.CloseSend()
		for {
			//temp := make([]byte, 0)
			select {
			case <-ctx.Done():
				return
			default:
				recv, er := cc.Recv()
				if er != nil {
					return
				}
				peerOut <- convertPb2AddrInfo(recv)
			}
		}
	}()

	return peerOut

}

func (g grpcRouting) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	addr, err := g.client.FindPeer(ctx, &pb.Peer{Buf: string(id)})
	if err != nil {
		return peer.AddrInfo{}, err
	}
	addrs := make([]ma.Multiaddr, len(addr.Multiaddr))
	for i, b := range addr.Multiaddr {
		addrs[i], _ = ma.NewMultiaddrBytes(b)
	}
	return convertPb2AddrInfo(addr), nil
}

func (g grpcRouting) PutValue(ctx context.Context, s string, bytes []byte, options ...routing.Option) error {
	op := &routing.Options{}
	op.Apply(options...)

	_, err := g.client.PutValue(ctx, &pb.PutReq{
		Key:  s,
		Body: bytes,
		Option: &pb.ValueStoreOption{
			Expired: op.Expired,
			Offline: op.Offline,
		},
	})
	return err
}

func (g grpcRouting) GetValue(ctx context.Context, s string, options ...routing.Option) ([]byte, error) {
	op := &routing.Options{}
	op.Apply(options...)
	value, err := g.client.GetValue(ctx, &pb.GetReq{
		Key: s,
		Option: &pb.ValueStoreOption{
			Expired: op.Expired,
			Offline: op.Offline,
		},
	})
	if err != nil {
		return nil, err
	}
	return value.Data, nil
}

func (g grpcRouting) SearchValue(ctx context.Context, s string, options ...routing.Option) (<-chan []byte, error) {
	op := &routing.Options{}
	op.Apply(options...)

	cc, err := g.client.SearchValue(ctx, &pb.SearchReq{
		Key: s,
		Option: &pb.ValueStoreOption{
			Expired: op.Expired,
			Offline: op.Offline,
		},
	})
	if err != nil {
		return nil, err
	}
	var ch = make(chan []byte)
	go func() {
		defer cc.CloseSend()
		for {
			//temp := make([]byte, 0)
			select {
			case <-ctx.Done():
				return
			default:
				recv, er := cc.Recv()
				if er != nil {
					return
				}
				ch <- recv.Data
			}
		}
	}()
	return ch, nil
}

// GetBlock TODO... Custom Session
func (g grpcRouting) GetBlock(ctx context.Context, k cid.Cid) (blocks.Block, error) {
	value, err := g.client.GetBlock(ctx, &pb.Cid{
		Str: k.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return blocks.NewBlockWithCid(value.Data, k)
}

// GetBlocks TODO... Custom Session
func (g grpcRouting) GetBlocks(ctx context.Context, keys []cid.Cid) (<-chan blocks.Block, error) {
	cids := make([]*pb.Cid, len(keys))
	for i, v := range keys {
		cids[i] = &pb.Cid{
			Str: v.Bytes(),
		}
	}
	cc, err := g.client.GetBlocks(ctx, &pb.Cids{Cid: cids})
	if err != nil {
		return nil, err
	}
	var ch = make(chan blocks.Block)
	go func() {
		defer cc.CloseSend()
		for {
			//temp := make([]byte, 0)
			select {
			case <-ctx.Done():
				return
			default:
				recv, er := cc.Recv()
				if er != nil {
					return
				}
				ccc, er := cid.Cast(recv.Cid.Str)
				if er != nil {
					return
				}
				block, er := blocks.NewBlockWithCid(recv.Data, ccc)
				if er != nil {
					return
				}
				ch <- block
			}
		}
	}()
	return ch, nil
}

func (g grpcRouting) NotifyNewBlocks(ctx context.Context, blocks ...blocks.Block) error {
	return nil
	for _, block := range blocks {
		err := g.Provide(ctx, block.Cid(), true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g grpcRouting) Bootstrap(ctx context.Context) error {
	return nil
}

func NewGrpcRouting(endpoint string) (GrpcNodeRouting, error) {
	tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc2.Dial(
		context.Background(),
		grpc2.WithEndpoint(endpoint),
		grpc2.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc2.WithTimeout(time.Minute),
		grpc2.WithOptions(tlsOption),
	)

	if err != nil {
		return nil, err
	}
	return &grpcRouting{
		client: pb.NewRoutingClient(conn),
	}, nil
}

func convertPb2AddrInfo(addr *pb.AddrInfo) peer.AddrInfo {
	addrs := make([]ma.Multiaddr, len(addr.Multiaddr))
	for i, b := range addr.Multiaddr {
		addrs[i], _ = ma.NewMultiaddrBytes(b)
	}
	return peer.AddrInfo{
		ID:    peer.ID(addr.Id),
		Addrs: addrs,
	}
}
