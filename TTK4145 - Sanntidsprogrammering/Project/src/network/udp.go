package network

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"utilities"
)

type UdpConnection struct {
	Addr  string
	Timer *time.Timer
}

var baddr *net.UDPAddr //Broadcast address

type udpMessage struct {
	raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	data   []byte
	length int //length of received data, in #bytes // N/A for sending
}

func udpInit(localListenPort, broadcastListenPort, message_size int, send_ch, receive_ch chan udpMessage) (err error) {
	//Generating broadcast address
	baddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}

	//Generating locaLocalIPess
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	localIP, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	localIP.Port = localListenPort
	utilities.LocalIP = localIP.String()

	//Creating local listening connections
	localListenConn, err := net.ListenUDP("udp4", localIP)
	if err != nil {
		return err
	}

	//Creating listener on broadcast connection
	broadcastListenConn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		localListenConn.Close()
		return err
	}

	go udp_receive_server(localListenConn, broadcastListenConn, message_size, receive_ch)
	go udp_transmit_server(localListenConn, broadcastListenConn, send_ch)
	go udp_connection_closer(localListenConn, broadcastListenConn)
	return err
}

func udp_transmit_server(lconn, bconn *net.UDPConn, send_ch <-chan udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_transmit_server: %s \n Closing connection.", r)
			lconn.Close()
			bconn.Close()
		}
	}()

	var err error
	var n int

	for {
		msg := <-send_ch
		if msg.raddr == "broadcast" {
			n, err = lconn.WriteToUDP(msg.data, baddr)
		} else {
			raddr, err := net.ResolveUDPAddr("udp", msg.raddr)
			if err != nil {
				fmt.Printf("Error: udp_transmit_server: could not resolve raddr\n")
				panic(err)
			}
			n, err = lconn.WriteToUDP(msg.data, raddr)
		}
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_transmit_server: writing\n")
			panic(err)
		}
	}
}

func udp_receive_server(lconn, bconn *net.UDPConn, message_size int, receive_ch chan<- udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_receive_server: %s \n Closing connection.", r)
			lconn.Close()
			bconn.Close()
		}
	}()

	bconn_rcv_ch := make(chan udpMessage)
	lconn_rcv_ch := make(chan udpMessage)

	go udp_connection_reader(lconn, message_size, lconn_rcv_ch)
	go udp_connection_reader(bconn, message_size, bconn_rcv_ch)

	for {
		select {

		case buf := <-bconn_rcv_ch:
			receive_ch <- buf

		case buf := <-lconn_rcv_ch:
			receive_ch <- buf
		}
	}
}

func udp_connection_reader(conn *net.UDPConn, message_size int, rcv_ch chan<- udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_connection_reader: Closing connection.", r)
			conn.Close()
		}
	}()

	for {
		buf := make([]byte, message_size)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_connection_reader: reading\n")
			panic(err)
		}
		rcv_ch <- udpMessage{raddr: raddr.String(), data: buf, length: n}
	}
}

func udp_connection_closer(lconn, bconn *net.UDPConn) {
	<-utilities.CloseConnectionChan
	lconn.Close()
	bconn.Close()
}
