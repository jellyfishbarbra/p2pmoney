package peering

import (
	"net"
	"time"
)

type peer struct {
	addr net.Addr
	version string
}

type daemon struct {

}

type connection struct {

}

func newPeer(addr net.Addr) (*Peer, error) {
	return &Peer{}, nil
}

func (p *Peer) connect() (*connection, error) {
	return nil, nil
}

func (d *daemon) listen() error {
	return nil
}
