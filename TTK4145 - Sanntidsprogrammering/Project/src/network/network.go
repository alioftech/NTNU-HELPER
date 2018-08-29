/*----------------------------------------------------------------------------------
network package implements a UDP broadcast service, setting up communication between
the connected nodes. A udpMessage struct is broadcasted with address, data and
length (when receiving.)
------------------------------------------------------------------------------------
MESSAGE-HANDLING:
	Package json implements encoding and decoding of JSON.
		-> Marshal function before sending.
		-> Unmarshal function when receiving.
----------------------------------------------------------------------------------*/
package network

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"utilities"
)

func Init(outgoingMsg, incomingMsg chan utilities.Message) {

	const messageSize = 1024

	var udpSend = make(chan udpMessage)
	var udpReceive = make(chan udpMessage, 10)
	err := udpInit(utilities.LocalListenPort, utilities.BroadcastListenPort, messageSize, udpSend, udpReceive)
	if err != nil {
		fmt.Print("UdpInit() error:\n", err)
	}

	go aliveHeartbeat(outgoingMsg)
	go forwardOutgoing(outgoingMsg, udpSend)
	go forwardIncoming(incomingMsg, udpReceive)

	log.Println("Network initialised.")
}

//aliveHeartbeat periodically sends messages on the network to notify all lifts that it's still alive.
func aliveHeartbeat(outgoingMsg chan<- utilities.Message) {
	alive := utilities.Message{Category: utilities.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		outgoingMsg <- alive
		time.Sleep(utilities.AliveHeartbeat)
	}
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp file
// as JSON.
func forwardOutgoing(outgoingMsg <-chan utilities.Message, udpSend chan<- udpMessage) {
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("json.Marshal error: %v\n", err)
		}
		udpSend <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func forwardIncoming(incomingMsg chan<- utilities.Message, udpReceive <-chan udpMessage) {
	for {
		udpMessage := <-udpReceive
		var message utilities.Message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}
		message.Addr = udpMessage.raddr
		incomingMsg <- message
	}
}
