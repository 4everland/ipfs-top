package dag

import (
	"context"
	pb "github.com/4everland/ipfs-servers/api/blockstore"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-libipfs/blocks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcBlockstore struct {
	client pb.BlockstoreClient
}

func NewBlockStore(endpoint string) (*grpcBlockstore, error) {
	tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(endpoint, tlsOption)
	if err != nil {
		return nil, err
	}
	return &grpcBlockstore{
		client: pb.NewBlockstoreClient(conn),
	}, nil
}

func (bs *grpcBlockstore) DeleteBlock(ctx context.Context, c cid.Cid) error {
	_, err := bs.client.DeleteBlock(ctx, &pb.Cid{Str: c.Bytes()})
	return err
}

func (bs *grpcBlockstore) Has(ctx context.Context, c cid.Cid) (bool, error) {
	has, err := bs.client.Has(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return false, err
	}
	return has.Value, err
}

func (bs *grpcBlockstore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	b, err := bs.client.Get(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return nil, err
	}
	return blocks.NewBlockWithCid(b.Data, c)
}

func (bs *grpcBlockstore) GetSize(ctx context.Context, c cid.Cid) (int, error) {
	b, err := bs.client.GetSize(ctx, &pb.Cid{Str: c.Bytes()})
	if err != nil {
		return 0, err
	}
	return int(b.Value), nil

}

func (bs *grpcBlockstore) Put(ctx context.Context, b blocks.Block) error {
	_, err := bs.client.Put(ctx, &pb.Block{
		Cid:  &pb.Cid{Str: b.Cid().Bytes()},
		Data: b.RawData(),
	})
	return err
}

func (bs *grpcBlockstore) PutMany(ctx context.Context, bks []blocks.Block) error {
	srv, err := bs.client.PutMany(ctx)
	if err != nil {
		return err
	}
	//defer srv.CloseSend()
	for _, b := range bks {
		err = srv.Send(&pb.Block{
			Cid:  &pb.Cid{Str: b.Cid().Bytes()},
			Data: b.RawData(),
		})
		if err != nil {
			return err
		}
	}
	_, err = srv.CloseAndRecv()
	return err
}

func (bs *grpcBlockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	allkeySt, err := bs.client.AllKeysChan(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
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
