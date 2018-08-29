package main

import( 
	"net"
	"fmt"
	"time"
)

const (
	servAddr = "129.241.187.255"
	udpPort = string(20000 + nr)
)

func udpSend(done chan bool, port , saddr *UDPAddr){
	conn, err := net.DialUDP("udp", nil, saddr)
	if err != nil {
		fmt.Println("Error connecting to" + saddr)
	}

	for {
		time.Sleep(1000*time.Millisecond)
		conn.Write([]byte("The cake is a lie"))
		fmt.Println("Msg sent on udp")
	}
	done <- true

}

func udpRecive(done chan bool, port , saddr * net.UDPAddr) {
	buff := make([]byte, 1024)

	l, err := net.ListenUDP("udp4", saddr)
	if err != nil {
		fmt.Println("Error listening to" + saddr)
	}

	_,_, err = l.ReadFromUDP(buff)

	if err != nil {
			fmt.Println(err)
	}
	fmt.Println(string(buff[:]))

}

func main() {
	done := make(chan bool)

	saddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(host, udpPort))

	if err != nil{
		fmt.Println("Failed to resolve address for: " + port)

	}

	go udpSend(done, port, saddr)
	go udpRecive(done, port, saddr)

	<-done
	<-done
}