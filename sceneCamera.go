package Cameras

//This is a camera library for 3D graphics. package cameralib

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position          mgl32.Vec3
	target            mgl32.Vec3
	up                mgl32.Vec3
	orientation       mgl32.Quat
	mode              int
	groundPlaneNormal mgl32.Vec3
}
//Choose the mode of the camera.
// 1 - Museum mode
// 2 - FPS mode
// 3 - RTS mode
func New(mode int) *Camera {

	c := &Camera{
		position:          mgl32.Vec3{0.0, 0.0, 5.0},
		target:            mgl32.Vec3{0.0, 0.0, 0.0},
		up:                mgl32.Vec3{0.0, 1.0, 0.0},
		orientation:       mgl32.QuatIdent(),
		mode:              mode,
		groundPlaneNormal: mgl32.Vec3{0.0, 0.0, 1.0},
	}
	if mode == 3 {
		c.up = mgl32.Vec3{0.0, 0.0, 1.0}
		c.position = mgl32.Vec3{5.0, 5.0, 5.0}
		//In RTS mode, set the initial target to a point on the ground plane
		forward := c.ForwardsVector()
		c.target = PlaneIntercept(c.groundPlaneNormal, c.position, forward)

	}
	c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	return c
}

func (c *Camera) SetUp(x, y, z float32) {
	c.up = mgl32.Vec3{x, y, z}
}

//Choose the mode of the camera.
// 1 - Museum mode
// 2 - FPS mode
// 3 - RTS mode
func (c *Camera) SetMode(mode int) {
	c.mode = mode
}

//Set the normal of the ground plane.  This is used in RTS mode, and ignored in other modes.
func (c *Camera) SetGroundPlaneNormal(x, y, z float32) {
	c.groundPlaneNormal = mgl32.Vec3{x, y, z}
}

//Print some information about the camera to stdout
func (c *Camera) Dump() {
	fmt.Println("Camera position:", c.position)
	fmt.Println("Camera target:", c.target)
	fmt.Println("Camera rotation:", c.orientation)
	fmt.Println("Camera mode:", c.mode)
	fmt.Println("Forward:", c.ForwardsVector())
	fmt.Println("Right:", c.RightWardsVector())
	fmt.Println("Up:", c.UpwardsVector())
}

// One of the more important functions, LookAt sets the target of the camera.
func (c *Camera) LookAt(x, y, z float32) {
	c.target = mgl32.Vec3{x, y, z}
	c.orientation = mgl32.Mat4ToQuat(mgl32.LookAtV(c.position, c.target, c.up))
}

//Returns the position of the camera in world space
func (c *Camera) Position() (float32, float32, float32) {
	return c.position.X(), c.position.Y(), c.position.Z()
}

//Returns the rotation matrix of the camera.  (the rotation part of the view matrix)
func (c *Camera) RotationMatrix() mgl32.Mat4 {
	return c.orientation.Mat4()
}

/*
func (c *Camera) EulerMatrix() mgl32.Mat4 {
	return c.orientation.Mat4()
}
*/

//Teleport to a position in world space
func (c *Camera) SetPosition(x, y, z float32) {
	c.position = mgl32.Vec3{x, y, z}
}

//Return the ViewMatrix for the camera.  This is the matrix that transforms world space to camera space.  It contains both the rotation and translation of the camera.  It can be passed directly to OpenGL as the ViewMatrix, and used in GLSL shaders as the ViewMatrix.
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	rotation := c.orientation.Mat4()
	translation := mgl32.Translate3D(-c.position.X(), -c.position.Y(), -c.position.Z())
	return rotation.Mul4(translation)
}

// Support 3D displays, by returning the view matrix for the left eye
func (c *Camera) LeftEyeViewMatrix() mgl32.Mat4 {
	ipd := float32(1.0)
	rightVec := c.RightWardsVector()
	eyepos := c.position.Sub(rightVec.Mul(ipd))
	rotation := c.orientation.Mat4()
	translation := mgl32.Translate3D(-eyepos.X(), -eyepos.Y(), -eyepos.Z())
	return rotation.Mul4(translation)
}
// Support 3D displays, by returning the view matrix for the left eye
func (c *Camera) RightEyeViewMatrix() mgl32.Mat4 {
	ipd := float32(1.0)
	rightVec := c.RightWardsVector()
	eyepos := c.position.Add(rightVec.Mul(ipd))
	rotation := c.orientation.Mat4()
	translation := mgl32.Translate3D(-eyepos.X(), -eyepos.Y(), -eyepos.Z())
	return rotation.Mul4(translation)
}

//Reset the camera to its initial position
func (c *Camera) Reset() {
	c.position = mgl32.Vec3{0.0, 0.0, 5.0}
	c.target = mgl32.Vec3{0.0, 0.0, 0.0}
	viewMatrix := mgl32.LookAtV(c.position, c.target, c.up)
	c.orientation = mgl32.Mat4ToQuat(viewMatrix)
	c.groundPlaneNormal = mgl32.Vec3{0.0, 0.0, 1.0}
}

// Move the camera, according to the parameter
// 0 - forward
// 1 - backward
// 2 - left
// 3 - right
// 4 - up
// 5 - down
// 6 - pitch up
// 7 - pitch down
// 8 - yaw left
// 9 - yaw right
// 10 - roll left
// 11 - roll right
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

//Move the camera through world space
func (c *Camera) Translate(x, y, z float32) {
	c.position = c.position.Add(mgl32.Vec3{x, y, z})
}

//Rotate the camera, probably not around the axes that you want
func (c *Camera) Rotate(x, y, z float32) {
	quatX := mgl32.QuatRotate(x, mgl32.Vec3{1, 0, 0})
	quatY := mgl32.QuatRotate(y, mgl32.Vec3{0, 1, 0})
	quatZ := mgl32.QuatRotate(z, mgl32.Vec3{0, 0, 1})
	c.orientation = c.orientation.Mul(quatX).Mul(quatY).Mul(quatZ)
}

func (c *Camera) moveMuseumMode(direction int, amount float32) {
	forward := c.ForwardsVector()
	relativePosition := c.position.Sub(c.target)

	switch direction {
	case 0: // Zoom in
		c.position = c.position.Add(forward.Mul(amount))
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 1: // Zoom out
		c.position = c.position.Sub(forward.Mul(amount))
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 2: // Orbit left
		//Rotate the camera around the target by the specified amount

		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		fmt.Printf("new_relative_position: %v\n", new_relative_position)
		c.position = c.target.Add(new_relative_position.Vec3())
		fmt.Printf("c.position: %v\n", c.position)
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 3: // Orbit right
		//Rotate the camera around the target by the specified amount

		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 4: //Orbit up

		new_relative_position := mgl32.HomogRotate3DX(-amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 5: // Orbit down

		new_relative_position := mgl32.HomogRotate3DX(amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())

	case 6: // Pitch up (Not applicable in museum mode)
	case 7: // Pitch down (Not applicable in museum mode)
	case 8: // Yaw left (Not applicable in museum mode)
	case 9: // Yaw right (Not applicable in museum mode)
	case 10: // Roll left (Not applicable in museum mode)
	case 11: // Roll right (Not applicable in museum mode)
	}
}

// The forward unit vector of the camera, in world space
func (c *Camera) ForwardsVector() mgl32.Vec3 {
	toTarget := c.target.Sub(c.position)
	forward := toTarget.Normalize()
	return forward
}

// The right unit vector of the camera, in world space
func (c *Camera) RightWardsVector() mgl32.Vec3 {
	toTarget := c.target.Sub(c.position).Normalize()
	forward := toTarget
	right := forward.Cross(c.up).Normalize()
	return right
}

// The up unit vector of the camera, in world space
func (c *Camera) UpwardsVector() mgl32.Vec3 {
	toTarget := c.target.Sub(c.position).Normalize()
	forward := toTarget
	right := forward.Cross(c.up).Normalize()
	up := right.Cross(forward).Normalize()
	return up
}

// Scenecam keeps an invisible target point to which the camera is always looking.  Not normalised.  This is the vector from the camera to the target.
// This is not the object that the camera is following
func (c *Camera) TargetVector() mgl32.Vec3 {
	toTarget := c.target.Sub(c.position)
	return toTarget
}

// The position of the target, in world space.
// This is not the object that the camera is following
func (c *Camera) TargetPosition() mgl32.Vec3 {
	return c.target
}

func (c *Camera) moveFPSMode(direction int, amount float32) {
	toTarget := c.TargetVector()
	forward := c.ForwardsVector()
	right := c.RightWardsVector()
	up := c.UpwardsVector()

	switch direction {
	case 0: // Move forward
		c.position = c.position.Add(forward.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 1: // Move backward
		c.position = c.position.Sub(forward.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 2: // Strafe left
		c.position = c.position.Sub(right.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 3: // Strafe right
		c.position = c.position.Add(right.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 4: // Move up
		c.position = c.position.Add(up.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 5: // Move down
		c.position = c.position.Sub(up.Mul(amount))
		c.target = c.position.Add(toTarget)
	case 6: // Pitch up
		//Rotate target around the camera's right vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, right).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	case 7: // Pitch down
		//Rotate target around the camera's right vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, right).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	case 8: // Yaw left
		//Rotate target around the camera's up vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, up).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	case 9: // Yaw right
		//Rotate target around the camera's up vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, up).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	case 10: // Roll left
		//Rotate target around the camera's forward vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, forward).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	case 11: // Roll right
		//Rotate target around the camera's forward vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, forward).Mul4x1(toTarget.Vec4(0))
		c.target = c.position.Add(newTarget.Vec3())
	}
	c.orientation = mgl32.Mat4ToQuat(mgl32.LookAtV(c.position, c.target, c.up))
}

//Project a vector onto a plane, given the normal of the plane, where v1 is the normal of the plane and v2 is the vector to be projected
func ProjectPlane(v1, v2 mgl32.Vec3) mgl32.Vec3 {
	//Project v2 onto v1
	//v1 is the normal of the plane
	//v2 is the vector to be projected
	v1 = v1.Normalize()
	v2 = v2.Normalize()
	d := v1.Dot(v2)
	return v2.Sub(v1.Mul(d))
}
//Find the point on the plane that the ray intercepts
	//groundNormal is the normal of the plane
	//rayOrigin is the origin of the ray
	//rayDirection is the direction of the ray
	//Returns the point on the plane that the ray intercepts
	//
	//The plane is assumed to pass through the origin
func PlaneIntercept(groundNormal, rayOrigin, rayDirection mgl32.Vec3) mgl32.Vec3 {
	
	groundNormal = groundNormal.Normalize()
	rayDirection = rayDirection.Normalize()
	d := groundNormal.Dot(rayDirection)
	if d == 0 {
		//Ray is parallel to the plane
		return mgl32.Vec3{0, 0, 0}
	}
	t := -groundNormal.Dot(rayOrigin) / d
	return rayOrigin.Add(rayDirection.Mul(t))
}

//Find the point on the plane that the ray intercepts
// as for PlaneIntercept, but the plane is not assumed to pass through the origin
func PlaneIntercept2(groundOrigin, groundNormal, rayOrigin, rayDirection mgl32.Vec3) mgl32.Vec3 {
	//Find the point on the plane that the ray intercepts
	//groundOrigin is a point on the plane
	//groundNormal is the normal of the plane
	//rayOrigin is the origin of the ray
	//rayDirection is the direction of the ray
	//Returns the point on the plane that the ray intercepts
	groundNormal = groundNormal.Normalize()
	rayDirection = rayDirection.Normalize()
	d := groundNormal.Dot(rayDirection)
	if d == 0 {
		//Ray is parallel to the plane
		return mgl32.Vec3{0, 0, 0}
	}
	t := groundNormal.Dot(groundOrigin.Sub(rayOrigin)) / d
	return rayOrigin.Add(rayDirection.Mul(t))
}

func (c *Camera) moveRTSMode(direction int, amount float32) {
	forward := c.ForwardsVector()
	up := c.UpwardsVector()

	//Project the camera's forward vector onto the ground plane, held in c.groundPlaneNormal
	groundForwardVec := ProjectPlane(c.groundPlaneNormal, forward).Normalize()
	groundRightVec := up.Cross(groundForwardVec).Normalize()
	target := PlaneIntercept(c.groundPlaneNormal, c.position, forward)
	c.target = target
	//Camera position relative to the target, in this case the ground intercept point
	relativePosition := c.position.Sub(c.target)

	switch direction {
	case 0: // Pan forward
		c.position = c.position.Add(groundForwardVec.Mul(amount))

	case 1: // Pan backward
		c.position = c.position.Sub(groundForwardVec.Mul(amount))
		c.target = c.position.Add(forward)

	case 2: // Pan left
		c.position = c.position.Add(groundRightVec.Mul(amount))
		c.target = c.position.Add(forward)

	case 3: // Pan right
		c.position = c.position.Sub(groundRightVec.Mul(amount))
		c.target = c.position.Add(forward)

	case 10: // Roll left (Not applicable in RTS mode)
	case 11: // Roll right (Not applicable in RTS mode)

	case 4: // Zoom in
		c.position = c.position.Add(forward.Mul(amount))
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
		c.target = c.position.Add(forward)
	case 5: // Zoom out
		c.position = c.position.Sub(forward.Mul(amount))
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
		c.target = c.position.Add(forward)
	case 8: // Orbit left
		//Rotate the camera around the target by the specified amount
		//FIXME rotate around the ground plane normal, not the axis
		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		fmt.Printf("new_relative_position: %v\n", new_relative_position)
		c.position = c.target.Add(new_relative_position.Vec3())
		fmt.Printf("c.position: %v\n", c.position)
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 9: // Orbit right
		//Rotate the camera around the target by the specified amount
		//FIXME rotate around the ground plane normal, not the axis
		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())

	case 6: //Orbit up
		//FIXME rotate around the camera's right vector, not the axis
		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	case 7: // Orbit down
		//FIXME rotate around the camera's right vector, not the axis
		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		c.position = c.target.Add(new_relative_position.Vec3())
		c.LookAt(c.target.X(), c.target.Y(), c.target.Z())
	}

}
