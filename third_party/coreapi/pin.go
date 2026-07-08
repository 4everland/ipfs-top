package coreapi

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/4everland/ipfs-top/api/pin"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	caopts "github.com/ipfs/kubo/core/coreiface/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcPin struct {
	client pb.PinClient
}

func NewPinAPI(endpoint string) (coreiface.PinAPI, error) {
	tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(endpoint, tlsOption)
	if err != nil {
		return nil, err
	}
	return &grpcPin{
		client: pb.NewPinClient(conn),
	}, nil
}

func (gp *grpcPin) Add(ctx context.Context, p path.Path, opts ...caopts.PinAddOption) error {
	settings, err := caopts.PinAddOptions(opts...)
	if err != nil {
		return err
	}

	_, err = gp.client.Add(ctx, &pb.AddReq{
		Path:      p.String(),
		Recursive: settings.Recursive,
	})

	return err
}

func (gp *grpcPin) IsPinned(ctx context.Context, p path.Path, opts ...caopts.PinIsPinnedOption) (string, bool, error) {
	settings, err := caopts.PinIsPinnedOptions(opts...)
	if err != nil {
		return "", false, err
	}

	resp, err := gp.client.IsPinned(ctx, &pb.IsPinnedReq{
		Path:     p.String(),
		WithType: settings.WithType,
	})
	if err != nil {
		return "", false, err
	}

	return resp.Cid, resp.IsPinned, nil
}

func (gp *grpcPin) Rm(ctx context.Context, p path.Path, opts ...caopts.PinRmOption) error {
	settings, err := caopts.PinRmOptions(opts...)
	if err != nil {
		return err
	}

	_, err = gp.client.Rm(ctx, &pb.RmReq{
		Path:      p.String(),
		Recursive: settings.Recursive,
	})

	return err
}

func (gp *grpcPin) Update(ctx context.Context, from path.Path, to path.Path, opts ...caopts.PinUpdateOption) error {
	settings, err := caopts.PinUpdateOptions(opts...)
	if err != nil {
		return err
	}

	_, err = gp.client.Update(ctx, &pb.UpdateReq{
		Form:  from.String(),
		To:    to.String(),
		Unpin: settings.Unpin,
	})

	return err
}

type pinInfo struct {
	pinType string
	path    path.ImmutablePath
	err     error
}

func (p *pinInfo) Path() path.ImmutablePath {
	return p.path
}

func (p *pinInfo) Name() string {
	return ""
}

func (p *pinInfo) Type() string {
	return p.pinType
}

func (gp *grpcPin) Ls(ctx context.Context, out chan<- coreiface.Pin, opts ...caopts.PinLsOption) error {
	settings, err := caopts.PinLsOptions(opts...)
	if err != nil {
		return err
	}

	switch settings.Type {
	case "all", "direct", "indirect", "recursive":
	default:
		return fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", settings.Type)
	}

	cc, err := gp.client.Ls(ctx, &pb.LsReq{Type: settings.Type})
	if err != nil {
		return err
	}
	defer cc.CloseSend()
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			recv, er := cc.Recv()
			if er != nil {
				return nil
			}
			if recv.Err != "" {
				return errors.New(recv.Err)
			}
			out <- &pinInfo{
				pinType: recv.PinType,
				path:    path.FromCid(cid.MustParse(recv.Cid)),
			}
		}
	}
}

func (gp *grpcPin) Verify(ctx context.Context) (<-chan coreiface.PinStatus, error) {
	//todo
	ch := make(chan coreiface.PinStatus)
	return ch, nil
}
