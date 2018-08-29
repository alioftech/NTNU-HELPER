package main

import (
        "fmt"
        "net"
        "time"
)

const (
    servAddr = "129.241.187.136"
    termPort = "33546"
    host = "129.241.187.104"
)

func TCPclient() {
    serv, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(servAddr, termPort))
    if err != nil {
    	fmt.Println("Error in resolving tcp address:", net.JoinHostPort(servAddr, termPort), err.Error())
    }
    
    conn, err := net.DialTCP("tcp", nil, serv)
    if err != nil {
    	fmt.Printf("Error in TCP: %s\n", err.Error())
    }
    
    message := []byte("Connect to:"+net.JoinHostPort(host, termPort)+"\x00")
    conn.Write(message)
    
}

func TCPserver() {

	tcpLocalAddr, err := net.ResolveTCPAddr("tcp","129.241.187.104:33546")
    tcpListener, err := net.ListenTCP("tcp", tcpLocalAddr)
    if err != nil {
        fmt.Println("Error opening connection:", err.Error())
    }
    defer tcpListener.Close()
    for {
        // Listen for an incoming connection.
        conn, err := tcpListener.AcceptTCP()
        fmt.Println("accepted connection")
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
        }
        
        go handleConnection(conn)
    	
    }
}

func handleConnection(conn * net.TCPConn){
    data := make([]byte, 1024)
    message := make([]byte, 1024)
    data = []byte("Cookies for parties\x00")
    for {
    	
    	_, err := conn.Write(data)
    	if err != nil {
    		fmt.Printf("Error in TCP: %s\n", err.Error())
    		break
    	}
    	
    	time.Sleep(100*time.Millisecond)
    	
    	conn.Read(message)
    	fmt.Println("r: ", string(message))
    	
    }
}

func main() {
	doneChan := make(chan bool)

	go TCPserver()
	time.Sleep(100*time.Millisecond)
	go TCPclient()
	
	<-doneChan
}