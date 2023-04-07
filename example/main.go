package main

import (
	Cameras ".."
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	"strings"
)

const (
	width  = 800
	height = 600
)

var vertexShaderSource = `#version 410
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vColor;

uniform mat4 u_MVP;

void main()
{
    gl_Position = u_MVP * vec4(aPos, 1.0);
    vColor = aColor;
}
` + "\x00"

var fragmentShaderSource = `#version 410
in vec3 vColor;

out vec4 FragColor;

void main()
{
    FragColor = vec4(vColor, 1.0);
}

` + "\x00"

func main() {
	runtime.LockOSThread()

	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}
	defer glfw.Terminate()

	// Create GLFW window
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(width, height, "Cube Demo", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create glfw window: %v", err))
	}
	win.MakeContextCurrent()

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("could not initialize OpenGL: %v", err))
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	// Setup OpenGL resources for rendering the cube
	setupOpenGL()

	camera := Cameras.New(1)
	camera.LookAt(0, 0, 0)
	camera.Move(2, 1.0)

	// Main loop
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Handle input
		handleInput(win, camera)

		// Render the cube
		//Get the window size

		windowWidth, windowHeight := win.GetFramebufferSize()
		renderCube(camera, windowWidth, windowHeight)

		// Swap buffers
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func handleInput(win *glfw.Window, camera *Cameras.Camera) {
	// Process input (WASD, mouse to pitch and yaw, and mouse button to switch modes)
}

func renderCube(camera *Cameras.Camera, windowWidth, windowHeight int) {
	// Set up the projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0)

	// Set up the model matrix
	model := mgl32.Ident4()

	// Calculate the MVP matrix
	MVP := projection.Mul4(camera.ViewMatrix()).Mul4(model)

	// Render the cube using the MVP matrix
	gl.UseProgram(shaderProgram)

	//Set the u_MVP uniform in the shader program to the MVP matrix
	gl.ProgramUniformMatrix4x2fv(shaderProgram, mvpLocation, 1, false, &MVP[0])

	gl.DrawArrays(gl.TRIANGLES, 0, 24)

}

var (
	cubeVAO       uint32
	cubeVBO       uint32
	cubeEBO       uint32
	shaderProgram uint32
	mvpLocation   int32
)

func setupOpenGL() {
	// Cube vertices and colors
	cubeVertices := []float32{
		// Positions          // Colors
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, -0.5, 1.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
	}

	cubeIndices := []uint32{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		0, 1, 5, 5, 4, 0,
		2, 3, 7, 7, 6, 2,
		0, 3, 7, 7, 4, 0,
		1, 2, 6, 6, 5, 1,
	}

	// Compile shaders
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	// Create shader program
	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// Get the MVP uniform location
	mvpLocation = gl.GetUniformLocation(shaderProgram, gl.Str("u_MVP\x00"))

	// Create vertex array object (VAO)
	gl.GenVertexArrays(1, &cubeVAO)
	gl.BindVertexArray(cubeVAO)

	// Create vertex buffer object (VBO)
	gl.GenBuffers(1, &cubeVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	// Create element buffer object (EBO)
	gl.GenBuffers(1, &cubeEBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cubeEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cubeIndices)*4, gl.Ptr(cubeIndices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Unbind VAO, VBO, and EBO
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to compile %v: %v", source, log))
	}

	return shader
}
