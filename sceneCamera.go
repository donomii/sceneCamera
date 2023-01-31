package Cameras

//This is a camera library for 3D graphics. 
//It implements a "statue viewer" camera.  This camera looks at a point in space, and can be rotated around that point.  It can also be moved around in space, but the camera always looks at the same point.

//Calling the rotate functions will rotate the camera around the point it is looking at.  Calling the translate functions will move the camera around in space, but it will always look at the same point.
import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)


type StatueCamera struct {
}

//New Creates a new camera and returns it.  The camera is positioned at 0.5 on the Y axis, and is looking directly down the axis in the -Y direction (i.e. it is looking at 0,0,0).
func New() *StatueCamera {
	
}

//Create a new camera with sensible defaults
func myinit() *StatueCamera {
}

//Dump prints the camera matrix, rotations and other details to stdout
func (s *StatueCamera) Dump() {

}

//LookAt sets the camera to look at the specified point in world coordinates
func (s *StatueCamera) LookAt(x, y, z float32) {

}

//Return the position of the camera in world coordinates
func (s *StatueCamera) Position() (float32, float32, float32) {

}

//Return the Euler rotation angles in radians
func (s *StatueCamera) Rotation() (float32, float32, float32) {

}



//Return a matrix that will rotate the camera to the specified Euler angles in radians
func (s *StatueCamera) EulerMatrix() mgl32.Mat4 {

}


func (s *StatueCamera) SetPosition(x, y, z float32) {
	
}


//ViewMatrix returns the camera matrix, aka the "view" matrix, which forms the V part of the MVP matrix for 3D graphics.
func (s *StatueCamera) ViewMatrix() mgl32.Mat4 {

}

func (s *StatueCamera) Reset() {
	
}

//Sets the internal view matrix, replacing it with your own mgl32.Mat4
func (s *StatueCamera) SetViewMatrix(newMatrix mgl32.Mat4) {
	
}

//Moves the camera
func (s *StatueCamera) Translate(x, y, z float32) {
}

func (s *StatueCamera) RotateZ(a float32) {


}

/*
// Checks if a matrix is a valid rotation matrix.
func isRotationMatrix(R *mgl32.Mat4) bool {
	Rt := R.Transpose()
	shouldBeIdentity := Rt.Mul4(*R)
	I := eye(3, 3, shouldBeIdentity.Type())

	return norm(I, shouldBeIdentity) < 1e-6

}
*/

func rotationMatrixToEulerAngles(R mgl32.Mat4) mgl32.Vec3 {

	//assert(isRotationMatrix(R));

	sy := Sqrt(R.At(0, 0)*R.At(0, 0) + R.At(1, 0)*R.At(1, 0))

	var x, y, z float32
	if !(sy < 1e-6) {
		x = Atan2(R.At(2, 1), R.At(2, 2))
		y = Atan2(-R.At(2, 0), sy)
		z = Atan2(R.At(1, 0), R.At(0, 0))
	} else {
		x = Atan2(-R.At(1, 2), R.At(1, 1))
		y = Atan2(-R.At(2, 0), sy)
		z = 0
	}
	return mgl32.Vec3{x, y, z}
}

func (s *StatueCamera) RotateX(a float32) {

}

//Rotate around the Y axis
func (s *StatueCamera) RotateY(a float32) {

}

