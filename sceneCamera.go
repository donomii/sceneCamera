package Cameras

//This is a camera library for 3D graphics. package cameralib

import (
	"log"
	"fmt"
	"math"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position          mgl32.Vec3  //The position of the camera in world space
	Target            mgl32.Vec3  //The target of the camera in world space.  Note: not the focal point
	Up                mgl32.Vec3  //The up vector of the camera
	Orientation       mgl32.Quat  //The orientation of the camera, quaternion
	Mode              int	 //The mode of the camera.  1 - Museum mode, 2 - FPS mode, 3 - RTS mode
	GroundPlaneNormal mgl32.Vec3	 //The normal of the ground plane
	IPD 			 float32	 //The inter-pupillary distance, in world space
	FocalLength		 float32 	 //The focal length of the camera, in world space
	Near 			 float32	 //The near clipping plane
	Far 			 float32	 //The far clipping plane
	Screenheight 	float32	 //The height of the screen, in pixels
	Screenwidth 	float32	 //The width of the screen, in pixels
	Aperture 		float32	 //The aperture of the camera, in world space
	FOV 			 float32	 //The field of view of the camera, in degrees
	
}

var PI = float32(3.1415927)
//Choose the mode of the camera.
// 1 - Museum mode
// 2 - FPS mode
// 3 - RTS mode
func New(mode int) *Camera {

	c := &Camera{
		Position:          mgl32.Vec3{0.0, 0.0, 5.0},
		Target:            mgl32.Vec3{0.0, 0.0, 0.0},
		Up:                mgl32.Vec3{0.0, 1.0, 0.0},
		Orientation:       mgl32.QuatIdent(),
		Mode:              mode,
		GroundPlaneNormal: mgl32.Vec3{0.0, 0.0, 1.0},
		FocalLength: 5.0,
		Near:  0.1,
		Far: 30.0,
		FOV : PI/2.0,
		IPD: 2.0,
		Screenheight: 1080.0,
		Screenwidth: 1920.0,
	}
	if mode == 3 {
		c.Up = mgl32.Vec3{0.0, 0.0, 1.0}
		c.Position = mgl32.Vec3{5.0, 5.0, 5.0}
		//In RTS mode, set the initial target to a point on the ground plane
		forward := c.ForwardsVector()
		c.Target = PlaneIntercept(c.GroundPlaneNormal, c.Position, forward)

	}
	c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	return c
}

func (c *Camera) SetUp(x, y, z float32) {
	c.Up = mgl32.Vec3{x, y, z}
}

//Choose the mode of the camera.
// 1 - Museum mode
// 2 - FPS mode
// 3 - RTS mode
func (c *Camera) SetMode(mode int) {
	c.Mode = mode
}

//Set the normal of the ground plane.  This is used in RTS mode, and ignored in other modes.
func (c *Camera) SetGroundPlaneNormal(x, y, z float32) {
	c.GroundPlaneNormal = mgl32.Vec3{x, y, z}
}

//Print some information about the camera to stdout
func (c *Camera) Dump() {
	fmt.Println("Camera position:", c.Position)
	fmt.Println("Camera target:", c.Target)
	fmt.Println("Camera rotation:", c.Orientation)
	fmt.Println("Camera mode:", c.Mode)
	fmt.Println("Forward:", c.ForwardsVector())
	fmt.Println("Right:", c.RightWardsVector())
	fmt.Println("Up:", c.UpwardsVector())
}

// One of the more important functions, LookAt sets the target of the camera.
func (c *Camera) LookAt(x, y, z float32) {
	c.Target = mgl32.Vec3{x, y, z}
	c.Orientation = mgl32.Mat4ToQuat(mgl32.LookAtV(c.Position, c.Target, c.Up))
}

//Returns the position of the camera in world space
func (c *Camera) WorldPosition() (float32, float32, float32) {
	return c.Position.X(), c.Position.Y(), c.Position.Z()
}

//Returns the rotation matrix of the camera.  (the rotation part of the view matrix)
func (c *Camera) RotationMatrix() mgl32.Mat4 {
	return c.Orientation.Mat4()
}

/*
func (c *Camera) EulerMatrix() mgl32.Mat4 {
	return c.orientation.Mat4()
}
*/

//Teleport to a position in world space
func (c *Camera) SetPosition(x, y, z float32) {
	c.Position = mgl32.Vec3{x, y, z}
}

//Set the inter-pupillary distance for 3D displays (in world coordinates)
func (c *Camera) SetIPD(ipd float32) {
	c.IPD = ipd
}

func (c *Camera) SetFocalLength(focalLength float32) {
	c.FocalLength = focalLength
}

//Return the ViewMatrix for the camera.  This is the matrix that transforms world space to camera space.  It contains both the rotation and translation of the camera.  It can be passed directly to OpenGL as the ViewMatrix, and used in GLSL shaders as the ViewMatrix.
func (c *Camera) ViewMatrix() mgl32.Mat4 {
	rotation := c.Orientation.Mat4()
	translation := mgl32.Translate3D(-c.Position.X(), -c.Position.Y(), -c.Position.Z())
	return rotation.Mul4(translation)
}

// Support 3D displays, by returning the view matrix for the left eye
func (c *Camera) LeftEyeViewMatrix() mgl32.Mat4 {
	
	rightVec := c.RightWardsVector()
	eyepos := c.Position.Sub(rightVec.Mul(c.IPD/2))
	rotation := c.Orientation.Mat4()
	log.Printf("Left eye position: %v", eyepos)
	translation := mgl32.Translate3D(-eyepos.X(), -eyepos.Y(), -eyepos.Z())
	log.Printf("Left eye translation matrix: %v", translation)
	return rotation.Mul4(translation)
}
// Support 3D displays, by returning the view matrix for the left eye
func (c *Camera) RightEyeViewMatrix() mgl32.Mat4 {
	
	rightVec := c.RightWardsVector()
	eyepos := c.Position.Add(rightVec.Mul(c.IPD/2))
	rotation := c.Orientation.Mat4()
	log.Printf("Right eye position: %v", eyepos)
	translation := mgl32.Translate3D(-eyepos.X(), -eyepos.Y(), -eyepos.Z())
	log.Printf("Right eye translation matrix: %v", translation)
	return rotation.Mul4(translation)
}

// Calculate the frustrum matrix for the right eye
func (c *Camera) RightEyeFrustrum() mgl32.Mat4 {
	if c.Screenheight == 0 {
		panic("Screen height is zero")
	}
	if c.Screenwidth == 0 {
		panic("Screen width is zero")
	}
	if c.IPD == 0 {
		panic("IPD is zero")
	}
	if c.Near == 0 {
		panic("Near is zero")
	}
	if c.Far == 0 {
		panic("Far is zero")
	}
	if c.FOV == 0 {
		panic("FOV is zero")
	}
	aspect_ratio  := c.Screenwidth / c.Screenheight
	frustumshift := (c.IPD/2)*c.Near/c.Far
	top := c.Near * float32(math.Tan(float64(c.FOV/2)))
	right := aspect_ratio*top+frustumshift
	left := aspect_ratio*top-frustumshift
	bottom := -top
	frustrum := mgl32.Frustum(left, right, bottom, top, c.Near, c.Far)
	return frustrum
}

//Reset the camera to its initial position
func (c *Camera) Reset() {
	c.Position = mgl32.Vec3{0.0, 0.0, 5.0}
	c.Target = mgl32.Vec3{0.0, 0.0, 0.0}
	viewMatrix := mgl32.LookAtV(c.Position, c.Target, c.Up)
	c.Orientation = mgl32.Mat4ToQuat(viewMatrix)
	c.GroundPlaneNormal = mgl32.Vec3{0.0, 0.0, 1.0}
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
	switch c.Mode {
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
	c.Position = c.Position.Add(mgl32.Vec3{x, y, z})
}

//Rotate the camera, probably not around the axes that you want
func (c *Camera) Rotate(x, y, z float32) {
	quatX := mgl32.QuatRotate(x, mgl32.Vec3{1, 0, 0})
	quatY := mgl32.QuatRotate(y, mgl32.Vec3{0, 1, 0})
	quatZ := mgl32.QuatRotate(z, mgl32.Vec3{0, 0, 1})
	c.Orientation = c.Orientation.Mul(quatX).Mul(quatY).Mul(quatZ)
}

func (c *Camera) moveMuseumMode(direction int, amount float32) {
	forward := c.ForwardsVector()
	relativePosition := c.Position.Sub(c.Target)

	switch direction {
	case 0: // Zoom in
		c.Position = c.Position.Add(forward.Mul(amount))
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 1: // Zoom out
		c.Position = c.Position.Sub(forward.Mul(amount))
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 2: // Orbit left
		//Rotate the camera around the target by the specified amount

		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		fmt.Printf("new_relative_position: %v\n", new_relative_position)
		c.Position = c.Target.Add(new_relative_position.Vec3())
		fmt.Printf("c.position: %v\n", c.Position)
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 3: // Orbit right
		//Rotate the camera around the target by the specified amount

		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 4: //Orbit up

		new_relative_position := mgl32.HomogRotate3DX(-amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 5: // Orbit down

		new_relative_position := mgl32.HomogRotate3DX(amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())

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
	toTarget := c.Target.Sub(c.Position)
	forward := toTarget.Normalize()
	return forward
}

// The right unit vector of the camera, in world space
func (c *Camera) RightWardsVector() mgl32.Vec3 {
	toTarget := c.Target.Sub(c.Position).Normalize()
	forward := toTarget
	right := forward.Cross(c.Up).Normalize()
	return right
}

// The up unit vector of the camera, in world space
func (c *Camera) UpwardsVector() mgl32.Vec3 {
	toTarget := c.Target.Sub(c.Position).Normalize()
	forward := toTarget
	right := forward.Cross(c.Up).Normalize()
	up := right.Cross(forward).Normalize()
	return up
}

// Scenecam keeps an invisible target point to which the camera is always looking.  Not normalised.  This is the vector from the camera to the target.
// This is not the object that the camera is following
func (c *Camera) TargetVector() mgl32.Vec3 {
	toTarget := c.Target.Sub(c.Position)
	return toTarget
}

// The position of the target, in world space.
// This is not the object that the camera is following
func (c *Camera) TargetPosition() mgl32.Vec3 {
	return c.Target
}

func (c *Camera) moveFPSMode(direction int, amount float32) {
	toTarget := c.TargetVector()
	forward := c.ForwardsVector()
	right := c.RightWardsVector()
	up := c.UpwardsVector()

	switch direction {
	case 0: // Move forward
		c.Position = c.Position.Add(forward.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 1: // Move backward
		c.Position = c.Position.Sub(forward.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 2: // Strafe left
		c.Position = c.Position.Sub(right.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 3: // Strafe right
		c.Position = c.Position.Add(right.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 4: // Move up
		c.Position = c.Position.Add(up.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 5: // Move down
		c.Position = c.Position.Sub(up.Mul(amount))
		c.Target = c.Position.Add(toTarget)
	case 6: // Pitch up
		//Rotate target around the camera's right vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, right).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	case 7: // Pitch down
		//Rotate target around the camera's right vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, right).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	case 8: // Yaw left
		//Rotate target around the camera's up vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, up).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	case 9: // Yaw right
		//Rotate target around the camera's up vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, up).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	case 10: // Roll left
		//Rotate target around the camera's forward vector by the specified amount
		newTarget := mgl32.HomogRotate3D(amount, forward).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	case 11: // Roll right
		//Rotate target around the camera's forward vector by the specified amount
		newTarget := mgl32.HomogRotate3D(-amount, forward).Mul4x1(toTarget.Vec4(0))
		c.Target = c.Position.Add(newTarget.Vec3())
	}
	c.Orientation = mgl32.Mat4ToQuat(mgl32.LookAtV(c.Position, c.Target, c.Up))
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
	groundForwardVec := ProjectPlane(c.GroundPlaneNormal, forward).Normalize()
	groundRightVec := up.Cross(groundForwardVec).Normalize()
	target := PlaneIntercept(c.GroundPlaneNormal, c.Position, forward)
	c.Target = target
	//Camera position relative to the target, in this case the ground intercept point
	relativePosition := c.Position.Sub(c.Target)

	switch direction {
	case 0: // Pan forward
		c.Position = c.Position.Add(groundForwardVec.Mul(amount))

	case 1: // Pan backward
		c.Position = c.Position.Sub(groundForwardVec.Mul(amount))
		c.Target = c.Position.Add(forward)

	case 2: // Pan left
		c.Position = c.Position.Add(groundRightVec.Mul(amount))
		c.Target = c.Position.Add(forward)

	case 3: // Pan right
		c.Position = c.Position.Sub(groundRightVec.Mul(amount))
		c.Target = c.Position.Add(forward)

	case 10: // Roll left (Not applicable in RTS mode)
	case 11: // Roll right (Not applicable in RTS mode)

	case 4: // Zoom in
		c.Position = c.Position.Add(forward.Mul(amount))
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
		c.Target = c.Position.Add(forward)
	case 5: // Zoom out
		c.Position = c.Position.Sub(forward.Mul(amount))
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
		c.Target = c.Position.Add(forward)
	case 8: // Orbit left
		//Rotate the camera around the target by the specified amount
		//FIXME rotate around the ground plane normal, not the axis
		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		fmt.Printf("new_relative_position: %v\n", new_relative_position)
		c.Position = c.Target.Add(new_relative_position.Vec3())
		fmt.Printf("c.position: %v\n", c.Position)
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 9: // Orbit right
		//Rotate the camera around the target by the specified amount
		//FIXME rotate around the ground plane normal, not the axis
		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())

	case 6: //Orbit up
		//FIXME rotate around the camera's right vector, not the axis
		new_relative_position := mgl32.HomogRotate3DY(-amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	case 7: // Orbit down
		//FIXME rotate around the camera's right vector, not the axis
		new_relative_position := mgl32.HomogRotate3DY(amount).Mul4x1(relativePosition.Vec4(0))
		c.Position = c.Target.Add(new_relative_position.Vec3())
		c.LookAt(c.Target.X(), c.Target.Y(), c.Target.Z())
	}

}
