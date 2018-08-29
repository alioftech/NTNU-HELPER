/*----------------------------------------------------------------------------------
utilities package contains constants, type structs and variables used in other packages.
----------------------------------------------------------------------------------*/
package utilities

import (
	"os/exec"
	"time"
)

// Messaage sent over the network
type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

// Network message category constants
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

const (
	BtnUp int = iota
	BtnDown
	BtnInside
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

// Local IP address
var LocalIP string

// Hardware constants
const MotorSpeed = 2800
const NButtons = 3
const NFloors = 4

//Network constants
const LocalListenPort = 35555
const BroadcastListenPort = 35557

//Timers
const PollDelayButtons = 100 * time.Millisecond  //inside driver
const LiftAliveTimeout = 500 * time.Millisecond  //inside cordinator
const AliveHeartbeat = 400 * time.Millisecond    //inside network
const DoorOpenInterval = 2500 * time.Millisecond //inside fsm
const OrderTimeout = 10000 * time.Millisecond    //inside queue- methods
const CostTimeoutTimer = 300 * time.Millisecond  //inside liftSelector (

var CloseConnectionChan = make(chan bool)

// Start a new terminal when Restart.Run() is called
var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
