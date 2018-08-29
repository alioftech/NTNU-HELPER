package main

import (
	"coordinator"
	"fsm"
	"hw"
	"liftSelector"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
	"utilities"
)

func main() {

	var (
		numOnline int
		floor     int
		err       error
	)
	floor, err = hw.Init()
	if err != nil {
		utilities.Restart.Run()
		log.Fatal(err)
	}
	go safeKill()

	ch := fsm.Channels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan int, 10),
		FloorLamp:    make(chan int, 10),
		DoorLamp:     make(chan bool, 10),
		OutgoingMsg:  make(chan utilities.Message, 10),
		IncomingMsg:  make(chan utilities.Message, 10),
		NOnline:      numOnline,
		OnlineLifts:  make(map[string]network.UdpConnection),
		CostChan:     make(chan utilities.Message),
		DeadChan:     make(chan network.UdpConnection),
	}

	network.Init(ch.OutgoingMsg, ch.IncomingMsg)
	fsm.Init(ch, floor)
	queue.Init(ch.NewOrder, ch.OutgoingMsg)
	coordinator.Init(ch)

	go liftSelector.Run(ch.CostChan, &numOnline)

	select {} //keep the program running.
}

//safeKill shutdowns the motor when an interupt signal is recieved
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDir(utilities.DirStop)
	log.Fatal("\nExiting program goodbye \t(•‿•) ")
}
