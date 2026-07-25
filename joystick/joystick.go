package joystick

import (
	//messages "github.com/donomii/sceneCamera/messages"

	"github.com/donomii/goof"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Default to a games controller
var joystick_type = "F310"
var latches = make(map[int]bool)

func Setup_joystick() {

	//Attempt to detect steam deck
	if goof.Exists("/sys/devices/virtual/dmi/id/board_vendor") {

		if goof.FileContains("/sys/devices/virtual/dmi/id/board_vendor", "Valve") ||
			goof.FileContains("/sys/devices/virtual/dmi/id/board_vendor", "valve") {
			joystick_type = "steamdeck"
		}
	}

}

func DoJoystick() {
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
					//messages.SendMessage("Button", i)

				}
			} else {
				//Button is released, so reset the latch
				latches[i] = false
			}

		}
	}
}
