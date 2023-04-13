package main

import (
	"log"
"sync"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var oldXpos float64
var MouseWheelValue float32

func handleMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	//log.Printf("Mouse moved: %v,%v", xpos, ypos)
	//diff := xpos - oldXpos

}
func handleMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	log.Printf("Got mouse button %v,%v,%v", button, mod, action)
	//handleKey(w, key, scancode, action, mods)
}

func handleMouseWheel(w *glfw.Window, xoff float64, yoff float64) {
	log.Printf("Got mouse wheel %v,%v", xoff, yoff)
	MouseWheelValue += float32(yoff)

	//handleKey(w, key, scancode, action, mods)
}

var keyLatch = NewGenericMap[glfw.Key,bool]()

func handleKey(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if mods == 2 && (action == 1 || action == 2) && key != 341 { //action 1 is press, 2 is repeat
		mask := ^byte(64 + 128)
		log.Printf("key mask: %#b", mask)
		val := byte(key)
		log.Printf("key val: %#b", val)
		b := mask & val
		log.Printf("key byte: %#b", b)

	}

	if (action == 1 || action == 2) && mods == 0 { //FIXME use latch
		keyLatch.Set(key, true)
	}

	if action == 0 && mods == 0 {
		keyLatch.Set(key,false)
	}


	camera.Dump()

}


func MoveStep(win *glfw.Window, size float32) {
	
	//Left
	if keyLatch.Get(65) {
		camera.Move(2, 0.2*size)
	}
	//Right
	if keyLatch.Get(68) {
		camera.Move(3, 0.2*size)
	}
	//Up
	if keyLatch.Get(87) {
		camera.Move(0, 0.2*size)
	}
	//Down
	if keyLatch.Get(83) {
		camera.Move(1, 0.2*size)
	}
	//Q
	if keyLatch.Get(81) {
		camera.Move(8, 0.1*size)
	}
	//E
	if keyLatch.Get(69) {
		camera.Move(9, 0.1*size)

	}
	// X
	if keyLatch.Get(88) {
		camera.Move(4, 0.1*size)
	}
	// space bar
	if keyLatch.Get(32) {
		camera.Move(5, 0.1*size)
	}
	// R
	if keyLatch.Get(82) {
		camera.Move(6, 0.1*size)
	}
	// F
	if keyLatch.Get(70){
		camera.Move(7, 0.1*size)
	}
	// Esc
	if keyLatch	.Get(256) {
		log.Println("Quitting")
		shutdown()
		win.SetShouldClose(true)
	}
}

// Make a generic hashmap for key presses

type GenericMap[K comparable, V any] struct {
	M map[K]V
	Mutex sync.Mutex
}

func (m *GenericMap[K, V]) Get(key K) V {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return m.M[key]
}

func (m *GenericMap[K, V]) Load(key K) (V, bool) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	v, ok := m.M[key]
	return v, ok
}

func (m *GenericMap[K, V]) Set(key K, value V) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.M[key] = value
}

func (m *GenericMap[K, V]) Delete(key K) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	delete(m.M, key)
}

func (m *GenericMap[K, V]) Has(key K) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	_, ok := m.M[key]
	return ok
}

func (m *GenericMap[K, V]) Keys() []K {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	keys := make([]K, 0, len(m.M))
	for k := range m.M {
		keys = append(keys, k)
	}
	return keys
}

func (m *GenericMap[K, V]) Values() []V {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	values := make([]V, 0, len(m.M))
	for _, v := range m.M {
		values = append(values, v)
	}
	return values
}

func (m *GenericMap[K, V]) Len() int {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return len(m.M)
}

func (m *GenericMap[K, V]) Clear() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.M = make(map[K]V)
}

func (m *GenericMap[K, V]) Lock() {
	m.Mutex.Lock()
}

func (m *GenericMap[K, V]) Unlock() {
	m.Mutex.Unlock()
}

func NewGenericMap[K comparable, V any]() *GenericMap[K, V] {
	return &GenericMap[K, V]{M: make(map[K]V)}
}

