package main

import (
	"encoding/json"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/justincampbell/anybar"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
}

func setStatus(message Message) {
	log.Println("Got server_message: ", message.Event)
	if message.Event == "door_open" {
		anybar.Green()
	} else if message.Event == "door_close" {
		anybar.Red()
	} else {
		anybar.Question()
	}
}

func main() {
	fmt.Printf("Starting AnyBar...")
	err := exec.Command("open", "-a", "AnyBar").Start()
	if err != nil {
		log.Fatal("Error starting AnyBar")
	} else {
		log.Print("Done")
	}

	log.Print("Connecting to event stream...")
	//connect to server, you can use your own transport settings
	client, err := gosocketio.Dial(
		gosocketio.GetUrl("www.aaron-santos.com", 8081, false),
		transport.GetDefaultWebsocketTransport(),
	)

	if err != nil {
		log.Fatal("Error connecting: ", err)
	} else {
		log.Print("Connected")
	}

	//do something, handlers and functions are same as server ones
	//custom event handler
	client.On("server_message", func(c *gosocketio.Channel, args Message) string {
		log.Println("Got server_message: ", args.Event)
		setStatus(args)

		//you can return result of handler, in this case
		//handler will be converted from "emit" to "ack"
		return "result"
	})

	log.Printf("Getting initial status...")
	resp, err := http.Get("http://www.aaron-santos.com:8081/status")
	if err != nil {
		log.Fatal("Error getting status: ", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading status body: ", err)
		} else {
			var message Message
			json.Unmarshal(body, &message)
			setStatus(message)
		}
	}

	//sleep forever
	select {}
	//close connection
	client.Close()
}
