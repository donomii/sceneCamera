package Cameras

//This is a camera library for 3D graphics. package cameralib

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position    mgl32.Vec3
	target      mgl32.Vec3
	up          mgl32.Vec3
	orientation mgl32.Quat
	mode        int
}

func New(mode int) *Camera {

	c := &Camera{
		position:    mgl32.Vec3{0.0, 0.0, 5.0},
		target:      mgl32.Vec3{0.0, 0.0, 0.0},
		up:          mgl32.Vec3{0.0, 1.0, 0.0},
		orientation: mgl32.QuatIdent(),
		mode:        mode,
	}
	viewMatrix := mgl32.LookAtV(c.position, c.target, c.up)
	c.orientation = mgl32.Mat4ToQuat(viewMatrix)
	return c

}

func (c *Camera) Dump() {
	fmt.Println("Camera position:", c.position)
	fmt.Println("Camera target:", c.target)
	fmt.Println("Camera rotation:", c.orientation)
	fmt.Println("Camera mode:", c.mode)
}

func (c *Camera) LookAt(x, y, z float32) {
	c.target = mgl32.Vec3{x, y, z}
}

func (c *Camera) Position() (float32, float32, float32) {
	return c.position.X(), c.position.Y(), c.position.Z()
}

/*
func (c *Camera) Rotation() (float32, float32, float32) {
	return c.orientation.EulerAngles()
}
*/

func (c *Camera) EulerMatrix() mgl32.Mat4 {
	return c.orientation.Mat4()
}

func (c *Camera) SetPosition(x, y, z float32) {
	c.position = mgl32.Vec3{x, y, z}
}

func (c *Camera) ViewMatrix() mgl32.Mat4 {

	return c.orientation.Mat4()
}

func (c *Camera) Reset() {
	c.position = mgl32.Vec3{0.0, 0.0, 5.0}
	c.target = mgl32.Vec3{0.0, 0.0, 0.0}
	c.orientation = mgl32.QuatIdent()
}

func (c *Camera) Move(direction int, amount float32) {
	switch c.mode {
	case 1:
		c.moveMuseumMode(direction, amount)
	case 2:
		c.moveFPSMode(direction, amount)
	case 3:
		c.moveRTSMode(direction, amount)
	}
}

func (c *Camera) Translate(x, y, z float32) {
	c.position = c.position.Add(mgl32.Vec3{x, y, z})
}

func (c *Camera) Rotate(x, y, z float32) {
	quatX := mgl32.QuatRotate(x, mgl32.Vec3{1, 0, 0})
	quatY := mgl32.QuatRotate(y, mgl32.Vec3{0, 1, 0})
	quatZ := mgl32.QuatRotate(z, mgl32.Vec3{0, 0, 1})
	c.orientation = c.orientation.Mul(quatX).Mul(quatY).Mul(quatZ)
}

/*
func (c *Camera) EulerAngles() (float32, float32, float32) {
	return c.Rotation()
}
*/

func (c *Camera) moveMuseumMode(direction int, amount float32) {
	// Assuming c.orientation is a quaternion representing the camera's orientation
	forward := c.orientation.Rotate(mgl32.Vec3{0, 0, 1}).Normalize() 

	switch direction {
	case 0: // Zoom in
		c.position = c.position.Add(forward.Mul(amount))
	
	case 1: // Zoom out
		c.position = c.position.Sub(forward.Mul(amount))
		
	case 2: // Orbit left
		//Rotate the camera around the target by the specified amount
        c.position = c.position.Sub(c.target)
		c.Rotate(0, -amount, 0)
	case 3: // Orbit right
		c.Rotate(0, amount, 0)
	case 4: // Orbit up
		c.Rotate(-amount, 0, 0)
	case 5: // Orbit down
		c.Rotate(amount, 0, 0)
	case 6: // Pitch up (Not applicable in museum mode)
	case 7: // Pitch down (Not applicable in museum mode)
	case 8: // Yaw left (Not applicable in museum mode)
	case 9: // Yaw right (Not applicable in museum mode)
	case 10: // Roll left (Not applicable in museum mode)
	case 11: // Roll right (Not applicable in museum mode)
	}
}

func (c *Camera) moveFPSMode(direction int, amount float32) {
	forward := c.orientation.Rotate(mgl32.Vec3{0, 0, 1}).Normalize() // Rotate the negative z-axis using the camera's orientation
	right := c.orientation.Rotate(mgl32.Vec3{1, 0, 0}).Normalize()    // Rotate the x-axis using the camera's orientation
	up := c.orientation.Rotate(mgl32.Vec3{0, 1, 0}).Normalize()       // Rotate the y-axis using the camera's orientation

	switch direction {
	case 0: // Move forward
		c.position = c.position.Add(forward.Mul(amount))

	case 1: // Move backward
		c.position = c.position.Sub(forward.Mul(amount))

	case 2: // Strafe left
		c.position = c.position.Sub(right.Mul(amount))

	case 3: // Strafe right
		c.position = c.position.Add(right.Mul(amount))

	case 4: // Move up
		c.position = c.position.Add(up.Mul(amount))

	case 5: // Move down
		c.position = c.position.Sub(up.Mul(amount))

	case 6: // Pitch up
		c.Rotate(-amount, 0, 0)
	case 7: // Pitch down
		c.Rotate(amount, 0, 0)
	case 8: // Yaw left
		c.Rotate(0, -amount, 0)
	case 9: // Yaw right
		c.Rotate(0, amount, 0)
	case 10: // Roll left
		c.Rotate(0, 0, -amount)
	case 11: // Roll right
		c.Rotate(0, 0, amount)
	}
}

func (c *Camera) moveRTSMode(direction int, amount float32) {
	forward := c.orientation.Rotate(mgl32.Vec3{0, 0, -1}).Normalize() // Rotate the negative z-axis using the camera's orientation
	right := c.orientation.Rotate(mgl32.Vec3{1, 0, 0}).Normalize()    // Rotate the x-axis using the camera's orientation
	up := c.orientation.Rotate(mgl32.Vec3{0, 1, 0}).Normalize()       // Rotate the y-axis using the camera's orientation

	switch direction {
	case 0: // Pan forward
		c.position = c.position.Add(forward.Mul(amount))
		
	case 1: // Pan backward
		c.position = c.position.Sub(forward.Mul(amount))
	
	case 2: // Pan left
		c.position = c.position.Sub(right.Mul(amount))
		
	case 3: // Pan right
		c.position = c.position.Add(right.Mul(amount))
		
	case 4: // Zoom in
		c.position = c.position.Sub(up.Mul(amount))
	
	case 5: // Zoom out
		c.position = c.position.Add(up.Mul(amount))
	
	case 6: // Rotate up
		c.Rotate(-amount, 0, 0)
	case 7: // Rotate down
		c.Rotate(amount, 0, 0)
	case 8: // Rotate left
		c.Rotate(0, -amount, 0)
	case 9: // Rotate right
		c.Rotate(0, amount, 0)
	case 10: // Roll left (Not applicable in RTS mode)
	case 11: // Roll right (Not applicable in RTS mode)
	}
}
