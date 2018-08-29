// Exercise1 project main.go
package main

import (
	"fmt"
	"runtime"
	"time"
)

var x = 0

func add(a *int) {
	for i := 0; i < 1000000; i++ {
		x++
		*a = i
	}
}

func sub(a *int) {
	for i := 0; i < 1000000; i++ {
		x--
		*a = i
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var a, b = 0, 0
	go add(&a)
	go sub(&b)

	for i := 0; i < 50; i++ {
		fmt.Println("Value: ", x)
	}

	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Done, value: ", x)
	fmt.Println("Done, number of loops including 0: ADD() = ", a, " SUB() = ", b)
}