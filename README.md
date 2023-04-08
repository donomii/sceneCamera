SceneCamera
===========

Easy 3D camera management

Description
===========

SceneCamera manages the camera view for 3D applications.  It provides the V in the GL trinity MVP.  It is designed to work with OpenGL, but does not rely on any graphics libraries (just a matrix lib).  It could be used with other graphics libraries, but you may have to copy the viewmatrix into your graphics library's matrix format.


It comes with 3 convenient modes, museum mode, first person mode, and RTS (Real Time Strategy) mode.

### Museum mode
Museum mode is a simple camera that orbits around a point, and can zoom in and out.  Draw your object at the origin, and the camera will orbit around it.  It is useful for inspecting 3D models from all angles.

### First person mode / flight mode

The classic game mode.  Move forwards, backwards, strafe and turn.  It can also pitch, roll, and yaw, making it useful for flight simulators.

### RTS mode

The camera floats over a ground plane.  Forwards, backwards and strafe slide the camera along the ground plane.  It can also spin around a point on the the map, using the pitch and yaw.

## Default position

The camera starts at z=5 (0,0,5), looking at the origin(0,0,0).  Positive Y is up (0,1,0).

In RTS mode, the camera starts at (10,10,10), looking at the origin(0,0,0).  Positive **Z** is up (0,0,1).  The ground plane is at z=0, covering the x and y axes.  i.e. the ground normal vector is (0,0,1).

All these settings can be changed by calling the appropriate method.

## Typical use

    "github.com/donomii/sceneCamera"

    //Create a new camera
    camera := Cameras.New(2)  //FPS mode

    camera.Move(0,0.5)  //Move forwards 0.5 world units

    viewMatrix := camera.ViewMatrix()  //Get the view matrix for the camera

	cameraUniform := gl.GetUniformLocation(state.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

    // Draw your scene

## Examples

An example program is included in the examples directory.  It can be run with

    cd examples
    go run .

## Known issues

RTS mode is not complete, it does not rotate correctly.