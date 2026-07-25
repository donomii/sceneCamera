package sceneCamera

import (
	"fmt"
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

const testEpsilon = 1e-5

func assertVec3(t *testing.T, actual, expected mgl32.Vec3) {
	t.Helper()
	if !actual.ApproxEqualThreshold(expected, testEpsilon) {
		t.Errorf("expected vector %v, got %v", expected, actual)
	}
}

func assertMat4(t *testing.T, actual, expected mgl32.Mat4) {
	t.Helper()
	if !actual.ApproxEqualThreshold(expected, testEpsilon) {
		t.Errorf("expected matrix %v, got %v", expected, actual)
	}
}

func assertFiniteMat4(t *testing.T, matrix mgl32.Mat4) {
	t.Helper()
	for index, value := range matrix {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			t.Errorf("matrix element %d is not finite: %v", index, value)
		}
	}
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func TestNew(t *testing.T) {
	testCases := []struct {
		mode     int
		position mgl32.Vec3
		up       mgl32.Vec3
	}{
		{mode: 1, position: mgl32.Vec3{0, 0, 5}, up: mgl32.Vec3{0, 1, 0}},
		{mode: 2, position: mgl32.Vec3{0, 0, 5}, up: mgl32.Vec3{0, 1, 0}},
		{mode: 3, position: mgl32.Vec3{5, 5, 5}, up: mgl32.Vec3{0, 0, 1}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("mode-%d", testCase.mode), func(t *testing.T) {
			camera := New(testCase.mode)
			if camera.Mode != testCase.mode {
				t.Errorf("expected mode %d, got %d", testCase.mode, camera.Mode)
			}
			assertVec3(t, camera.Position, testCase.position)
			assertVec3(t, camera.Target, mgl32.Vec3{0, 0, 0})
			assertVec3(t, camera.Up, testCase.up)
		})
	}
}

func TestCameraSettersAndReset(t *testing.T) {
	camera := New(1)
	camera.SetUp(0, 0, 1)
	camera.SetMode(2)
	camera.SetGroundPlaneNormal(0, 1, 0)
	camera.SetPosition(2, 3, 4)
	camera.SetIPD(0.8)
	camera.SetFocalLength(12)

	assertVec3(t, camera.Up, mgl32.Vec3{0, 0, 1})
	if camera.Mode != 2 {
		t.Errorf("expected mode 2, got %d", camera.Mode)
	}
	assertVec3(t, camera.GroundPlaneNormal, mgl32.Vec3{0, 1, 0})
	if camera.IPD != 0.8 {
		t.Errorf("expected IPD 0.8, got %v", camera.IPD)
	}
	if camera.FocalLength != 12 {
		t.Errorf("expected focal length 12, got %v", camera.FocalLength)
	}

	x, y, z := camera.WorldPosition()
	if x != 2 || y != 3 || z != 4 {
		t.Errorf("expected position (2, 3, 4), got (%v, %v, %v)", x, y, z)
	}

	camera.Reset()
	assertVec3(t, camera.Position, mgl32.Vec3{0, 0, 5})
	assertVec3(t, camera.Target, mgl32.Vec3{0, 0, 0})
	assertVec3(t, camera.GroundPlaneNormal, mgl32.Vec3{0, 0, 1})
}

func TestTranslate(t *testing.T) {
	camera := New(1)
	camera.Translate(1, 2, 3)
	assertVec3(t, camera.Position, mgl32.Vec3{1, 2, 8})
}

func TestLookAtAndMatrices(t *testing.T) {
	camera := New(2)
	camera.LookAt(1, 0, 0)
	assertVec3(t, camera.Target, mgl32.Vec3{1, 0, 0})

	expectedView := mgl32.LookAtV(camera.Position, camera.Target, camera.Up)
	assertMat4(t, camera.ViewMatrix(), expectedView)
	assertMat4(t, camera.RotationMatrix(), camera.Orientation.Mat4())

	leftView := camera.LeftEyeViewMatrix()
	rightView := camera.RightEyeViewMatrix()
	if leftView.ApproxEqualThreshold(rightView, testEpsilon) {
		t.Error("expected left and right eye view matrices to differ")
	}

	camera.SetIPD(0)
	assertMat4(t, camera.LeftEyeViewMatrix(), camera.ViewMatrix())
	assertMat4(t, camera.RightEyeViewMatrix(), camera.ViewMatrix())
}

func TestRotate(t *testing.T) {
	camera := New(2)
	original := camera.Orientation
	camera.Rotate(0.1, 0.2, 0.3)
	if camera.Orientation.ApproxEqual(original) {
		t.Error("expected Rotate to change the orientation")
	}
}

func TestFrustumMatrices(t *testing.T) {
	camera := New(2)
	left := camera.LeftEyeFrustum()
	right := camera.RightEyeFrustum()
	assertFiniteMat4(t, left)
	assertFiniteMat4(t, right)
	if left.ApproxEqualThreshold(right, testEpsilon) {
		t.Error("expected left and right eye frustum matrices to differ")
	}
	assertMat4(t, camera.LeftEyeFrustrum(), left)
	assertMat4(t, camera.RightEyeFrustrum(), right)
}

func TestFrustumValidation(t *testing.T) {
	testCases := []struct {
		name   string
		clear  func(*Camera)
		method func(*Camera) mgl32.Mat4
	}{
		{name: "right-height", clear: func(camera *Camera) { camera.Screenheight = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "right-width", clear: func(camera *Camera) { camera.Screenwidth = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "right-ipd", clear: func(camera *Camera) { camera.IPD = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "right-near", clear: func(camera *Camera) { camera.Near = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "right-far", clear: func(camera *Camera) { camera.Far = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "right-fov", clear: func(camera *Camera) { camera.FOV = 0 }, method: (*Camera).RightEyeFrustum},
		{name: "left-height", clear: func(camera *Camera) { camera.Screenheight = 0 }, method: (*Camera).LeftEyeFrustum},
		{name: "left-width", clear: func(camera *Camera) { camera.Screenwidth = 0 }, method: (*Camera).LeftEyeFrustum},
		{name: "left-ipd", clear: func(camera *Camera) { camera.IPD = 0 }, method: (*Camera).LeftEyeFrustum},
		{name: "left-near", clear: func(camera *Camera) { camera.Near = 0 }, method: (*Camera).LeftEyeFrustum},
		{name: "left-far", clear: func(camera *Camera) { camera.Far = 0 }, method: (*Camera).LeftEyeFrustum},
		{name: "left-fov", clear: func(camera *Camera) { camera.FOV = 0 }, method: (*Camera).LeftEyeFrustum},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			camera := New(2)
			testCase.clear(camera)
			assertPanics(t, func() { testCase.method(camera) })
		})
	}
}

func TestDirectionVectors(t *testing.T) {
	camera := New(2)
	assertVec3(t, camera.ForwardsVector(), mgl32.Vec3{0, 0, -1})
	assertVec3(t, camera.RightWardsVector(), mgl32.Vec3{1, 0, 0})
	assertVec3(t, camera.UpwardsVector(), mgl32.Vec3{0, 1, 0})
	assertVec3(t, camera.TargetVector(), mgl32.Vec3{0, 0, -5})
	assertVec3(t, camera.TargetPosition(), mgl32.Vec3{0, 0, 0})
}

func TestMuseumMovement(t *testing.T) {
	for direction := 0; direction <= 11; direction++ {
		t.Run(fmt.Sprintf("direction-%d", direction), func(t *testing.T) {
			camera := New(1)
			original := camera.Position
			originalDistance := camera.Position.Sub(camera.Target).Len()
			camera.Move(direction, 0.25)

			switch {
			case direction == 0:
				assertVec3(t, camera.Position, mgl32.Vec3{0, 0, 4.75})
			case direction == 1:
				assertVec3(t, camera.Position, mgl32.Vec3{0, 0, 5.25})
			case direction >= 2 && direction <= 5:
				if camera.Position.ApproxEqualThreshold(original, testEpsilon) {
					t.Error("expected orbit to move the camera")
				}
				if math.Abs(float64(camera.Position.Sub(camera.Target).Len()-originalDistance)) > testEpsilon {
					t.Error("expected orbit to preserve distance from target")
				}
			default:
				assertVec3(t, camera.Position, original)
			}
		})
	}
}

func TestFPSMovement(t *testing.T) {
	positions := map[int]mgl32.Vec3{
		0: {0, 0, 4},
		1: {0, 0, 6},
		2: {-1, 0, 5},
		3: {1, 0, 5},
		4: {0, 1, 5},
		5: {0, -1, 5},
	}
	for direction := 0; direction <= 11; direction++ {
		t.Run(fmt.Sprintf("direction-%d", direction), func(t *testing.T) {
			camera := New(2)
			originalTarget := camera.Target
			originalUp := camera.Up
			camera.Move(direction, 0.25)

			if expected, ok := positions[direction]; ok {
				expected = New(2).Position.Add(expected.Sub(mgl32.Vec3{0, 0, 5}).Mul(0.25))
				assertVec3(t, camera.Position, expected)
			} else {
				assertVec3(t, camera.Position, mgl32.Vec3{0, 0, 5})
			}
			if direction >= 6 && direction <= 9 && camera.Target.ApproxEqualThreshold(originalTarget, testEpsilon) {
				t.Error("expected pitch or yaw to change the target")
			}
			if direction >= 10 {
				assertVec3(t, camera.Target, originalTarget)
				assertVec3(t, camera.Up, originalUp)
			}
		})
	}
}

func TestRTSMovement(t *testing.T) {
	for direction := 0; direction <= 11; direction++ {
		t.Run(fmt.Sprintf("direction-%d", direction), func(t *testing.T) {
			camera := New(3)
			original := camera.Position
			camera.Move(direction, 0.25)

			if direction == 10 || direction == 11 {
				assertVec3(t, camera.Position, original)
			} else if camera.Position.ApproxEqualThreshold(original, testEpsilon) {
				t.Error("expected movement to change the camera position")
			}
			assertFiniteMat4(t, camera.ViewMatrix())
		})
	}
}

func TestPlaneGeometry(t *testing.T) {
	projected := ProjectPlane(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{1, 1, 1})
	if math.Abs(float64(projected.Z())) > testEpsilon {
		t.Errorf("expected projected vector on XY plane, got %v", projected)
	}

	intercept := PlaneIntercept(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{1, 2, 3}, mgl32.Vec3{0, 0, -1})
	assertVec3(t, intercept, mgl32.Vec3{1, 2, 0})
	parallel := PlaneIntercept(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{1, 2, 3}, mgl32.Vec3{1, 0, 0})
	assertVec3(t, parallel, mgl32.Vec3{0, 0, 0})

	intercept = PlaneIntercept2(
		mgl32.Vec3{0, 0, 2},
		mgl32.Vec3{0, 0, 1},
		mgl32.Vec3{1, 2, 3},
		mgl32.Vec3{0, 0, -1},
	)
	assertVec3(t, intercept, mgl32.Vec3{1, 2, 2})
	parallel = PlaneIntercept2(
		mgl32.Vec3{0, 0, 2},
		mgl32.Vec3{0, 0, 1},
		mgl32.Vec3{1, 2, 3},
		mgl32.Vec3{1, 0, 0},
	)
	assertVec3(t, parallel, mgl32.Vec3{0, 0, 0})
}

func ExampleCamera() {
	camera := New(2)
	camera.Move(0, 1)
	x, y, z := camera.WorldPosition()
	fmt.Printf("%.0f %.0f %.0f\n", x, y, z)
	// Output:
	// 0 0 4
}
