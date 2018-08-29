/*---------------------------------------------------------------------------
hw package contains functions to control the IO ports, as well to two polling functions
(in driver.go) that ping the hardware and registers the buttonpresses of the hardware.
---------------------------------------------------------------------------*/
package hw

import (
	"errors"
	"log"
	"time"
	"utilities"
)

var lampChannelMatrix = [utilities.NFloors][utilities.NButtons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var buttonChannelMatrix = [utilities.NFloors][utilities.NButtons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

type Keypress struct {
	Button int
	Floor  int
}

// Init initialises the lift hardware and moves the lift to a defined state.
// (Descending until it reaches a floor.)
func Init() (int, error) {
	// Init hardware
	if !ioInit() {
		return -1, errors.New("hardware driver: ioInit() failed")
	}

	clearLights()
	// Move to defined state (down)
	SetMotorDir(utilities.DirDown)
	floor := Floor()
	for floor == -1 {
		floor = Floor()
	}
	SetMotorDir(utilities.DirStop)
	SetFloorLamp(floor)
	log.Println("Hardware initialised.")

	return floor, nil
}

//SetMotorDir sets the motor directions (up, down or stop)
func SetMotorDir(dirn int) {
	if dirn == 0 {
		ioWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR, utilities.MotorSpeed)
	} else if dirn < 0 {
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR, utilities.MotorSpeed)
	}
}

//SetDoorLamp sets the lamp ON/OFF
func SetDoorLamp(value bool) {
	if value {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

//setStopLamp write to stopLamp IO
func setStopLamp(value bool) {
	if value {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

//SetFloorLamp sets the lamp of the selected floor (0-3)
func SetFloorLamp(floor int) {
	if floor < 0 || floor >= utilities.NFloors {
		//TODO: ADD ERROR LOG HANDLING
		log.Printf("Error: Floor %d out of range!\n", floor)
		log.Println("No floor indicator will be set.")
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

//Floor returns 0,1,2,3 or -1 for inbetween/invalid floor
func Floor() int {
	if ioReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

//readButton is used by the polling channels, returns true or false (ButtonPressed or NOT)
func readButton(floor int, button int) bool {
	if floor < 0 || floor >= utilities.NFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= utilities.NButtons {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == utilities.BtnUp && floor == utilities.NFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == utilities.BtnDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	if ioReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

//SetButtonLamp sets the Button lamp (UP or Down)
func SetButtonLamp(floor int, button int, value bool) {
	if floor < 0 || floor >= utilities.NFloors {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return
	}
	if button == utilities.BtnUp && floor == utilities.NFloors-1 {
		log.Println("Button up from top floor does not exist!")
		return
	}
	if button == utilities.BtnDown && floor == 0 {
		log.Println("Button down from ground floor does not exist!")
		return
	}
	if button != utilities.BtnUp &&
		button != utilities.BtnDown &&
		button != utilities.BtnInside {
		log.Printf("Invalid button %d\n", button)
		return
	}

	if value {
		ioSetBit(lampChannelMatrix[floor][button])
	} else {
		ioClearBit(lampChannelMatrix[floor][button])
	}
}

//PollButtons has a go-routine that countiously checks the pressed Buttons
//The Buttons are reccieved by the coordinator package
func PollButtons() <-chan Keypress {
	c := make(chan Keypress)
	go func() {
		var buttonState [utilities.NFloors][utilities.NButtons]bool

		for {
			for f := 0; f < utilities.NFloors; f++ {
				for b := 0; b < utilities.NButtons; b++ {
					if (f == 0 && b == utilities.BtnDown) ||
						(f == utilities.NFloors-1 && b == utilities.BtnUp) {
						continue
					}
					if readButton(f, b) {
						if !buttonState[f][b] {
							c <- Keypress{Button: b, Floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(utilities.PollDelayButtons) //TODO HOWMANY? POLLING DELAY.
		}
	}()
	return c
}

//PollFloors has a go-routine that countiously checks for which floor the elevator is on
//The Floor is reccieved by the coordinator package
func PollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := Floor()
		for {
			newFloor := Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(utilities.PollDelayButtons)
		}
	}()
	return c
}

// clearLights turns all the Buttons and Lamps off
func clearLights() {
	for f := 0; f < utilities.NFloors; f++ {
		if f != 0 {
			SetButtonLamp(f, utilities.BtnDown, false)
		}
		if f != utilities.NFloors-1 {
			SetButtonLamp(f, utilities.BtnUp, false)
		}
		SetButtonLamp(f, utilities.BtnInside, false)
	}
	SetDoorLamp(false)
	setStopLamp(false)
}

//getObstructionSignal Reads the IO (not implemented) returns True/false
func getObstructionSignal() bool {
	return ioReadBit(OBSTRUCTION)
}

//getStopSignal Reads the IO(not implemented) returns True/false
func getStopSignal() bool {
	return ioReadBit(STOP)
}
