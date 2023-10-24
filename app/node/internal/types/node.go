package types

import (
	"fmt"
	"time"
)

type ConnectPeer struct {
	Id        string
	Local     string `json:"local"`
	Addr      string `json:"addr,omitempty"`
	Direction string
	Opened    time.Time
	Transient bool
}

func (cp ConnectPeer) String() string {
	return fmt.Sprintf("%-59s, direct: %s, local: %s, remote: %s", cp.Id, cp.Direction, cp.Local, cp.Addr)
}
