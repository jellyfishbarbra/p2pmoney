package peering

import (
	"bufio"
	"bytes"
	"encoding/gob"
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

func (d *Daemon) Connect(remote string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:5555", remote))
	defer conn.Close()
	if err != nil {
		log.WithField("error", err.Error()).
			WithField("remote", remote).
			WithField("port", "5555").
			Error("[peer] failed connecting to remote host")
		return err
	}

	disc := NewDiscoveryMessage(d.c.Version)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(disc)
	written, err := conn.Write(buffer.Bytes())
	if err != nil {
		log.WithField("written", written).
			WithField("error", err.Error()).
			WithField("remote", remote).
			Error("[peer] failed to transmit DiscoveryMessage to remote host")
		return err
	}

	if written <= 4 {
		log.WithField("bytes written", written).
			Warn("[peer] bytes of DiscoveryMessage written to remote host are low")
	}

	log.WithField("remote", remote).
		WithField("port", 5555).
		Info("[peer] sent a DiscoveryMessage to remote host")
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

		// buffer := make([]byte, 128)
		// read, err := conn.Read(buffer)
		c := bufio.NewReader(conn)
		// size, err := c.ReadByte()
		// if err != nil {
		// 	log.WithField("error", err.Error()).
		// 		WithField("bytes read", size).
		// 		WithField("remote", conn.RemoteAddr().String).
		// 		Error("[peer] connection handler could not read payload")
		// 	err = conn.Close()
		// 	// TODO: figure out what scenarios lead to this + how to recover
		// 	if err != nil {
		// 		log.WithField("error", err.Error()).
		// 			Fatal("[peer] could not lose connection!")
		// 	}
		// 	continue
		// }

		// _, err = io.ReadFull(c, buffer[:int(size)])
		// if err != nil {
		// 	log.WithField("size", size).
		// 		Error("[peer] could not read full payload")
		// 	continue
		// }

		var msg DiscoveryMessage
		dec := gob.NewDecoder(c)
		err = dec.Decode(&msg)
		if err != nil {
			log.WithField("error", err.Error()).
				Error("[peer] could not decode payload")
			continue
		}
		log.WithField("remote version", msg.Version).
			WithField("remote", conn.RemoteAddr().String()).
			Info("[peer] received DiscoveryMessage from remote peer")
	}
}
