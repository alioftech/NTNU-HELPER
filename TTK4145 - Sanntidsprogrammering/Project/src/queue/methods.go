package queue

import (
	"hw"
	"log"
	"time"
	"utilities"
)

func (q *queue) startTimer(floor, button int) {
	q.OrdersList[floor][button].Timer = time.NewTimer(utilities.OrderTimeout)
	<-q.OrdersList[floor][button].Timer.C
	OrderTimeoutChan <- hw.Keypress{Button: button, Floor: floor}
}

func (q *queue) stopTimer(floor, button int) {
	if q.OrdersList[floor][button].Timer != nil {
		q.OrdersList[floor][button].Timer.Stop()
	}
}

func (q *queue) isEmpty() bool {
	for f := 0; f < utilities.NFloors; f++ {
		for b := 0; b < utilities.NButtons; b++ {
			if q.OrdersList[f][b].Active {
				return false
			}
		}
	}
	return true
}

func (q *queue) addOrder(floor, button int, status orderStatus) {
	if q.isOrder(floor, button) == status.Active {
		// Ignore if order is already in queue.
		return
	}
	q.OrdersList[floor][button] = status
	takeBackup <- true
	syncLightsChan <- true
	printQueues()
}

func (q *queue) isOrder(floor, button int) bool {
	return q.OrdersList[floor][button].Active
}

func (q *queue) isOrdersAbove(floor int) bool {
	for f := floor + 1; f < utilities.NFloors; f++ {
		for b := 0; b < utilities.NButtons; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < utilities.NButtons; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor, dir int) int {
	if q.isEmpty() {
		return utilities.DirStop
	}
	switch dir {
	case utilities.DirDown:
		if q.isOrdersBelow(floor) && floor > 0 {
			return utilities.DirDown
		} else {
			return utilities.DirUp
		}
	case utilities.DirUp:
		if q.isOrdersAbove(floor) && floor < utilities.NFloors-1 {
			return utilities.DirUp
		} else {
			return utilities.DirDown
		}
	case utilities.DirStop:
		if q.isOrdersAbove(floor) {
			return utilities.DirUp
		} else if q.isOrdersBelow(floor) {
			return utilities.DirDown
		} else {
			return utilities.DirStop
		}
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Printf("chooseDirection(): called with invalid direction %d, returning stop\n", dir)
		return 0
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case utilities.DirDown:
		return q.isOrder(floor, utilities.BtnDown) ||
			q.isOrder(floor, utilities.BtnInside) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case utilities.DirUp:
		return q.isOrder(floor, utilities.BtnUp) ||
			q.isOrder(floor, utilities.BtnInside) ||
			floor == utilities.NFloors-1 ||
			!q.isOrdersAbove(floor)
	case utilities.DirStop:
		return q.isOrder(floor, utilities.BtnDown) ||
			q.isOrder(floor, utilities.BtnUp) ||
			q.isOrder(floor, utilities.BtnInside)
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Fatalln("Non-existing direction, restarting ...")
	}
	return false
}

func (q *queue) deepCopy() *queue {
	queueCopy := new(queue)
	for f := 0; f < utilities.NFloors; f++ {
		for b := 0; b < utilities.NButtons; b++ {
			queueCopy.OrdersList[f][b] = q.OrdersList[f][b]
		}
	}
	return queueCopy
}
