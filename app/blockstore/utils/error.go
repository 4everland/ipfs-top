package utils

import (
	"github.com/4everland/ipfs-servers/api/blockstore"
	"github.com/go-kratos/kratos/v2/errors"
	ipld "github.com/ipfs/go-ipld-format"
)

func GrpcErrorWrapper(err error) *errors.Error {
	if ipld.IsNotFound(err) {
		return blockstore.ErrorIpldNotFound(err.Error())
	}
	return blockstore.ErrorUnknown(err.Error())
}
