package utils

import (
	"github.com/4everland/ipfs-servers/api/blockstore"
	ipld "github.com/ipfs/go-ipld-format"
)

func GrpcErrorWrapper(err error) error {
	if err == nil {
		return nil
	}
	if ipld.IsNotFound(err) {
		return blockstore.ErrorIpldNotFound(err.Error())
	}
	return blockstore.ErrorUnknown(err.Error())
}
