package coreunix

import (
	"fmt"

	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
)

func NewPath(p string) (path.Path, error) {
	contentPath, err := path.NewPath(p)
	if err == nil {
		return contentPath, nil
	}

	c, cidErr := cid.Decode(p)
	if cidErr != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return path.FromCid(c), nil
}
