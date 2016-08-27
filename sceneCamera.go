package sceneCamera

import "github.com/go-gl/mathgl/mgl32"
import "golang.org/x/mobile/exp/sensor"

type SceneCamera struct {
    Camera  mgl32.Mat4
    RotationMatrix mgl32.Mat4  
    PrevTime int64       //Device timestamp.  Does not hold correct time, is only useful for delta time
    flipCam  bool       //Draw the mirror version of a scene
}

//New Creates a new sora and returns it.  The camera is positioned at 0.5 on the Y axis, and is looking directly down the axis in the -Y direction (i.e. it is looking at 0,0,0).
func New () *SceneCamera {
    return myinit()
}

func myinit () *SceneCamera {
    s := SceneCamera{}
    s.Camera = mgl32.LookAt(0.0, 0.0, 0.8, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
    s.RotationMatrix = mgl32.Ident4()
    return &s
}

func (s *SceneCamera) FlipCam ()  {
    //sMat := mgl32.Scale3D(1.0,-1.0,1.0)
    //s.Camera = compose(sMat, s.Camera)
    s.flipCam = !s.flipCam
}


//ViewMatrix returns the camera matrix, aka the "view" matrix, which forms the V part of the MVP matrix for 3D graphics.
func (s *SceneCamera) ViewMatrix() mgl32.Mat4 {
    //sMat := mgl32.Scale3D(1.0,-1.0,1.0)
    //return compose(sMat, s.Camera)     //Fixme
    return s.Camera
}

func (s *SceneCamera) Reset() {
    s = myinit()      //Fixme
}


//Sets the internal view matrix, replacing it with your own mgl32.Mat4
func (s *SceneCamera) SetViewMatrix(newMatrix mgl32.Mat4) {
    s.Camera = newMatrix
}

//Moves the camera
func (s *SceneCamera) Translate(x,y,z float32) {
    s.Camera = compose(s.Camera, mgl32.Translate3D(x,y,z))
}


func compose (a, b mgl32.Mat4) mgl32.Mat4 {
    return a.Mul4(b)
}

//Expects events from the gomobile "app" module.  ProcessEvent will attempt to extract the movement events and process them.
func (s *SceneCamera) ProcessEvent(e sensor.Event) {
                  delta := e.Timestamp - s.PrevTime
                  s.PrevTime = e.Timestamp
                  scale := float32(360000000.0/float32(delta))   //Arbitrary scale, works for my phone, not sure if universal?
                  var sora mgl32.Vec3   //The real sora
                  sora = mgl32.Vec3{float32(-e.Data[1])/scale, float32(e.Data[0])/scale,float32(-e.Data[2])/scale/float32(3.14)/2.0}
                  s_norm := sora.Normalize()
                  rotMatrix := mgl32.HomogRotate3D(sora.Len(), s_norm)
                  s.Camera = compose(rotMatrix, s.Camera)
}
