package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	assets "../assets"
	"log"
	"embed"
)
func glInit( state *State, winWidth, winHeight int, embeddedFS embed.FS) {
	var err error
	//Activate the program we just created.  This means we will use the render and fragment shaders we compiled above
	gl.UseProgram(state.Program)

	//Set a default projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHeight), 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(state.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	/*
	//Setup the camera
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(state.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
*/
	//Setup the cube
	model := mgl32.Ident4()
	state.ModelUniform = gl.GetUniformLocation(state.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(state.ModelUniform, 1, false, &model[0])

	//Find the location of the texture, so we can upload a picture to it
	state.TextureUniform = gl.GetUniformLocation(state.Program, gl.Str("tex\x00"))
	gl.Uniform1i(state.TextureUniform, 0)

	//This is the variable in the fragment shader that will hold the colour for each pixel
	gl.BindFragDataLocation(state.Program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	state.Texture, err = newTexture(assets.GetFileByName(embeddedFS, "assets/grass.png"))
	if err != nil {
		log.Fatalln(err)
	}

	//Load textures into texture bank

	state.TextureBank = make([]uint32, 2)
	for i, textureFile := range []string{"grass.png", "tree.jpg"} {
		log.Printf("Loading texture %v", textureFile)
		//Load an image from a file
		data := assets.GetFileByName(embeddedFS, "assets/"+textureFile)
		state.TextureBank[i], _ = newTexture(data)

	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &state.Vao)
	gl.BindVertexArray(state.Vao)

	gl.GenBuffers(1, &state.Vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, state.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	checkGlError()
	state.VertAttrib = uint32(gl.GetAttribLocation(state.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(state.VertAttrib)
	gl.VertexAttribPointer(state.VertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	checkGlError()
	state.TexCoordAttrib = uint32(gl.GetAttribLocation(state.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(state.TexCoordAttrib)
	gl.VertexAttribPointer(state.TexCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	checkGlError()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)

	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//gl.Enable(gl.TEXTURE_2D)

	gl.UseProgram(state.Program)
	gl.ClearColor(1.0, 1.0, 1.0, 0.0)

	//Activate the cube data, which will be drawn
	gl.BindVertexArray(state.Vao)

	//Choose the texture we just created and uploaded
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, state.Texture)
}