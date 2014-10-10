package main

import (
	"../config"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(0)
	cfgFilePtr := flag.String("config", "config.yml", "Path to the config file.")
	flag.Parse()

	cfg, err := config.Read(*cfgFilePtr)
	if err != nil {
		log.Fatalf("failed to process config file: %s\n", err)
	}

	obs, err := NewObserver(cfg)
	if err != nil {
		log.Fatalf("failed to initialize: %s\n", err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Println("stopping")
		obs.Stop()
	}()

	log.Println("starting")
	obs.Run()
	log.Println("stopped")
}
