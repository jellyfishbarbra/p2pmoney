package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/jellyfishbarbra/verzo/network/peering"
)

var config = struct {
	debugMode   bool
	portPeering int
	logLevel    log.Level
	connectSeed bool
	version     string
}{version: "0.0.1"}

func main() {
	setup()
	padlock, err := ioutil.ReadFile("motd")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(padlock))
	log.SetHandler(text.New(os.Stdout))
	log.WithField("version", config.version).Info("initializing client")
	log.Info("initializing peering daemon")
	peeringDaemon, err := peering.NewDaemon(peering.NewConfig(config.version, config.portPeering, config.logLevel))
	if err != nil {
		log.WithField("error", err.Error()).WithField("port", config.portPeering).Fatal("could not start peering daemon")
	}

	err = peeringDaemon.Listen()
	if err != nil {
		log.WithField("error", err.Error()).WithField("port", config.portPeering).Fatal("failed to launch peering daemon")
	}
	log.Infof("started peering daemon on port %d", config.portPeering)

	if config.connectSeed {
		log.Info("[client] requesting connection to seed host")
		peeringDaemon.Connect("172.105.75.16")
	}
	for {
	}
}

func setup() {
	flag.BoolVar(&config.debugMode, "debug", false, "show debug log")
	flag.IntVar(&config.portPeering, "port", 5555, "custom peering daemon port")
	flag.BoolVar(&config.connectSeed, "seed", false, "connect to seed host")
	flag.Parse()

	if config.debugMode {
		config.logLevel = log.DebugLevel
	} else {
		config.logLevel = log.InfoLevel
	}
	log.SetLevel(config.logLevel)
}
