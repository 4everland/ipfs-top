package services

import (
	"context"
	"errors"
	"github.com/4everland/ipfs-servers/app/blockstore/biz"
	"io"

	pb "github.com/4everland/ipfs-servers/api/blockstore"
	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type BlockstoreService struct {
	pb.UnimplementedBlockstoreServer
	backend biz.BlockStore
	index   biz.BlockIndex
}

func NewBlockstoreService(backend biz.BlockStore, index biz.BlockIndex) *BlockstoreService {
	return &BlockstoreService{
		backend: backend,
		index:   index,
	}
}

func (s *BlockstoreService) DeleteBlock(ctx context.Context, req *pb.Cid) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, err
	}
	err = s.index.Delete(ctx, c.String())
	if err != nil {
		return nil, err
	}
	err = s.backend.Delete(ctx, c.String())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BlockstoreService) Has(ctx context.Context, req *pb.Cid) (*wrapperspb.BoolValue, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, err
	}
	exists, err := s.index.Has(ctx, c.String())
	if err != nil {
		return nil, err
	}
	return &wrapperspb.BoolValue{Value: exists}, nil
}

func (s *BlockstoreService) Get(ctx context.Context, req *pb.Cid) (*pb.Block, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, err
	}
	exists, err := s.index.Has(ctx, c.String())
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	r, err := s.backend.Get(ctx, c.String())
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &pb.Block{
		Cid:  req,
		Data: data,
	}, nil
}

func (s *BlockstoreService) GetSize(ctx context.Context, req *pb.Cid) (*wrapperspb.Int32Value, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, err
	}
	iv, err := s.index.Get(ctx, c.String())
	if err != nil {
		return nil, err
	}
	size, _ := iv.Size()
	return &wrapperspb.Int32Value{Value: int32(size)}, nil
}

func (s *BlockstoreService) Put(ctx context.Context, req *pb.Block) (*emptypb.Empty, error) {
	c, err := cid.Cast([]byte(req.Cid.Str))
	if err != nil {
		return nil, err
	}
	err = s.backend.Put(ctx, c.String(), c.String(), req.Data)
	if err != nil {
		return nil, err
	}
	err = s.index.Put(ctx, c.String(), biz.NewIndexValue(0, uint32(len(req.Data)), 0, c.String()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BlockstoreService) PutMany(srv pb.Blockstore_PutManyServer) error {
	for {
		b, err := srv.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			_ = srv.SendAndClose(&emptypb.Empty{})
			return err
		}

		_, err = s.Put(context.Background(), b)
		if err != nil {
			return err
		}
	}

}

func (s *BlockstoreService) AllKeysChan(_ *emptypb.Empty, conn pb.Blockstore_AllKeysChanServer) error {
	//TODO... get all keys
	var i = 100
	for {
		i--
		err := conn.Send(&pb.Cid{})
		if err != nil {
			return err
		}
		if i <= 0 {
			break
		}
	}
	return nil
}
