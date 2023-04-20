package dag

import (
	"github.com/ipfs/boxo/ipld/merkledag"
	format "github.com/ipfs/go-ipld-format"
)

func NewGrpcDagService(blockEndpoint, exchangeEndpoint string) (format.DAGService, error) {
	service, err := NewGrpcBlockService(blockEndpoint, exchangeEndpoint)

	if err != nil {
		return nil, err
	}
	return merkledag.NewDAGService(service), nil
}
