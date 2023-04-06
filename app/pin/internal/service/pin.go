package service

import (
	"context"
	"fmt"
	pb "github.com/4everland/ipfs-servers/api/pin"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	pin "github.com/ipfs/go-ipfs-pinner"
	"github.com/ipfs/go-ipfs-pinner/dspinner"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/interface-go-ipfs-core/path"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PinService struct {
	pb.UnimplementedPinServer
	//pinning coreiface.PinAPI

	dagResolver coreunix.DagResolve
	blockstore  blockstore.GCBlockstore
	dag         ipld.DAGService

	pinning pin.Pinner
}

type emptyGCBlockstore struct {
	blockstore.GCLocker
	blockstore.Blockstore
}

func NewPinService(blockStore blockstore.Blockstore, datastore datastore.Datastore, exchange exchange.Interface) (*PinService, error) {
	bs := blockservice.New(blockStore, exchange)
	dag := merkledag.NewDAGService(bs)
	pinner, err := dspinner.New(context.Background(), datastore, dag)
	if err != nil {
		return nil, err
	}

	return &PinService{
		dagResolver: coreunix.NewDagResolver(context.Background(), dag, bs),
		blockstore: emptyGCBlockstore{
			blockstore.NewGCLocker(),
			blockStore,
		},
		dag:     dag,
		pinning: pinner,
	}, nil

}
func (s *PinService) Add(ctx context.Context, req *pb.AddReq) (*emptypb.Empty, error) {
	dagNode, err := s.dagResolver.ResolveNode(ctx, path.New(req.Path))
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("pin: %s", err)
	}

	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	err = s.pinning.Pin(ctx, dagNode, req.Recursive)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("pin: %s", err)
	}

	//if err := api.provider.Provide(dagNode.Cid()); err != nil {
	//	return err
	//}

	return &emptypb.Empty{}, nil
}

func (s *PinService) IsPinned(ctx context.Context, req *pb.IsPinnedReq) (*pb.IsPinnedResp, error) {
	resolved, err := s.dagResolver.ResolvePath(ctx, path.New(req.Path))
	if err != nil {
		return &pb.IsPinnedResp{}, fmt.Errorf("error resolving path: %s", err)
	}

	mode, ok := pin.StringToMode(req.WithType)
	if !ok {
		return &pb.IsPinnedResp{}, fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", req.WithType)
	}

	cid, isPinned, err := s.pinning.IsPinnedWithType(ctx, resolved.Cid(), mode)
	return &pb.IsPinnedResp{
		Cid:      cid,
		IsPinned: isPinned,
	}, nil
}

func (s *PinService) Rm(ctx context.Context, req *pb.RmReq) (*emptypb.Empty, error) {
	rp, err := s.dagResolver.ResolvePath(ctx, path.New(req.Path))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Note: after unpin the pin sets are flushed to the blockstore, so we need
	// to take a lock to prevent a concurrent garbage collection
	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	if err = s.pinning.Unpin(ctx, rp.Cid(), req.Recursive); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, s.pinning.Flush(ctx)
}

func (s *PinService) Update(ctx context.Context, req *pb.UpdateReq) (*emptypb.Empty, error) {
	fp, err := s.dagResolver.ResolvePath(ctx, path.New(req.Form))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	tp, err := s.dagResolver.ResolvePath(ctx, path.New(req.To))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	err = s.pinning.Update(ctx, fp.Cid(), tp.Cid(), req.Unpin)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, s.pinning.Flush(ctx)
}

func (s *PinService) Ls(req *pb.LsReq, server pb.Pin_LsServer) error {
	switch req.Type {
	case "all", "direct", "indirect", "recursive":
	default:
		return fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", req.Type)
	}

	keys := cid.NewSet()

	AddToResultKeys := func(keyList []cid.Cid, typeStr string) error {
		for _, c := range keyList {
			if keys.Visit(c) {
				select {
				case <-server.Context().Done():
					return server.Context().Err()
				default:
					server.Send(&pb.PinInfo{
						PinType: typeStr,
						Cid:     c.String(),
					})
				}
			}
		}
		return nil
	}

	VisitKeys := func(keyList []cid.Cid) {
		for _, c := range keyList {
			keys.Visit(c)
		}
	}

	go func() {
		var dkeys, rkeys []cid.Cid
		var err error
		if req.Type == "recursive" || req.Type == "all" {
			rkeys, err = s.pinning.RecursiveKeys(server.Context())
			if err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
			if err = AddToResultKeys(rkeys, "recursive"); err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
		}
		if req.Type == "direct" || req.Type == "all" {
			dkeys, err = s.pinning.DirectKeys(server.Context())
			if err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
			if err = AddToResultKeys(dkeys, "direct"); err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
		}
		if req.Type == "all" {
			set := cid.NewSet()
			for _, k := range rkeys {
				err = merkledag.Walk(
					server.Context(), merkledag.GetLinksWithDAG(s.dag), k,
					set.Visit,
					merkledag.SkipRoot(), merkledag.Concurrent(),
				)
				if err != nil {
					server.Send(&pb.PinInfo{Err: err.Error()})
					return
				}
			}
			if err = AddToResultKeys(set.Keys(), "indirect"); err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
		}
		if req.Type == "indirect" {
			// We need to first visit the direct pins that have priority
			// without emitting them

			dkeys, err = s.pinning.DirectKeys(server.Context())
			if err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
			VisitKeys(dkeys)

			rkeys, err = s.pinning.RecursiveKeys(server.Context())
			if err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
			VisitKeys(rkeys)

			set := cid.NewSet()
			for _, k := range rkeys {
				err = merkledag.Walk(
					server.Context(), merkledag.GetLinksWithDAG(s.dag), k,
					set.Visit,
					merkledag.SkipRoot(), merkledag.Concurrent(),
				)
				if err != nil {
					server.Send(&pb.PinInfo{Err: err.Error()})
					return
				}
			}
			if err = AddToResultKeys(set.Keys(), "indirect"); err != nil {
				server.Send(&pb.PinInfo{Err: err.Error()})
				return
			}
		}
	}()

	return nil
}

func (s *PinService) Verify(empty *emptypb.Empty, server pb.Pin_VerifyServer) error {
	// todo
	return nil
}
