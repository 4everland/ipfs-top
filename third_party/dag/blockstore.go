package dag

import (
	"context"
	pb "github.com/4everland/ipfs-servers/api/blockstore"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	grpc2 "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"time"
)

type grpcBlockstore struct {
	client pb.BlockstoreClient
}

func unWrapperError(err error) error {
	if pb.IsIpldNotFound(err) {
		return ipld.ErrNotFound{}
	}
	return err
}

func NewBlockStore(endpoint, cert string) (blockstore.Blockstore, error) {
	var o []grpc.DialOption
	if cert == "" {
		o = append(o, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		cred, err := credentials.NewClientTLSFromFile(cert, strings.Split(endpoint, ":")[0])
		if err != nil {
			return nil, err
		}
		o = append(o, grpc.WithTransportCredentials(cred))
	}
	o = append(o, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                30 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             2 * time.Second,  // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}))

	conn, err := grpc2.Dial(
		context.Background(),
		grpc2.WithEndpoint(endpoint),
		grpc2.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc2.WithTimeout(time.Minute),
		grpc2.WithOptions(o...),
	)
	if err != nil {
		return nil, err
	}
	return &grpcBlockstore{
		client: pb.NewBlockstoreClient(conn),
	}, nil
}

func (bs *grpcBlockstore) DeleteBlock(ctx context.Context, c cid.Cid) error {
	_, err := bs.client.DeleteBlock(ctx, &pb.Cid{Str: c.Bytes()})
	return unWrapperError(err)
}

func (bs *grpcBlockstore) Has(ctx context.Context, c cid.Cid) (bool, error) {
	has, err := bs.client.Has(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return false, unWrapperError(err)
	}
	return has.Value, unWrapperError(err)
}

func (bs *grpcBlockstore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	b, err := bs.client.Get(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return nil, unWrapperError(err)
	}
	bb, err := blocks.NewBlockWithCid(b.Data, c)
	return bb, unWrapperError(err)
}

func (bs *grpcBlockstore) GetSize(ctx context.Context, c cid.Cid) (int, error) {
	b, err := bs.client.GetSize(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return 0, unWrapperError(err)
	}
	return int(b.Value), nil
}

func (bs *grpcBlockstore) Put(ctx context.Context, b blocks.Block) error {
	_, err := bs.client.Put(ctx, &pb.Block{
		Cid:  &pb.Cid{Str: b.Cid().Bytes()},
		Data: b.RawData(),
	})
	return unWrapperError(err)
}

func (bs *grpcBlockstore) PutMany(ctx context.Context, bks []blocks.Block) error {
	srv, err := bs.client.PutMany(ctx)
	if err != nil {
		return unWrapperError(err)
	}
	//defer srv.CloseSend()
	for _, b := range bks {
		err = srv.Send(&pb.Block{
			Cid:  &pb.Cid{Str: b.Cid().Bytes()},
			Data: b.RawData(),
		})
		if err != nil {
			return unWrapperError(err)
		}
	}
	_, err = srv.CloseAndRecv()
	return unWrapperError(err)
}

func (bs *grpcBlockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	allkeySt, err := bs.client.AllKeysChan(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, unWrapperError(err)
	}
	ch := make(chan cid.Cid)
	go func() {
		defer close(ch)
		for {
			c, er := allkeySt.Recv()
			if er != nil {
				return
			}
			cc, _ := cid.Cast(c.Str)
			ch <- cc
		}

	}()
	return ch, nil
}

func (bs *grpcBlockstore) HashOnRead(enabled bool) {

}
