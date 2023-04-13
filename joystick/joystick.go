package joystick

import (
	messages "../messages"
	"fmt"
	"github.com/donomii/goof"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Global flag to switch on steam deck features, such as correct button mapping
var steamDeck = false

// Default to a games controller
var joystick_type = "F310"
var latches = make(map[int]bool)
var joystickHistory [][]float64
var joystickHistoryIndex int

func Setup_joystick() {

	//Attempt to detect steam deck
	if goof.Exists("/sys/devices/virtual/dmi/id/board_vendor") {

		if goof.FileContains("/sys/devices/virtual/dmi/id/board_vendor", "Valve") ||
			goof.FileContains("/sys/devices/virtual/dmi/id/board_vendor", "valve") {
			steamDeck = true
			joystick_type = "steamdeck"
		}
	}

}

func outOfDeadZone(deadZone, xdeflect, ydeflect float64) bool {

	return xdeflect > deadZone || xdeflect < -deadZone || ydeflect > deadZone || ydeflect < -deadZone
}
func DoJoystick() {
	deadZone := 0.2
	joy := glfw.Joystick(0)
	if joy.Present() {

		//fmt.Println("Joystick 1 present.  Axes: ", joy.GetAxes(), " Buttons: ", joy.GetButtons(), " Hats: ", joy.GetHats(), "Latches: ", latches)

		//Fetch button states
		buttons := joy.GetButtons()

		//Check for button presses
		for i := 0; i < len(buttons); i++ {

			butt := buttons[i]

			//If the button is available

			//If the button is pressed, and it wasn't pressed last time
			if butt == glfw.Press {

				if !latches[i] {
					latches[i] = true
					messages.SendMessage("Button", i)

				}
			} else {
				//Button is released, so reset the latch
				latches[i] = false
			}

		}
		if joystick_type == "steamdeck" {
			for i := 0; i < len(joy.GetAxes()); i = i + 3 {
				var ydeflect float64 = 0
				var trigger float64 = 0
				xdeflect := float64(joy.GetAxes()[i])
				if i+1 < len(joy.GetAxes()) {
					ydeflect = float64(joy.GetAxes()[i+1])
				}
				if i+2 < len(joy.GetAxes()) {
					trigger = float64(joy.GetAxes()[i+2])
				}
				if outOfDeadZone(deadZone, xdeflect, ydeflect) {
					messages.SendMessage("JoystickY", ydeflect/-100)
					messages.SendMessage("JoystickX", xdeflect/-100)
					messages.SendMessage("JoystickTrigger", trigger)
				}
				/*
					if xdeflect > deadZone || xdeflect < -deadZone || ydeflect > deadZone || ydeflect < -deadZone {
						b.X = b.X + (xdeflect)*10.0
						b.Y = b.Y + (ydeflect)*10.0

					}
				*/

			}

			//fmt.Printf("Axes: %+v, hats: %+v, buttons: %+v\n", joy.GetAxes(), joy.GetButtons(), joy.GetHats())
		}

		if joystick_type == "F310" {
			for i := 0; i < len(joy.GetAxes()); i = i + 2 {
				var ydeflect float64 = 0
				var trigger float64 = 0
				xdeflect := float64(joy.GetAxes()[i])
				if i+1 < len(joy.GetAxes()) {
					ydeflect = float64(joy.GetAxes()[i+1])
				}
				if i+2 < len(joy.GetAxes()) {
					trigger = float64(joy.GetAxes()[i+2])
				}
				if outOfDeadZone(deadZone, xdeflect, ydeflect) {
					messages.SendMessage("JoystickY", ydeflect/-100)
					messages.SendMessage("JoystickX", xdeflect/-100)
					messages.SendMessage("JoystickTrigger", trigger)
				}

				/*
					if xdeflect > deadZone || xdeflect < -deadZone || ydeflect > deadZone || ydeflect < -deadZone {
						b.X = b.X + (xdeflect)*10.0
						b.Y = b.Y + (ydeflect)*10.0

					}
				*/

			}
		}
	}
}
