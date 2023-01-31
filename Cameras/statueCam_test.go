package Cameras

import (
	"testing"
	"math"

	"fmt"
	"github.com/go-gl/mathgl/mgl32"

)

func TestNew(t *testing.T) {
	c := New()
	x,y,z := c.Position()
	if x != 0.0 || y != 0.0 || z != 0.5 {
		fmt.Println("x,y,z = ", x, y, z)
		t.Errorf("Expected position to be 0,0,0.5, but got %v,%v,%v", x, y, z)
	}
	//Test rotation
	//...

}

func TestPositionRoundTrip(t *testing.T) {
	c := New()
	c.SetPosition(1, 2, 3)
	x, y, z := c.Position()
	c.Dump()
	//Compare floats with a tolerance
	if math.Abs(float64(x-1)) > 0.0001 || math.Abs(float64(y-2)) > 0.0001 || math.Abs(float64(z-3)) > 0.0001 {
		t.Errorf("Expected position to be 1,2,3, but got %v,%v,%v", x, y, z)
	}
}

func TestLookAt(t *testing.T) {
	c := New()
	c.LookAt(-4, 0,4)
	c.SetPosition(0,0,0)
	c.LookAt(-4, 0,4)

v := mgl32.SphericalToCartesian(1, mgl32.DegToRad(0), mgl32.DegToRad(45))
fmt.Printf("v = %v\n", v)


}

func TestPosition(t *testing.T) {
	c := New()
	x, y, z := c.Position()
	if x != 0.0 || y != 0.0 || z != 0.5 {
		t.Errorf("Expected position to be 0,0,0.5, but got %v,%v,%v", x, y, z)
	}
}

func TestRotateZ(t *testing.T) {
	c := New()
	c.RotateZ(math.Pi / 2)
	if c.rotation[2] != math.Pi/2 {
		t.Errorf("Expected rotation to be %v, but got %v", math.Pi/2, c.rotation[2])
	}
}

func TestViewMatrix(t *testing.T) {
	c := New()
	c.Translate(1, 2, 3)
	c.RotateX(mgl32.DegToRad(30))
	c.RotateY(mgl32.DegToRad(45))

	c.RotateY(mgl32.DegToRad(-45))
	c.RotateX(mgl32.DegToRad(-30))
	c.Translate(-1, -2, -3)


	view := c.ViewMatrix()

	expected := New().ViewMatrix()

	if view != expected {
		t.Errorf("Expected view matrix to be %v, but got %v", expected, view)
	}
}


