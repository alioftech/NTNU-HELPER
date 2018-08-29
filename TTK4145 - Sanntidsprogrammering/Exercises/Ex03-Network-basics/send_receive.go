package main

import( 
	"net"
	"fmt"
	"time"
)

func main(){
	saddr, err := net.ResolveUDPAddr("udp4",port)
	l, err := ListenUDP("udp", saddr)
	var buff []byte = make([]byte, 1024)
	var ipAddr := "localhost"
	//var saddr := "127.0.0.1"
	conn, err := net.Dial("udp", saddr)

	defer conn.Close()
	        

	if err != nil {
		// handle error
	}

	n, err := conn.Write([]byte("Connect to:" + ipAddr ":" + port + "\0"))
	l.ReadFromUDP(buff)
	fmt.Println("Buffer udp: ", buff)

	 for {

                time.Sleep(1000*time.Millisecond)
                n, err := conn.Write([]byte("Connect to:" + ipAddr ":" + port + "\0"))
                if err != nil {
                        fmt.Println("error writing data to server", saddr)
                        fmt.Println(err)
                        return
                }

                if n > 0 {
                        fmt.Println("Wrote ",n, " bytes to server at ", saddr)
                }
            }
}
