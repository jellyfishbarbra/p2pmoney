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
}{}

func main() {
	setup()
	padlock, err := ioutil.ReadFile("padlock")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(padlock))
	log.SetHandler(text.New(os.Stdout))
	log.Info("initializing client")
	log.Info("initializing peering daemon")
	peeringDaemon, err := peering.NewDaemon(peering.NewConfig("0.0.1", config.portPeering, config.logLevel))
	if err != nil {
		log.WithField("error", err.Error()).WithField("port", config.portPeering).Fatal("could not start peering daemon")
	}

	err = peeringDaemon.Listen()
	if err != nil {
		log.WithField("error", err.Error()).WithField("port", config.portPeering).Fatal("failed to launch peering daemon")
	}
	log.Infof("started peering daemon on port %d", config.portPeering)
	for {
	}
}

func setup() {
	flag.BoolVar(&config.debugMode, "debug", false, "show debug log")
	flag.IntVar(&config.portPeering, "port", 5555, "custom peering daemon port")
	flag.Parse()

	if config.debugMode {
		config.logLevel = log.DebugLevel
	} else {
		config.logLevel = log.InfoLevel
	}
	log.SetLevel(config.logLevel)
}
