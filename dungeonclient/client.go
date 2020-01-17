// Felipe Ryan - yep the variable, etc naming in here is meant as a "tongue in cheek"
// collection of D&D nerdyness.

package dungeonclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var fileName string
var dungeonMasterMode = false

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
	case "/audio":
		return message{Action: "audio", Payload: tokens[1], Recipient: ""}
	case "/d":
		return message{Action: "d", Payload: "", Recipient: ""}
	}

	return message{Action: "say", Payload: action, Recipient: "all"}

}

func decodeMessage(received []byte, messageToUI chan interface{}) {
	m := &message{}
	if err := json.Unmarshal(received, m); err != nil {
		log.Println("error unmarshalling: ", err)
	}

	switch m.Action {
	case "say":
		fmt.Printf("%s -> %s \n", m.Sender, m.Payload)
		if !dungeonMasterMode {
			messageToUI <- m.Payload
		}
	case "file":
		fileName = m.Payload
	}

}

func sendMessage(rawMessage string, connection *websocket.Conn) error {
	msg := createMessage(rawMessage)
	packet, err := json.Marshal(msg)
	if err != nil {
		log.Println("write:", err)
		return err
	}

	//test for special case, audio will be sent as a binary
	if msg.Action == "audio" {
		dat, err := ioutil.ReadFile(msg.Payload)
		if err != nil {
			log.Printf("Error opening file: %s\n", err)
			return err
		}
		err = connection.WriteMessage(websocket.BinaryMessage, dat)
		if err != nil {
			log.Println("write:", err)
			return err
		}
		return nil
	}

	err = connection.WriteMessage(websocket.TextMessage, []byte(packet))
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

func RunClient(serverAddress, mode string, interrupt chan os.Signal,
	waiting chan bool, fromUI chan interface{}, toUI chan interface{}, audioToUI chan interface{}) {

	connectTo := serverAddress
	if mode != "dm" {
		fmt.Println("blocking on fromUI channel")
		select {
		case textFromUI := <-fromUI: // server address
			connectTo = fmt.Sprintf("%v", textFromUI)
		case <-interrupt:
			close(waiting)
			return
		}

	}

	u := url.URL{Scheme: "ws", Host: connectTo, Path: "/receive"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error connecting:", err)
	}
	defer c.Close()

	toUI <- "Connected to Server"

	done := make(chan struct{})

	// this goroutine pings every 45 seconds to keep the connection alive on Heroku
	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
				if err != nil {
					return
				}
			}
		}
	}()

	// this goroutine receives messages from the websocket
	go func() {
		defer close(done)
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("error read: ", err)
				return
			}
			if msgType == websocket.TextMessage {
				decodeMessage(msg, toUI)
			}
			if msgType == websocket.BinaryMessage {
				if !dungeonMasterMode {
					// send bytes to audio channel for UI to pick up and play
					audioToUI <- msg
				}
			}
		}
	}()

	// this goroutine reads messages from stdin so the DM can send commands
	if mode == "dm" {
		dungeonMasterMode = true
		fmt.Print("*** Welcome master! ***\n")
		reader := bufio.NewReader(os.Stdin)
		go func() {
			for {
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				text = strings.Replace(text, "\n", "", -1)

				err := sendMessage(text, c)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}()
	}

	// this for-loop selects on different channels
	for {
		select {
		case textFromUI := <-fromUI: // atm only used for the /setname coming from the UI
			err := sendMessage(fmt.Sprintf("%v", textFromUI), c)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-done: // Probably because the server closed the connection
			fmt.Println("Looks like the server closed the connection")
			close(waiting)
			return
		case <-interrupt:
			// Prob a ctrl+c or window close
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
			close(waiting)
			return
		}
	}
}
