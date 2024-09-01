package main

//go:generate go mod init github.com/donomii/splash-screen
//go:generate go mod tidy

import (
	"context"
	"flag"
	"fmt"
	"embed"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"time"
	joystick "github.com/donomii/sceneCamera/joystick"
	messages "github.com/donomii/sceneCamera/messages"

	"github.com/donomii/goof"

	"github.com/mattn/go-shellwords"


	//"time"

	"github.com/go-gl/mathgl/mgl32"

	Cameras "github.com/donomii/sceneCamera"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	_ "embed"
)

//go:embed assets/*
var embeddedFS embed.FS

var MainWin *glfw.Window
var WantSBS bool

var (
	// ... other variables ...
	cameraMode    int
	camera        *Cameras.Camera
	switchModeKey glfw.Key = glfw.KeyTab // Default key to switch camera mode
)

// Arrange that main.main runs on main thread.
func init() {
	flag.BoolVar(&WantSBS, "sbs", false, "Side by side 3D")
	flag.IntVar(&cameraMode, "camera-mode", 1, "Set initial camera mode (1: Museum, 2: FPS, 3: RTS)")
	flag.Parse()
	runtime.LockOSThread()
	debug.SetGCPercent(-1)
}

type State struct {
	prop           int32
	Program        uint32
	Vao            uint32
	Vbo            uint32
	Texture        uint32
	TextureBank    []uint32
	TextureUniform int32
	VertAttrib     uint32
	Angle          float64
	PreviousTime   float64
	ModelUniform   int32
	TexCoordAttrib uint32
}

type tree_struct struct {
	X float32
	Y float32
	Z float32
}

var winWidth = 180
var winHeight = 180
var lasttime float64
var trees []tree_struct

var launchShellList arrayFlags
var launchList arrayFlags
var runningProcs []context.CancelFunc

func drainChannel(ch chan []byte) {
	for {
		<-ch
	}
}


func switchCameraMode() {
	cameraMode = (cameraMode % 3) + 1
	camera.SetMode(cameraMode)
	log.Printf("Switched to camera mode: %d", cameraMode)
}


func main() {
	flag.Var(&launchShellList, "launchShell", "Run shell command at start.  May be repeated to launch multiple commands.")
	flag.Var(&launchList, "launch", "Command line to start an app.  May be repeated to launch multiple apps.")
	flag.Parse()
	for _, commandStr := range launchShellList {
		ctx, cancel := context.WithCancel(context.Background())
		command := exec.CommandContext(ctx, "/bin/sh", "-c", commandStr)
		_, out, err := goof.WrapCmd(command, 3)
		runningProcs = append(runningProcs, cancel)
		go drainChannel(out)
		go drainChannel(err)
	}


	camera = Cameras.New(cameraMode)
	camera.SetPosition(12, 14, 2)
	camera.SetUp(0, 0, 1)
	camera.SetIPD(1.0)
	camera.LookAt(0, 0, 2)
	PI := 3.1415927
	camera.FOV = float32(PI / 4)
	camera.Near = 1.0
	camera.Far = 100
	camera.IPD = 1.0

	currentDir, _ := os.Getwd()
	for _, commandStr := range launchList {
		log.Printf("Launching %v", commandStr)
		os.Chdir(currentDir)
		args, _ := shellwords.Parse(commandStr)
		directory := filepath.Dir(args[0])
		os.Chdir(directory)
		exe := "./" + filepath.Base(args[0])
		log.Printf("Exe %v", exe)
		log.Printf("In dir %v", directory)
		ctx, cancel := context.WithCancel(context.Background())
		command := exec.CommandContext(ctx, exe, args[1:]...)
		_, out, err := goof.WrapCmd(command, 3)
		runningProcs = append(runningProcs, cancel)
		go drainChannel(out)
		go drainChannel(err)
	}
	os.Chdir(currentDir)
	log.Println("Starting windowing system")
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(winWidth, winHeight, "Demo", nil, nil)
	if err != nil {
		panic(err)
	}

	MainWin = win
	go func() {
		time.Sleep(5 * time.Second)
		win.Iconify()
	}()

	win.MakeContextCurrent()

	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		log.Printf("Got key %c,%v,%v,%v", key, key, mods, action)
		handleKey(w, key, scancode, action, mods)
	})

	win.SetMouseButtonCallback(handleMouseButton)

	win.SetCursorPosCallback(handleMouseMove)

	win.SetScrollCallback(handleMouseWheel)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	state := &State{
		prop: 1,
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	state.Program, err = newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	glInit(state, winWidth, winHeight, embeddedFS	)

	//Position some trees
	trees = make([]tree_struct, 0)
	for i := 0; i < 10; i++ {
		//make random location between -10 and 10
		x := rand.Float32()*20 - 10
		y := rand.Float32()*20 - 10
		z := rand.Float32()*20 - 10
		trees = append(trees, tree_struct{X: x, Y: y, Z: z})
	}

	go func() {
		for {
			scale := float32(0.5)  //Multiplier for the camera movement
			//Move the camera
			MoveStep(win, scale)
			time.Sleep(25 * time.Millisecond)
		}
	}()
	joystick.Setup_joystick()
	messages.Register("JoystickY", "JoystickY", func(name , id string, args interface{}) {
		amount := args.(float64)
		camera.Move(0, float32(amount))
	})
	messages.Register("JoystickX", "JoystickX", func(name , id string, args interface{}) {
		amount := args.(float64)
		camera.Move(2, float32(amount))
	})
	for !win.ShouldClose() {
		joystick.DoJoystick()
		mode := glfw.GetPrimaryMonitor().GetVideoMode()
		screenW, screenH := mode.Width, mode.Height
		if screenW >= screenH*2-1 {
			//SBS display mode
			if !WantSBS {
				WantSBS = true
				fmt.Printf("SBS mode %v %v\n", screenW, screenH)
				win.Maximize()
			}
		} else {
			if WantSBS {
				//Normal display mode
				WantSBS = false
				win.Restore()
			}
		}

	
		gfxMain(win, state)
		glfw.PollEvents()
		time.Sleep(1 * time.Millisecond)
		glfw.PollEvents()

	}
	shutdown()

}

func shutdown() {
	for x, cancelF := range runningProcs {
		log.Printf("Stopping sub-process %v", x)
		cancelF()
	}
}

func gfxMain(win *glfw.Window, state *State) {

	//fmt.Println("Draw")
	//width, height := win.GetSize()
	//gl.Viewport(0, 0, int32(width-1), int32(height-1))

	// Render

	// Update

	now := glfw.GetTime()
	elapsed := now - state.PreviousTime

	if elapsed > 0.050 && 1 != win.GetAttrib(glfw.Iconified) {
		//gl.ClearColor(0.0, 1.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		//fmt.Printf("elapsed: %v\n", elapsed)
		state.PreviousTime = now
		angle := state.Angle
		angle += elapsed
		state.Angle = angle

		viewMatrix := camera.ViewMatrix()

		if WantSBS {
			RenderStereoFrame(state, viewMatrix)
		} else {
			projectionMatrix := mgl32.Perspective(camera.FOV, float32(camera.Screenwidth)/float32(camera.Screenheight), camera.Near, camera.Far)
			RenderFrame(state, viewMatrix, projectionMatrix)
		}
		win.SwapBuffers()

	}

}

func RenderStereoFrame(state *State, discard mgl32.Mat4) {
	if MouseWheelValue == 0 {
		MouseWheelValue = 0.01
	}
	camera.SetIPD(MouseWheelValue)
	//get window width and height
	width, height := MainWin.GetSize()
	camera.Screenwidth = float32(width) / 2
	camera.Screenheight = float32(height)
	// Set viewport to left half of window
	gl.Viewport(0, 0, int32(width/2), int32(height))
	LeftviewMatrix := camera.LeftEyeViewMatrix()
	//fmt.Println("Left Eye View Matrix", LeftviewMatrix)
	LeftEyeFrustrum := camera.LeftEyeFrustrum()
	RenderFrame(state, LeftviewMatrix, LeftEyeFrustrum)
	//Set viewport to right half of window
	gl.Viewport(int32(width/2), 0, int32(width/2), int32(height))
	viewMatrix := camera.RightEyeViewMatrix()
	//fmt.Println("Right Eye View Matrix", viewMatrix)
	RightProjectionMatrix := camera.RightEyeFrustrum()
	RenderFrame(state, viewMatrix, RightProjectionMatrix)
	//Set viewport to whole window
	gl.Viewport(0, 0, int32(width), int32(height))
}

func RenderFrame(state *State, viewMatrix mgl32.Mat4, projectionMatrix mgl32.Mat4) {
	cameraUniform := gl.GetUniformLocation(state.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

	//Set a default projection matrix
	projectionUniform := gl.GetUniformLocation(state.Program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])

	//Set the texture to use
	gl.Uniform1i(gl.GetUniformLocation(state.Program, gl.Str("tex\x00")), 0)
	//Bind the texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, state.Texture)

	gl.Disable(gl.BLEND)

	//Draw the ground layer
	for i := -10; i < 11; i++ {
		for j := -10; j < 11; j++ {

			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(float32(i)*2, float32(j)*2, 0))

			gl.UniformMatrix4fv(state.ModelUniform, 1, false, &model[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
		}
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, state.TextureBank[1])

	PI := float32(3.14159265358979323846)

	//Sort the trees by distance from camera
	sort.Slice(trees, func(i, j int) bool {
		//Subtract the camera position from the tree position
		treeVeci := mgl32.Vec3{trees[i].X, trees[i].Y, trees[i].Z}
		cameraVeci := mgl32.Vec3{camera.Position.X(), camera.Position.Y(), camera.Position.Z()}
		treeDisVeci := treeVeci.Sub(cameraVeci)
		treeDisi := treeDisVeci.Len()

		treeVecj := mgl32.Vec3{trees[j].X, trees[j].Y, trees[j].Z}
		cameraVecj := mgl32.Vec3{camera.Position.X(), camera.Position.Y(), camera.Position.Z()}
		treeDisVecj := treeVecj.Sub(cameraVecj)
		treeDisj := treeDisVecj.Len()

		return treeDisi > treeDisj
	})
	// Draw the trees
	for _, tree := range trees {
		model := mgl32.Ident4()

		model = model.Mul4(mgl32.Translate3D(tree.X, tree.Y, 2.0))

		model = model.Mul4(mgl32.HomogRotate3DY(PI))

		gl.UniformMatrix4fv(state.ModelUniform, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 2*3)
	}
}

func checkGlError() {

	err := gl.GetError()
	if err > 0 {
		errStr := fmt.Sprintf("GLerror: %v\n", err)
		fmt.Printf(errStr)
		panic(errStr)
	}

}
