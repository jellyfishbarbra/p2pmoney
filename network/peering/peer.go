package peering

import (
	"fmt"
	"net"

	"github.com/apex/log"

	"time"
)

// DiscoveryMessage is ...
type DiscoveryMessage struct {
	Version string
	ts      time.Time
}

// DiscoveryAck is ...
type DiscoveryAck struct {
	Accept bool
}

// Heartbeat is ...
type Heartbeat struct {
	Version string
}

// Config is ...
type Config struct {
	Version  string
	port     int
	logLevel log.Level
}

// Daemon is a ...
type Daemon struct {
	c Config
}

// NewDiscoveryMessage is ...
func NewDiscoveryMessage(v string) DiscoveryMessage {
	return DiscoveryMessage{
		v,
		time.Now(),
	}
}

// NewDiscoveryAck is a ...
func NewDiscoveryAck(b bool) DiscoveryAck {
	return DiscoveryAck{b}
}

// NewConfig is ..
func NewConfig(ver string, port int, lvl log.Level) Config {
	return Config{
		Version:  ver,
		port:     port,
		logLevel: lvl,
	}
}

// NewDaemon is ...
func NewDaemon(c Config) (*Daemon, error) {
	return &Daemon{c}, nil
}

func (d *Daemon) formatPort() string {
	return fmt.Sprintf(":%d", d.c.port)
}

func (d *Daemon) Listen() error {
	listener, err := net.Listen("tcp", d.formatPort())
	if err != nil {
		log.WithField("error", err).Error("[peer] failed listening")
		return err
	}

	go handle(listener)
	return nil
}

func handle(incoming net.Listener) {
	for {
		conn, err := incoming.Accept()
		if err != nil {
			log.WithField("error", err.Error()).
				Error("[peer] failed to accept incoming connection")
		}
		log.WithField("from", conn.RemoteAddr().String()).
			Info("[peer] handling incoming connection")
	}
}
