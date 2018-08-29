package queue

import (
	"log"
	"utilities"
)

// CalculateCost returns the cost for a given elevator to go through with a given order.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q := local.deepCopy()
	q.addOrder(targetFloor, utilities.BtnInside, orderStatus{true, "", nil})

	cost := 0
	floor := prevFloor
	dir := currDir
	// Between floors adds cost +1.
	if currFloor == -1 {
		cost++
		// Not between floors adds cost +2.
	} else if dir != utilities.DirStop {
		cost += 2
	}

	floor, dir = incrementFloor(floor, dir)
	// Adds cost +2 for each stop and +2 for each travel between adjacent floors.
	for n := 0; !(floor == targetFloor && q.shouldStop(floor, dir)); n++ {
		if q.shouldStop(floor, dir) {
			cost += 2
			q.addOrder(floor, utilities.BtnUp, blankOrder)
			q.addOrder(floor, utilities.BtnDown, blankOrder)
			q.addOrder(floor, utilities.BtnInside, blankOrder)
		}
		dir = q.chooseDirection(floor, dir)
		floor, dir = incrementFloor(floor, dir)
		cost += 2

		if n > 20 {
			break
		}
	}
	return cost
}

func incrementFloor(floor, dir int) (int, int) {
	switch dir {
	case utilities.DirDown:
		floor--
	case utilities.DirUp:
		floor++
	case utilities.DirStop:
		// No incremention.
	default:
		utilities.CloseConnectionChan <- true
		utilities.Restart.Run()
		log.Fatalln("incrementFloor(): invalid direction, not incremented")
	}

	if floor <= 0 && dir == utilities.DirDown {
		dir = utilities.DirUp
		floor = 0
	}
	if floor >= utilities.NFloors-1 && dir == utilities.DirUp {
		dir = utilities.DirDown
		floor = utilities.NFloors - 1
	}
	return floor, dir
}
