package main

import (
	"ddclient/dungeonclient"
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
	mode = flag.String("mode", "player", "player or dm mode")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dungeonclient.RunClient(*addr, *mode, interrupt)

}
