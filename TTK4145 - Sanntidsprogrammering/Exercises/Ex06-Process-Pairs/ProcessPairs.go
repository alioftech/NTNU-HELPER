package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"time"
)

var buffer = make([]byte, 8)
var counter uint64 = 0

func spawnBackUp() {
	command := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run ProcessPairs.go")
	_ = command.Run()
	fmt.Println("New primary is running")
}

func initialize(udpaddr *net.UDPAddr) {
	var primary bool = false
	connection, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		fmt.Println("Problem with listening to the local adress")
	}
	fmt.Println("UDP connection established")
	for !primary {
		connection.SetReadDeadline(time.Now().Add(time.Second * 2))
		b, _, err := connection.ReadFromUDP(buffer)
		if err == nil {
			counter = binary.BigEndian.Uint64(buffer[0:b])
		} else {
			primary = true
		}
	}
	connection.Close()
}

func main() {
	udpaddr, _ := net.ResolveUDPAddr("udp", "129.241.187.156:30089")
	initialize(udpaddr)
	spawnBackUp()
	connection, _ := net.DialUDP("udp", nil, udpaddr)

	for {
		fmt.Println(counter)
		counter++
		binary.BigEndian.PutUint64(buffer, counter)
		_, _ = connection.Write(buffer)

		time.Sleep(time.Second)
	}

}
