/*
	coordinator package coordinate between the different fsm Channels
	depending on the type of event. It's the brain of the program.
*/
package coordinator

import (
	"fsm"
	"hw"
	"log"
	"network"
	"queue"
	"time"
	"utilities"
)

//Init starts the brains of the program
func Init(ch fsm.Channels) {
	go coordinator(ch)
}

//coordinator is the brains of the progam, it coordinates all of the fsm channels to the right method
func coordinator(ch fsm.Channels) {
	buttonChan := hw.PollButtons()
	floorChan := hw.PollFloors()

	for {
		select {
		case keyPress := <-buttonChan:
			switch keyPress.Button {
			case utilities.BtnInside:
				queue.AddLocalOrder(keyPress.Floor, keyPress.Button)
			case utilities.BtnUp, utilities.BtnDown:
				ch.OutgoingMsg <- utilities.Message{Category: utilities.NewOrder, Floor: keyPress.Floor, Button: keyPress.Button}
			}
		case floor := <-floorChan:
			ch.FloorReached <- floor
		case ListenMsg := <-ch.IncomingMsg:
			msgCategoryHandler(ListenMsg, ch)
		case dead := <-ch.DeadChan:
			DeadConnnectionHandler(dead.Addr, ch)
		case order := <-queue.OrderTimeoutChan:
			log.Println("Order timeout, I'll handle it :)")
			queue.RemoveRemoteOrdersAt(order.Floor)
			queue.AddRemoteOrder(order.Floor, order.Button, utilities.LocalIP)
		case motorDir := <-ch.MotorDir:
			hw.SetMotorDir(motorDir)
		case value := <-ch.DoorLamp:
			hw.SetDoorLamp(value)
		case floorLamp := <-ch.FloorLamp:
			hw.SetFloorLamp(floorLamp)
		}
	}
}

//msgCategoryHandler handles the connection messages, Lift status: Alive, NewOrder, CompleteOrder, cost msg
func msgCategoryHandler(msg utilities.Message, ch fsm.Channels) {
	switch msg.Category {
	case utilities.Alive:
		if connection, exist := ch.OnlineLifts[msg.Addr]; exist {
			connection.Timer.Reset(utilities.LiftAliveTimeout)
		} else {
			newConnection := network.UdpConnection{msg.Addr, time.NewTimer(utilities.LiftAliveTimeout)}
			ch.OnlineLifts[msg.Addr] = newConnection
			ch.NOnline = len(ch.OnlineLifts)
			go connectionTimer(&newConnection, ch)
			log.Printf("Connection to IP %s established!", msg.Addr[0:15])
		}
	case utilities.NewOrder:
		cost := queue.CalculateCost(msg.Floor, msg.Button, fsm.Floor(), hw.Floor(), fsm.Direction())
		ch.OutgoingMsg <- utilities.Message{Category: utilities.Cost, Floor: msg.Floor, Button: msg.Button, Cost: cost}
	case utilities.Cost:
		ch.CostChan <- msg
	case utilities.CompleteOrder:
		queue.RemoveRemoteOrdersAt(msg.Floor)
	}
}

//connectionTimer waits for connetion.Timer.C channels if recieved it sends a msg to the coordinator
func connectionTimer(connection *network.UdpConnection, ch fsm.Channels) {
	<-connection.Timer.C
	ch.DeadChan <- *connection
}

//DeadConnnectionHandler reassigns the order to another elevator and deleted the elevator from the OnlineLifts map
func DeadConnnectionHandler(deadAddr string, ch fsm.Channels) {
	log.Printf("Connection to IP %s is dead!", deadAddr[0:15])
	delete(ch.OnlineLifts, deadAddr)
	ch.NOnline = len(ch.OnlineLifts)
	queue.ReassignOrders(deadAddr, ch.OutgoingMsg)
}
