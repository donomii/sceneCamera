package sceneCamera

import "github.com/go-gl/mathgl/mgl32"
import "golang.org/x/mobile/exp/sensor"

//import "log"

import "fmt"

type SceneCamera struct {
	camera         mgl32.Mat4
	rotationMatrix mgl32.Mat4
	pos            mgl32.Vec4
	PrevTime       int64 //Device timestamp.  Does not hold correct time, is only useful for delta time
	flipCam        bool  //Draw the mirror version of a scene
}

//New Creates a new sora and returns it.  The camera is positioned at 0.5 on the Y axis, and is looking directly down the axis in the -Y direction (i.e. it is looking at 0,0,0).
func New() *SceneCamera {
	return myinit()
}

func myinit() *SceneCamera {
	s := SceneCamera{}
	s.camera = mgl32.LookAt(0.0, 0.0, -0.81, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
	s.pos = mgl32.Vec4{0.0, 0.0, -0.81, 1.0}
	s.rotationMatrix = mgl32.Ident4()
	return &s
}

func (s *SceneCamera) Dump() {
	fmt.Println("Camera matrix:", s.camera)
	x, y, z := s.Position()
	fmt.Println("Position - X: ", x, "Y:", y, "Z:", z)
	fmt.Println("Col 4 ", s.camera.Col(3))
}
func (s *SceneCamera) LookAt(x, y, z float32) {
	//vec := s.RotationMatrix.Mul4x1(mgl32.Vec4{1.0, 1.0, 1.0, 1.0})
	xx, yy, zz := s.Position()
	//fmt.Printf("Eye vector: %v\n", vec)
	s.camera = mgl32.LookAt(xx, yy, zz, x, y, z, 0.0, 1.0, 0.0)
}

func (s *SceneCamera) Position() (float32, float32, float32) {

	//vec := s.Camera.Mul4x1(mgl32.Vec4{0.0, 0.0, 0.0, 1.0})
	//return vec.X(), vec.Y(), vec.Z()
	return s.pos.X(), s.pos.Y(), s.pos.Z()
}

func (s *SceneCamera) SetPosition(x, y, z float32) {
	s.pos = mgl32.Vec4{0.0, 0.0, 0.0, 1.0}
}

func (s *SceneCamera) FlipCam() {
	//sMat := mgl32.Scale3D(1.0,-1.0,1.0)
	//s.Camera = compose(sMat, s.Camera)
	s.flipCam = !s.flipCam
}

//ViewMatrix returns the camera matrix, aka the "view" matrix, which forms the V part of the MVP matrix for 3D graphics.
func (s *SceneCamera) ViewMatrix() mgl32.Mat4 {
	//sMat := mgl32.Scale3D(1.0,-1.0,1.0)
	//return compose(sMat, s.Camera)     //Fixme
	ret := compose(mgl32.Ident4(), s.rotationMatrix)
	ret.SetCol(3, s.pos)
	return ret
	//return s.Camera
}

func (s *SceneCamera) Reset() {
	s = myinit() //Fixme
}

//Sets the internal view matrix, replacing it with your own mgl32.Mat4
func (s *SceneCamera) SetViewMatrix(newMatrix mgl32.Mat4) {
	s.camera = newMatrix
}

//Moves the camera
func (s *SceneCamera) Translate(x, y, z float32) {

	s.camera = compose(s.camera, mgl32.Translate3D(x, y, z))
	//fmt.Println("Translate matrix: ", mgl32.Translate3D(x, y, z))
	s.pos = s.pos.Add(s.rotationMatrix.Mul4x1(mgl32.Vec4{x, y, z, 0.0}))
	fmt.Println("Done")
}

//Rotate around the Y axis
//FIXME translate to the origin, do the rotate, then translate back
//Maybe we should start storing the MVP matrices separately?
func (s *SceneCamera) RotateY(a float32) {
	s.rotationMatrix = compose(s.rotationMatrix, mgl32.HomogRotate3DY(a))
	/*
		fmt.Println("-----------------------")
		fmt.Println("Dump prior to rotate: ")
		s.Dump()
		fmt.Println("-----------------------")
		x, y, z := s.Position()
		positionVec := mgl32.Vec4{x, y, z, 0.0}
		//Find our current lookat target
		targetPos := s.Camera.Mul4x1(mgl32.Vec4{0.0, 0.0, 1.0, 1.0})
		fmt.Println("Unit vector in eyespace", targetPos)
		//Move our position to the origin
		shifted := targetPos.Sub(positionVec)
		fmt.Println("Vector with eyepos removed", shifted)
		//Rotate it
		rotated := mgl32.HomogRotate3DY(a).Mul4x1(shifted)
		fmt.Println("Rotated", rotated)
		//Move it back into position
		newTarget := rotated.Sub(positionVec)
		fmt.Println("New target in world space", newTarget)

		s.Camera = mgl32.LookAt(x, y, z, newTarget.X(), newTarget.Y(), newTarget.Z(), 0.0, 1.0, 0.0)
		fmt.Println("-----------------------")
		fmt.Println("Dump post to rotate: ")
		s.Dump()
		fmt.Println("-----------------------")

			fmt.Println("-----------------------")
			fmt.Println("Dump prior to translate: ")
			s.Dump()
			//s.Translate(x, y, z)
			fmt.Println("Dump post translate: ")

			s.Dump()
			fmt.Println("Translating - X: ", x, "Y:", y, "Z:", z)
			fmt.Println("-----------------------")

			//trans := mgl32.Translate3D(x, y, z)
			//invTrans := trans.Inv()
			col := s.Camera.Col(3)
			s.Camera.SetCol(3, mgl32.Vec4{0.0, 0.0, 0.0, 1.0})
			s.Camera = compose(s.Camera, mgl32.HomogRotate3DY(a))
			vec := s.Camera.Mul4x1(col)
			s.Translate(-vec.X(), -vec.Y(), -vec.Z())
			//s.Camera = compose(s.Camera, invTrans)
			//s.Camera.SetCol(3,col)
	*/
}

func compose(a, b mgl32.Mat4) mgl32.Mat4 {
	return a.Mul4(b)
}

//Expects events from the gomobile "app" module.  ProcessEvent will attempt to extract the movement events and process them.
func (s *SceneCamera) ProcessEvent(e sensor.Event) {
	delta := e.Timestamp - s.PrevTime
	s.PrevTime = e.Timestamp
	scale := float32(360000000.0 / float32(delta)) //Arbitrary scale, works for my phone, not sure if universal?
	var sora mgl32.Vec3                            //The real sora
	sora = mgl32.Vec3{float32(-e.Data[1]) / scale, float32(e.Data[0]) / scale, float32(-e.Data[2]) / scale / float32(3.14) / 2.0}
	s_norm := sora.Normalize()
	rotMatrix := mgl32.HomogRotate3D(sora.Len(), s_norm)
	s.camera = compose(rotMatrix, s.camera)
}
