package main

impoort(
	"io"
	"log"
	"net"
)

func main(){
	listner, err := net.Listen("udp", ":16569")
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}