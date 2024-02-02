package service

import (
	"context"
	"fmt"
	pb "github.com/4everland/ipfs-top/api/pin"
	"github.com/4everland/ipfs-top/third_party/coreunix"
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/path"

	//"github.com/ipfs/boxo/coreiface/path"
	exchange "github.com/ipfs/boxo/exchange"
	"github.com/ipfs/boxo/ipld/merkledag"
	pin "github.com/ipfs/boxo/pinning/pinner"
	"github.com/ipfs/boxo/pinning/pinner/dspinner"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	ipld "github.com/ipfs/go-ipld-format"

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
	p, err := path.NewPath(req.Path)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	dagNode, err := s.dagResolver.ResolveNode(ctx, p)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("pin: %s", err)
	}

	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	err = s.pinning.Pin(ctx, dagNode, req.Recursive, "")
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("pin: %s", err)
	}

	//if err := api.provider.Provide(dagNode.Cid()); err != nil {
	//	return err
	//}

	return &emptypb.Empty{}, nil
}

func (s *PinService) IsPinned(ctx context.Context, req *pb.IsPinnedReq) (*pb.IsPinnedResp, error) {
	p, err := path.NewPath(req.Path)
	if err != nil {
		return &pb.IsPinnedResp{}, err
	}

	resolved, err := s.dagResolver.ResolvePath(ctx, p)
	if err != nil {
		return &pb.IsPinnedResp{}, fmt.Errorf("error resolving path: %s", err)
	}

	mode, ok := pin.StringToMode(req.WithType)
	if !ok {
		return &pb.IsPinnedResp{}, fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", req.WithType)
	}

	c, isPinned, err := s.pinning.IsPinnedWithType(ctx, resolved.RootCid(), mode)
	return &pb.IsPinnedResp{
		Cid:      c,
		IsPinned: isPinned,
	}, nil
}

func (s *PinService) Rm(ctx context.Context, req *pb.RmReq) (*emptypb.Empty, error) {
	p, err := path.NewPath(req.Path)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	rp, err := s.dagResolver.ResolvePath(ctx, p)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Note: after unpin the pin sets are flushed to the blockstore, so we need
	// to take a lock to prevent a concurrent garbage collection
	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	if err = s.pinning.Unpin(ctx, rp.RootCid(), req.Recursive); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, s.pinning.Flush(ctx)
}

func (s *PinService) Update(ctx context.Context, req *pb.UpdateReq) (*emptypb.Empty, error) {
	p, err := path.NewPath(req.Form)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	p1, err := path.NewPath(req.To)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	fp, err := s.dagResolver.ResolvePath(ctx, p)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	tp, err := s.dagResolver.ResolvePath(ctx, p1)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	defer s.blockstore.PinLock(ctx).Unlock(ctx)

	err = s.pinning.Update(ctx, fp.RootCid(), tp.RootCid(), req.Unpin)
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

	AddToResultKeys := func(c cid.Cid, typeStr string) error {
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
		return nil
	}

	go func() {
		var (
			rkeys []cid.Cid
			err   error
		)
		if req.Type == "recursive" || req.Type == "all" {
			for streamedCid := range s.pinning.RecursiveKeys(server.Context(), true) {
				if streamedCid.Err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
				if err = AddToResultKeys(streamedCid.Pin.Key, "recursive"); err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
				rkeys = append(rkeys, streamedCid.Pin.Key)
			}
		}
		if req.Type == "direct" || req.Type == "all" {
			for streamedCid := range s.pinning.DirectKeys(server.Context(), false) {
				if streamedCid.Err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
				if err = AddToResultKeys(streamedCid.Pin.Key, "direct"); err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
			}
		}
		if req.Type == "indirect" {
			// We need to first visit the direct pins that have priority
			// without emitting them

			for streamedCid := range s.pinning.DirectKeys(server.Context(), false) {
				if streamedCid.Err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
				keys.Add(streamedCid.Pin.Key)
			}

			for streamedCid := range s.pinning.RecursiveKeys(server.Context(), true) {
				if streamedCid.Err != nil {
					server.Send(&pb.PinInfo{Err: streamedCid.Err.Error()})
					return
				}
				keys.Add(streamedCid.Pin.Key)
				rkeys = append(rkeys, streamedCid.Pin.Key)
			}
		}
		if req.Type == "indirect" || req.Type == "all" {
			walkingSet := cid.NewSet()
			for _, k := range rkeys {
				err = merkledag.Walk(
					server.Context(), merkledag.GetLinksWithDAG(s.dag), k,
					func(c cid.Cid) bool {
						if !walkingSet.Visit(c) {
							return false
						}
						if keys.Has(c) {
							return true // skipped
						}
						err := AddToResultKeys(c, "indirect")
						if err != nil {
							server.Send(&pb.PinInfo{Err: err.Error()})
							return false
						}
						return true
					},
					merkledag.SkipRoot(), merkledag.Concurrent(),
				)
				if err != nil {
					server.Send(&pb.PinInfo{Err: err.Error()})
					return
				}
			}
		}
	}()

	return nil
}

func (s *PinService) Verify(empty *emptypb.Empty, server pb.Pin_VerifyServer) error {
	// todo
	return nil
}
