package main

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var oldXpos float64

func handleMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	//log.Printf("Mouse moved: %v,%v", xpos, ypos)
	//diff := xpos - oldXpos
	
}
func handleMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	log.Printf("Got mouse button %v,%v,%v", button, mod, action)
	//handleKey(w, key, scancode, action, mods)
}

func handleKey(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if mods == 2 && (action == 1 || action == 2) && key != 341 {  //action 1 is press, 2 is repeat
		mask := ^byte(64 + 128)
		log.Printf("key mask: %#b", mask)
		val := byte(key)
		log.Printf("key val: %#b", val)
		b := mask & val
		log.Printf("key byte: %#b", b)

	}

	if (action == 1 || action == 2) && mods == 0 {  //FIXME use latch
		log.Printf("Acting on  key %c,%v", key, key)
		switch key {
			case 65:
				//Left
				camera.Move(2,0.1)
			case 68:
				//Right
				camera.Move(3,0.1)
			case 87:
				//Up
				camera.Move(0,0.2)
			case 83:
				//Down
				camera.Move(1,0.2)
			case 81:
				//Q
				camera.Move(8,0.1)
			case 69:
				//E
				camera.Move(9,0.1)
			// X
			case 88:
				camera.Move(4,0.1)
			// space bar
			case 32:
				camera.Move(5,0.1)

				


		case 256:
			log.Println("Quitting")
			shutdown()
			win.SetShouldClose(true)

		}
		camera.Dump()
	}

}
