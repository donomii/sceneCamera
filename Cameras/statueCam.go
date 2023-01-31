package Cameras

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// StatueCamera implements a "statue viewer" camera for 3D graphics.
// The camera looks at a target point in space and can be rotated around it or moved in space.

// StatueCamera holds the camera information
type StatueCamera struct {
	target   mgl32.Vec3 // Target position in world coordinates
	radius   float32    // Distance from camera to target
	rotation mgl32.Vec3 // Euler rotation angles around target in radians
	UpVector mgl32.Vec3 // Up vector for the camera
}

// New creates a new camera with default values
func New() *StatueCamera {
	return &StatueCamera{
		target:   mgl32.Vec3{0, 0, 0},
		radius:   0.5,
		rotation: mgl32.Vec3{0, 0, 0},
		UpVector: mgl32.Vec3{0, 1, 0},
	}
}

// Dump prints camera details to stdout
func (s *StatueCamera) Dump() {
	fmt.Println("Target: ", s.target)
	fmt.Println("Radius: ", s.radius)
	fmt.Println("Rotation: ", s.rotation)
	fmt.Println("Up Vector: ", s.UpVector)
}

// LookAt sets the target point for the camera to look at
func (s *StatueCamera) LookAt(x, y, z float32) {
	//Get position in world coordinates
	px, py, pz := s.Position()
	//Set the target
	s.target = mgl32.Vec3{x, y, z}
	//Set the position
	s.SetPosition(px, py, pz)
}

// Position returns the camera position in world coordinates
func (s *StatueCamera) Position() (float32, float32, float32) {
	res:= s.target.Add( mgl32.SphericalToCartesian(s.radius, s.rotation[0], s.rotation[1]))
	return res[0], res[1], res[2]
}

// returns the Euler rotation angles in radians
func (s *StatueCamera) EulerAngles() (float32, float32, float32) {
	return s.rotation[0], s.rotation[1], s.rotation[2]
}

// Set the angles in radians
func (s *StatueCamera) SetEulerAngles(pitch, yaw, roll float32) {
	s.rotation[0] = pitch
	s.rotation[1] = yaw
	s.rotation[2] = roll
}

func (s *StatueCamera) SetPosition(x, y, z float32) {
	
	v := mgl32.Vec3{x - s.target[0], y - s.target[1], z - s.target[2]}
	r, ele,azi :=mgl32.CartesianToSpherical(v)

	s.radius = r
	s.rotation[0] = ele
	s.rotation[1] = azi
}

func (s *StatueCamera) ViewMatrix() mgl32.Mat4 {
	x, y, z := s.Position()
	position := mgl32.Vec3{x, y, z}
	view := mgl32.LookAtV(position, s.target, s.UpVector)
	return view
}

// Translate moves the camera by x, y, and z units in the world space.
func (s *StatueCamera) Translate(x, y, z float32) {
	//Transform s.radius and s.rotation into cartesian coordinates
	v := mgl32.SphericalToCartesian(s.radius, s.rotation[0], s.rotation[1])
	//Add x,y,z
	v = v.Add(mgl32.Vec3{x, y, z})
	//Convert back to polar coordinates
	r, ele, azi := mgl32.CartesianToSpherical(v)
	s.radius = r
	s.rotation[0] = ele
	s.rotation[1] = azi
}

// RotateZ rotates the camera by `a` radians around the Z axis.
func (s *StatueCamera) RotateZ(a float32) {
	s.rotation[2] += a
}

// RotateX rotates the camera by `a` radians around the X axis.
func (s *StatueCamera) RotateX(a float32) {
	s.rotation[0] += a
}

// RotateY rotates the camera by `a` radians around the Y axis.
func (s *StatueCamera) RotateY(a float32) {
	s.rotation[1] += a
}

func (s *StatueCamera) SetRadius(a float32) {
	s.radius = a
}
