package services

import (
	"context"
	"errors"
	pb "github.com/4everland/ipfs-servers/api/blockstore"
	biz2 "github.com/4everland/ipfs-servers/app/blockstore/internal/biz"
	"github.com/4everland/ipfs-servers/app/blockstore/internal/utils"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
)

type BlockstoreService struct {
	pb.UnimplementedBlockstoreServer
	backend biz2.BlockStore
	index   biz2.BlockIndex
}

func NewBlockstoreService(backend biz2.BlockStore, index biz2.BlockIndex) *BlockstoreService {
	return &BlockstoreService{
		backend: backend,
		index:   index,
	}
}

func (s *BlockstoreService) DeleteBlock(ctx context.Context, req *pb.Cid) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	err = s.index.Delete(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	err = s.backend.Delete(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *BlockstoreService) Has(ctx context.Context, req *pb.Cid) (*wrapperspb.BoolValue, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	exists, err := s.index.Has(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	return &wrapperspb.BoolValue{Value: exists}, nil
}

func (s *BlockstoreService) Get(ctx context.Context, req *pb.Cid) (*pb.Block, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	exists, err := s.index.Has(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if !exists {
		return nil, utils.GrpcErrorWrapper(ipld.ErrNotFound{})
	}
	r, err := s.backend.Get(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	return &pb.Block{
		Cid:  req,
		Data: data,
	}, nil
}

func (s *BlockstoreService) GetSize(ctx context.Context, req *pb.Cid) (*wrapperspb.Int32Value, error) {
	c, err := cid.Cast(req.Str)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	iv, err := s.index.Get(ctx, c.String())
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	size, _ := iv.Size()
	return &wrapperspb.Int32Value{Value: int32(size)}, nil
}

func (s *BlockstoreService) Put(ctx context.Context, req *pb.Block) (*emptypb.Empty, error) {
	c, err := cid.Cast(req.Cid.Str)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	if c.Version() == 0 {
		c = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	}
	err = s.backend.Put(ctx, c.String(), c.String(), req.Data)
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
	}
	err = s.index.Put(ctx, c.String(), biz2.NewIndexValue(0, uint32(len(req.Data)), 0, c.String()))
	if err != nil {
		return nil, utils.GrpcErrorWrapper(err)
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
			return utils.GrpcErrorWrapper(err)
		}

		if _, err = s.Put(context.Background(), b); err != nil {
			return utils.GrpcErrorWrapper(err)
		}
	}

}

func (s *BlockstoreService) AllKeysChan(_ *emptypb.Empty, conn pb.Blockstore_AllKeysChanServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keyCh := s.index.List(ctx)
	for key := range keyCh {
		c, err := cid.Cast([]byte(key))
		if err != nil {
			continue
		}
		err = conn.Send(&pb.Cid{Str: c.Bytes()})
		if err != nil {
			return utils.GrpcErrorWrapper(err)
		}
	}
	return nil
}
