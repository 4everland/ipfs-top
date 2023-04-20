package types

import (
	"fmt"
	"github.com/libp2p/go-libp2p/core/network"
)

type ConnectPeer struct {
	Id      string
	Local   string `json:"local"`
	Addr    string `json:"addr,omitempty"`
	Network network.Stats
}

func (cp ConnectPeer) String() string {

	return fmt.Sprintf("%-59s, direct: %s, local: %s, remote: %s", cp.Id, cp.Network.Direction, cp.Local, cp.Addr)
}
