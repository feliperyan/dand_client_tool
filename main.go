package main

import (
	"ddclient/dungeonclient"
	"ddclient/uicomponents"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hajimehoshi/ebiten"
)

var (
	addr           = flag.String("addr", "localhost:8080", "http service address")
	mode           = flag.String("mode", "player", "player or dm mode")
	messagesFromUI chan interface{}
	messagesToUI   chan interface{}
	audioToUI      chan interface{}
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	waitToClose := make(chan bool)

	messagesFromUI = make(chan interface{})
	messagesToUI = make(chan interface{}, 100) // need to buffer here as ebiten's Update wont be waiting/blocking on the other side
	audioToUI = make(chan interface{}, 10)     // used for audio []byte

	game := &uicomponents.Game{}
	game.OutgoingMSG = messagesFromUI
	game.IncomingMSG = messagesToUI
	game.AudioInMSG = audioToUI

	update := game.Update

	// run the client in a separate go routine as ebiten.Run didn't like it.
	go dungeonclient.RunClient(*addr, *mode, interrupt, waitToClose, messagesFromUI, messagesToUI, audioToUI)

	// GM interacts via console so only players see the UI
	if *mode != "dm" {
		fmt.Println("Client running - not in DM mode.")
		ebiten.SetRunnableInBackground(true)
		if err := ebiten.Run(update, uicomponents.ScreenWidth, uicomponents.ScreenHeight, 1, "D and D"); err != nil {
			log.Fatal(err)
		}
		// ebiten.Run will return nil after we close the window so pop an os.Interrupt sigal
		// into the interrupt channel to simulate ctrl+c in the terminal
		interrupt <- os.Interrupt
	}

	<-waitToClose
}
