package coreapi

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/4everland/ipfs-top/api/pin"
	"github.com/4everland/ipfs-top/third_party/coreunix"
	"github.com/4everland/ipfs-top/third_party/coreunix/options"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcPin struct {
	client pb.PinClient
}

func NewPinAPI(endpoint string) (coreunix.PinAPI, error) {
	tlsOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(endpoint, tlsOption)
	if err != nil {
		return nil, err
	}
	return &grpcPin{
		client: pb.NewPinClient(conn),
	}, nil
}

func (gp *grpcPin) Add(ctx context.Context, p path.Path, opts ...options.PinAddOption) error {
	settings, err := options.PinAddOptions(opts...)
	if err != nil {
		return err
	}

	_, err = gp.client.Add(ctx, &pb.AddReq{
		Path:      p.String(),
		Recursive: settings.Recursive,
	})

	return err
}

func (gp *grpcPin) IsPinned(ctx context.Context, p path.Path, opts ...options.PinIsPinnedOption) (string, bool, error) {
	settings, err := options.PinIsPinnedOptions(opts...)
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

func (gp *grpcPin) Rm(ctx context.Context, p path.Path, opts ...options.PinRmOption) error {
	settings, err := options.PinRmOptions(opts...)
	if err != nil {
		return err
	}

	_, err = gp.client.Rm(ctx, &pb.RmReq{
		Path:      p.String(),
		Recursive: settings.Recursive,
	})

	return err
}

func (gp *grpcPin) Update(ctx context.Context, from path.Path, to path.Path, opts ...options.PinUpdateOption) error {
	settings, err := options.PinUpdateOptions(opts...)
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
	name    string
	path    path.ImmutablePath
	err     error
}

func (p *pinInfo) Path() path.ImmutablePath {
	return p.path
}

func (p *pinInfo) Type() string {
	return p.pinType
}

func (p *pinInfo) Name() string {
	return p.name
}

func (p *pinInfo) Err() error {
	return p.err
}

func (gp *grpcPin) Ls(ctx context.Context, opts ...options.PinLsOption) (<-chan coreunix.Pin, error) {
	settings, err := options.PinLsOptions(opts...)
	if err != nil {
		return nil, err
	}

	switch settings.Type {
	case "all", "direct", "indirect", "recursive":
	default:
		return nil, fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", settings.Type)
	}

	cc, err := gp.client.Ls(ctx, &pb.LsReq{Type: settings.Type})
	ch := make(chan coreunix.Pin)
	go func() {
		defer cc.CloseSend()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				recv, er := cc.Recv()
				if er != nil {
					return
				}
				if recv.Err != "" {
					ch <- &pinInfo{
						err: errors.New(recv.Err),
					}
				} else {
					ch <- &pinInfo{
						pinType: recv.PinType,
						path:    path.FromCid(cid.MustParse(recv.Cid)),
					}
				}

			}
		}
	}()

	return ch, nil
}

func (gp *grpcPin) Verify(ctx context.Context) (<-chan coreunix.PinStatus, error) {
	//todo
	ch := make(chan coreunix.PinStatus)
	return ch, nil
}
