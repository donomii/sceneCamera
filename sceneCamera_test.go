package Cameras

import (
	
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	c := New(1)
	if c == nil {
		t.Errorf("Expected a new camera instance, got nil")
	}
}

func TestPosition(t *testing.T) {
	c := New(1)
	x, y, z := c.Position()
	if x != 0 || y != 1 || z != 0 {
		t.Errorf("Expected position (0, 1, 0), got (%f, %f, %f)", x, y, z)
	}
}

func TestSetPosition(t *testing.T) {
	c := New(1)
	c.SetPosition(2, 2, 2)
	x, y, z := c.Position()
	if x != 2 || y != 2 || z != 2 {
		t.Errorf("Expected position (2, 2, 2), got (%f, %f, %f)", x, y, z)
	}
}

/*
func TestRotation(t *testing.T) {
	c := New(1)
	pitch, yaw, roll := c.Rotation()
	const epsilon = 1e-6
	if math.Abs(float64(pitch)) > epsilon || math.Abs(float64(yaw-90)) > epsilon || math.Abs(float64(roll)) > epsilon {
		t.Errorf("Expected rotation (0, 90, 0), got (%f, %f, %f)", pitch, yaw, roll)
	}
}
*/

func TestReset(t *testing.T) {
	c := New(1)
	c.SetPosition(2, 2, 2)
	c.Reset()
	x, y, z := c.Position()
	if x != 0 || y != 1 || z != 0 {
		t.Errorf("Expected reset position (0, 1, 0), got (%f, %f, %f)", x, y, z)
	}
}

func TestTranslate(t *testing.T) {
	c := New(1)
	c.Translate(1, 1, 1)
	x, y, z := c.Position()
	if x != 1 || y != 2 || z != 1 {
		t.Errorf("Expected translated position (1, 2, 1), got (%f, %f, %f)", x, y, z)
	}
}

func TestViewMatrix(t *testing.T) {
	c := New(1)
	viewMatrix := c.ViewMatrix()
	expectedViewMatrix := mgl32.LookAtV(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, -1})

	if !viewMatrix.ApproxEqualThreshold(expectedViewMatrix, 1e-6) {
		t.Errorf("Expected view matrix to match the calculated view matrix")
	}
}

/*
func TestEulerAngles(t *testing.T) {
	c := New(1)
	c.SetPosition(2, 2, 2)
	pitch, yaw, roll := c.EulerAngles()
	const epsilon = 1e-6
	if math.Abs(float64(pitch)) > epsilon || math.Abs(float64(yaw-90)) > epsilon || math.Abs(float64(roll)) > epsilon {
		t.Errorf("Expected rotation (0, 90, 0), got (%f, %f, %f)", pitch, yaw, roll)
	}
}
*/

func TestMove(t *testing.T) {
	const epsilon = 1e-6
	testCases := []struct {
		name      string
		mode      int
		direction int
		amount    float32
		expected  mgl32.Vec3
	}{
		{"Museum mode - forward", 1, 0, 1, mgl32.Vec3{0, 0, -1}},
		{"Museum mode - backward", 1, 1, 1, mgl32.Vec3{0, 2, 0}},
		{"FPS mode - forward", 2, 0, 1, mgl32.Vec3{0, 1, 1}},
		{"FPS mode - left", 2, 2, 1, mgl32.Vec3{1, 1, 0}},
		{"RTS mode - forward", 3, 0, 1, mgl32.Vec3{0, 0, -1}},
		{"RTS mode - up", 3, 4, 1, mgl32.Vec3{0, 2, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New(tc.mode)
			c.Move(tc.direction, tc.amount)
			pos := mgl32.Vec3{c.Position.X(), c.Position.Y(), c.Position.Z()}

			if !pos.ApproxEqualThreshold(tc.expected, epsilon) {
				t.Errorf("Expected position %v, got %v", tc.expected, pos)
			}
		})
	}
}