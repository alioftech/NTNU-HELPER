package main



import (
    . "fmt"
    "runtime"
)

var i int = 0


func threadFunction_1(message chan int, done_message_1 chan bool){
	for k := 0; k < 1000002; k++{
		i := <- message
		i ++
		message <- i
	}
	done_message_1 <- true
}


func threadFunction_2(message chan int, done_message_2 chan bool){
	for k := 0; k < 1000000; k++{
		i := <- message
		i --
		message <- i
	}
	done_message_2 <- true
}


func main(){

	message := make(chan int, 1)

	done_message := make(chan bool)

	runtime.GOMAXPROCS(runtime.NumCPU())

	message <- 0
	go threadFunction_1(message, done_message)
	go threadFunction_2(message, done_message)

	<-done_message
	<-done_message

	i := <- message
	Println("Value after both threads running, i = ", i)
}