# SceneCamera

Easy 3D camera management

# Description

SceneCamera manages the camera view for 3D applications, especially for Android and GLES applications. 

It provides the V in the GL trinity MVP.

The  OpenGl removes all the convenience functions for working with 3D cameras, this library puts them back. It also supports headtracking (only tested on Android). You can send your Android events directly to the camera (with ProcessEvent), and SceneCam will handle all the details of maintaining a the correct View Matrix.

# Use

For a golang mobile app

    app.Main(func(a app.App) {
    for e := range a.Events() {
        sensor.Notify(a)
        sc := sceneCamera.New()
        sc.ProcessEvent(e)
        view := sc.ViewMatrix()
        glctx.UniformMatrix4fv(uniform_view, view[0:16])
        
For any other application:

        sc := sceneCamera.New()
        sc.Translate(1.0, 2.0, 3.0)
        view := sc.ViewMatrix()
        glctx.UniformMatrix4fv(uniform_view, view[0:16])

# Notes

The camera starts at position 0,0,0.5, looking at 0,0,0 (the origin). I.e. it is staring directly down the Z axis, in the negative direction. 

Positive Y should always be up, towards the top of  monitor, and positive X should be to the right.

# Bugs

At the moment, the camera rotates around the origin, instead of the origin rotating around the camera.  I'll fix tht soon.
