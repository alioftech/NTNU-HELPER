/*----------------------------------------------------------------------------------
liftSelector package assigns elevators to orders based on cost values calculated
from the cost module in the queue package.
-----------------------------------------------------------------------------------*/
package liftSelector

import (
	"log"
	"queue"
	"time"
	"utilities"
)

type order struct {
	floor  int
	button int
	timer  *time.Timer
}
type reply struct {
	cost int
	lift string
}

// Run receive cost-values from all elevators for each new order, assigning
// a suitable elevator to each order.
func Run(costReply <-chan utilities.Message, numOnline *int) {
	// Gathered cost data for all orders is stored until an elevator is assigned.
	unassigned := make(map[order][]reply)
	var timeout = make(chan *order)

	for {
		select {
		case message := <-costReply:
			newOrder := order{floor: message.Floor, button: message.Button}
			newReply := reply{cost: message.Cost, lift: message.Addr}
			for oldOrder := range unassigned {
				if equal(oldOrder, newOrder) {
					newOrder = oldOrder
				}
			}
			// Check if order exists in queue.
			if replyList, exist := unassigned[newOrder]; exist {
				// Check if newReply already is registered.
				found := false
				for _, reply := range replyList {
					if reply == newReply {
						found = true
					}
				}
				// Not registered -> Add order.
				if !found {
					unassigned[newOrder] = append(unassigned[newOrder], newReply)
					newOrder.timer.Reset(utilities.CostTimeoutTimer)
				}
			} else {
				// If order not in queue, initialise orderlist along with order.
				newOrder.timer = time.NewTimer(utilities.CostTimeoutTimer)
				unassigned[newOrder] = []reply{newReply}
				go costTimer(&newOrder, timeout)
			}
			chooseBestLift(unassigned, numOnline, false)

		case <-timeout:
			log.Println("Time-out collecting costs!")
			chooseBestLift(unassigned, numOnline, true)
		}
	}
}

// chooseBestLift checks if awaiting orders have collected necessary information and are ready
// for lift-assignment. Ready orders are assigned a lift and added to the queue.
func chooseBestLift(unassigned map[order][]reply, numOnline *int, orderTimedOut bool) {
	const maxInt = int(^uint(0) >> 1)
	for order, replyList := range unassigned {

		if len(replyList) == *numOnline || orderTimedOut {
			lowestCost := maxInt
			var bestLift string

			// Cost loop in each complete list.
			for _, reply := range replyList {
				if reply.cost < lowestCost {
					lowestCost = reply.cost
					bestLift = reply.lift
				} else if reply.cost == lowestCost {
					// IP priority if cost is the same (lowest IP first).
					if reply.lift < bestLift {
						lowestCost = reply.cost
						bestLift = reply.lift
					}
				}
			}
			queue.AddRemoteOrder(order.floor, order.button, bestLift)
			order.timer.Stop()
			delete(unassigned, order)
		}
	}
}

func costTimer(newOrder *order, timeout chan<- *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

//equal checks if two orders are the same and returns true or false
func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
}
