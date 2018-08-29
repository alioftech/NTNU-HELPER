/*----------------------------------------------------------------------------------
fsm package implements a "finite-state-machine" for the elevator behaviour.
The elevators takes orders based on a queue stored and managed by the queue package.

The state machine has 3 main states:
	Idle: 		Elevator is stationary, at a floor, door closed, waiting for orders.
	Moving: 	Lift is moving, inbetween or passing floors.
	Door open: 	Lift is at a floor with the door open.
----------------------------------------------------------------------------------*/
package fsm

import (
	"log"
	"network"
	"queue"
	"time"
	"utilities"
)

const (
	idle int = iota
	moving
	doorOpen
)

var (
	state int
	floor int
	dir   int
)

type Channels struct {
	// Hardware Interaction
	MotorDir  chan int
	FloorLamp chan int
	DoorLamp  chan bool

	// Orders
	NewOrder       chan bool
	FloorReached   chan int
	doorTimeout    chan bool
	doorTimerReset chan bool

	// Network Interaction
	OnlineLifts map[string]network.UdpConnection
	NOnline     int
	OutgoingMsg chan utilities.Message
	IncomingMsg chan utilities.Message
	DeadChan    chan network.UdpConnection
	CostChan    chan utilities.Message
}

func Init(ch Channels, startFloor int) {
	state = idle
	dir = utilities.DirStop
	floor = startFloor

	ch.doorTimeout = make(chan bool)
	ch.doorTimerReset = make(chan bool)

	go doorTimer(ch.doorTimeout, ch.doorTimerReset)
	go run(ch)
	log.Println("FSM initialised.")
}

//Direction returns the stored dir inside of the state machine
func Direction() int {
	return dir
}

//Floor returns the stored floor inside of the state machine
func Floor() int {
	return floor
}

func run(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			nextOrder(ch)
		case floor := <-ch.FloorReached:
			completeOrder(ch, floor)
		case <-ch.doorTimeout:
			doorTimeout(ch)
		}
	}
}

//nextOrder assures that the system is in one of the three main states or else it will
//send close connection signal and restart the program.
func nextOrder(ch Channels) {
	log.Printf("EVENT: New order in state.%v", liftStates(state))
	switch state {
	case moving:
		//do nothing
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
			ch.DoorLamp <- true
			state = doorOpen
		} else {
			ch.MotorDir <- dir
			state = moving
		}
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
		}
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Fatalf("This state doesn't exist, restarting...")
	}
}

//completeOrder checks if the lfit should stop on the currentFloor
//and removes the order from the orderList (restarts if in a weird state)
func completeOrder(ch Channels, newFloor int) {
	log.Printf("EVENT: Floor %d reached in state %s.", newFloor+1, liftStates(state))
	floor = newFloor
	ch.FloorLamp <- floor
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
			ch.DoorLamp <- true
			dir = utilities.DirStop
			ch.MotorDir <- dir
			state = doorOpen
		}
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Fatalf("Lift is in an unrecognized state, restarting %s ...\n", liftStates(state))
	}
}

//doorTimeout changes the state of the fsm once the lift is ready again
//turns off the DoorLamp and selects the next Direction
func doorTimeout(ch Channels) {
	log.Printf("EVENT: Door timeout in state %s.", liftStates(state))
	switch state {
	case doorOpen:
		ch.DoorLamp <- false
		dir = queue.ChooseDirection(floor, dir)
		ch.MotorDir <- dir
		if dir == utilities.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Fatalln("Lift is in an unrecognized state, restarting ...")
	}
}

//liftStates returns a string depending on the state
func liftStates(state int) string {
	switch state {
	case idle:
		return "idle"
	case doorOpen:
		return "door open"
	case moving:
		return "moving"
	default:
		return "error: bad state"
	}
}

//doorTimer keeps track of door open timeout  and sends to the timeout channel (inside fsm).
func doorTimer(timeout chan<- bool, reset <-chan bool) {
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(utilities.DoorOpenInterval)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}
