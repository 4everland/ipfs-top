package dag

import (
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
)

func NewGrpcDagService(blockEndpoint, exchangeEndpoint string) (format.DAGService, error) {
	service, err := NewGrpcBlockService(blockEndpoint, exchangeEndpoint)

	if err != nil {
		return nil, err
	}
	return merkledag.NewDAGService(service), nil
}
