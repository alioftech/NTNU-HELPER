/*---------------------------------------------------------------------------------------
queue package implements order operations; add, remove, update and terminal output.
---------------------------------------------------------------------------------------*/
package queue

import (
	"fmt"
	"hw"
	"log"
	"time"
	"utilities"
)

type orderStatus struct {
	Active bool
	Addr   string      `json:"-"`
	Timer  *time.Timer `json:"-"`
}

type queue struct {
	OrdersList [utilities.NFloors][utilities.NButtons]orderStatus
}

var (
	blankOrder = orderStatus{Active: false, Addr: "", Timer: nil}

	local  queue
	remote queue

	updateLocal      = make(chan bool)
	takeBackup       = make(chan bool, 10)
	OrderTimeoutChan = make(chan hw.Keypress)
	syncLightsChan   = make(chan bool)
	newOrder         chan bool
)

func Init(newOrderTemp chan bool, outgoingMsg chan utilities.Message) {

	newOrder = newOrderTemp
	go updateLocalQueue()
	go syncLights()
	runBackup(outgoingMsg)
	log.Println("Queue initialised.")
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.addOrder(floor, button, orderStatus{true, "", nil})
	newOrder <- true
}

// AddRemoteOrder adds an order to the remote queue, and spawns a timer
// for the order. (If the order times out, it will be taken care of.)
func AddRemoteOrder(floor, button int, addr string) {
	alreadyExist := IsRemoteOrder(floor, button)
	remote.addOrder(floor, button, orderStatus{true, addr, nil})
	if !alreadyExist {
		go remote.startTimer(floor, button)
	}
	updateLocal <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from remote queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < utilities.NButtons; b++ {
		remote.stopTimer(floor, b)
		remote.addOrder(floor, b, blankOrder)
	}
	updateLocal <- true
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int, outgoingMsg chan<- utilities.Message) {
	for b := 0; b < utilities.NButtons; b++ {
		remote.stopTimer(floor, b)
		local.addOrder(floor, b, blankOrder)
		remote.addOrder(floor, b, blankOrder)
	}
	outgoingMsg <- utilities.Message{Category: utilities.CompleteOrder, Floor: floor, Button: -1, Cost: -1}
}

// ShouldStop returns whether the lift should stop when it reaches the given
// floor, going in the given direction.
func ShouldStop(floor, dir int) bool {
	return local.shouldStop(floor, dir)
}

// ChooseDirection returns the direction the lift should continue after the
// current floor, going in the given direction.
func ChooseDirection(floor, dir int) int {
	return local.chooseDirection(floor, dir)
}

// IsLocalOrder returns whether there in an order with the given floor and
// button in the local queue.
func IsLocalOrder(floor, button int) bool {
	return local.isOrder(floor, button)
}

// IsRemoteOrder returns true if there is a order with the given floor and
// button in the remote queue.
func IsRemoteOrder(floor, button int) bool {
	return remote.isOrder(floor, button)
}

// ReassignOrders finds all orders assigned to a dead lift, removes them from
// the remote queue, and sends them on the network as new, unassigned orders.
func ReassignOrders(deadAddr string, outgoingMsg chan<- utilities.Message) {
	for f := 0; f < utilities.NFloors; f++ {
		for b := 0; b < utilities.NButtons; b++ {
			if remote.OrdersList[f][b].Addr == deadAddr {
				remote.addOrder(f, b, blankOrder)
				outgoingMsg <- utilities.Message{
					Category: utilities.NewOrder,
					Floor:    f,
					Button:   b}
			}
		}
	}
}

// printQueues prints the queues as simple art to terminal.
func printQueues() {
	fmt.Println("\t# Local\t\tRemote #      (IP)")
	for f := utilities.NFloors - 1; f >= 0; f-- {

		s1 := "    \t"
		if local.isOrder(f, utilities.BtnUp) {
			s1 += "  ↑"
		} else {
			s1 += "  "
		}
		if local.isOrder(f, utilities.BtnInside) {
			s1 += " ×"
		} else {
			s1 += "  "
		}
		fmt.Printf(s1)
		if local.isOrder(f, utilities.BtnDown) {
			fmt.Printf(" ↓    %d  ", f+1)
		} else {
			fmt.Printf("      %d  ", f+1)

		}

		s2 := "   \t"
		if remote.isOrder(f, utilities.BtnUp) {
			fmt.Printf("    ↑")
			s2 += "    (↑ " + remote.OrdersList[f][utilities.BtnUp].Addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isOrder(f, utilities.BtnDown) {
			fmt.Printf("   ↓")
			s2 += "    (↓ " + remote.OrdersList[f][utilities.BtnDown].Addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", s2)
		fmt.Println()
		fmt.Println("\t| -------------------- |")
	}
}

// updateLocalQueue checks remote queue for new orders assigned to this lift
// and copies them to the local queue.
func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < utilities.NFloors; f++ {
			for b := 0; b < utilities.NButtons; b++ {
				if remote.isOrder(f, b) {
					if b != utilities.BtnInside && remote.OrdersList[f][b].Addr == utilities.LocalIP {
						if !local.isOrder(f, b) {
							local.addOrder(f, b, orderStatus{true, "", nil})
							newOrder <- true
						}
					}
				}
			}
		}
	}
}

//syncLights checks the queues and updates all order lamps accordingly.
func syncLights() {
	for {
		<-syncLightsChan
		for f := 0; f < utilities.NFloors; f++ {
			for b := 0; b < utilities.NButtons; b++ {
				if (b == utilities.BtnUp && f == utilities.NFloors-1) || (b == utilities.BtnDown && f == 0) {
					continue
				} else {
					switch b {
					case utilities.BtnInside:
						hw.SetButtonLamp(f, b, IsLocalOrder(f, b))
					case utilities.BtnUp, utilities.BtnDown:
						hw.SetButtonLamp(f, b, IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
