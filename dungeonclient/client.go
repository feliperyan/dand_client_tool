package dungeonclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var fileName string

type message struct {
	Action    string `json:"action"`
	Payload   string `json:"payload"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

func createMessage(command string) message {
	action := strings.Trim(command, " ")
	tokens := strings.Split(action, " ")

	switch tokens[0] {
	case "/setname":
		text := strings.Join(tokens[1:len(tokens)], " ")
		return message{Action: "setname", Payload: text, Recipient: ""}
	case "/whisper":
		text := strings.Join(tokens[2:len(tokens)], " ")
		return message{Action: "whisper", Payload: text, Recipient: tokens[1]}
	case "/list":
		return message{Action: "list", Payload: "", Recipient: ""}
	case "/file":
		return message{Action: "file", Payload: tokens[1], Recipient: ""}
	}

	return message{Action: "say", Payload: action, Recipient: "all"}

}

func decodeMessage(received []byte) {
	m := &message{}
	if err := json.Unmarshal(received, m); err != nil {
		log.Println("error unmarshalling: ", err)
	}

	switch m.Action {
	case "say":
		fmt.Printf("%s -> %s \n", m.Sender, m.Payload)
	case "file":
		fileName = m.Payload
	}

}

func writeFile(message []byte, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("error creating file: %s\n", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	n4, err := w.Write(message)
	if err != nil {
		log.Printf("error writing file: %s\n", err)
		return
	}
	log.Printf("Wrote %d bytes\n", n4)
	w.Flush()
}

func RunClient(serverAddress, mode string, interrupt chan os.Signal) {

	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/receive"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("error read: ", err)
				return
			}
			if msgType == websocket.TextMessage {
				decodeMessage(msg)
			}
			if msgType == websocket.BinaryMessage {
				if fileName != "" {
					writeFile(msg, "example.mp3")
					fileName = ""
				}
			}
		}
	}()

	if mode == "dm" {
		fmt.Print("*** Welcome master! ***\n")
		reader := bufio.NewReader(os.Stdin)
		go func() {
			for {
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				text = strings.Replace(text, "\n", "", -1)

				msg := createMessage(text)
				packet, err := json.Marshal(msg)
				if err != nil {
					log.Println("write:", err)
					return
				}

				err = c.WriteMessage(websocket.TextMessage, []byte(packet))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}()
	}

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
