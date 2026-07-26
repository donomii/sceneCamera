package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	demoWidth      = 480
	demoHeight     = 270
	demoFrameCount = 72
	demoFrameDelay = 7
)

func validateDemoMode(mode string) error {
	switch mode {
	case "", "rts", "flight":
		return nil
	default:
		return fmt.Errorf("record demo mode %q is unsupported; expected rts or flight", mode)
	}
}

func makeDemoTrees() []tree_struct {
	positions := [][2]float32{
		{-9, -7}, {-7, -4}, {-5, -8}, {-2, -5}, {1, -8}, {4, -5}, {7, -8}, {9, -4},
		{-9, 0}, {-6, 2}, {-3, 0}, {2, 2}, {5, 0}, {9, 2},
		{-8, 7}, {-4, 6}, {0, 8}, {4, 6}, {8, 8},
	}
	result := make([]tree_struct, 0, len(positions))
	for _, position := range positions {
		result = append(result, tree_struct{X: position[0], Y: position[1]})
	}
	return result
}

func recordDemoGIF(win *glfw.Window, state *State, mode string) error {
	outputPath, err := demoOutputPath(mode)
	if err != nil {
		return err
	}

	gl.ClearColor(0.58, 0.8, 0.98, 1)
	configureDemoCamera(mode, 0)
	renderDemoFrame(win, state)
	win.SwapBuffers()
	glfw.PollEvents()

	animation := &gif.GIF{
		Image:     make([]*image.Paletted, 0, demoFrameCount),
		Delay:     make([]int, 0, demoFrameCount),
		LoopCount: 0,
	}
	for frameIndex := 0; frameIndex < demoFrameCount; frameIndex++ {
		progress := float32(frameIndex) / float32(demoFrameCount-1)
		configureDemoCamera(mode, progress)
		renderDemoFrame(win, state)
		animation.Image = append(animation.Image, captureDemoFrame(demoWidth, demoHeight))
		animation.Delay = append(animation.Delay, demoFrameDelay)
		win.SwapBuffers()
		glfw.PollEvents()
	}

	return writeDemoGIF(outputPath, animation)
}

func demoOutputPath(mode string) (string, error) {
	fileName := "camerademo_" + mode + ".gif"
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("determine current directory for %s demo output: %w", mode, err)
	}
	if filepath.Base(currentDirectory) == "example" {
		return filepath.Join("..", fileName), nil
	}
	return fileName, nil
}

func configureDemoCamera(mode string, progress float32) {
	camera.FOV = mgl32.DegToRad(52)
	camera.Near = 0.1
	camera.Far = 100
	camera.Screenwidth = demoWidth
	camera.Screenheight = demoHeight

	switch mode {
	case "rts":
		configureRTSDemoCamera(progress)
	case "flight":
		configureFlightDemoCamera(progress)
	}
}

func configureRTSDemoCamera(progress float32) {
	angle := progress * 2 * float32(math.Pi)
	mapOffset := mgl32.Vec3{
		6 * float32(math.Cos(float64(angle))),
		6 * float32(math.Sin(float64(angle))),
		0,
	}
	position := mapOffset.Add(mgl32.Vec3{11, -14, 13})
	target := mapOffset.Add(mgl32.Vec3{0, 0, 0.8})

	camera.SetMode(3)
	camera.SetUp(0, 0, 1)
	camera.SetPosition(position.X(), position.Y(), position.Z())
	camera.LookAt(target.X(), target.Y(), target.Z())
}

func configureFlightDemoCamera(progress float32) {
	position := flightPosition(progress)
	lookAheadDistance := 0.22 - 0.17*progress
	lookAheadProgress := min(progress+lookAheadDistance, 1)
	direction := flightPosition(lookAheadProgress).Sub(position)
	if direction.Len() == 0 {
		direction = position.Sub(flightPosition(max(progress-0.035, 0)))
	}
	direction = direction.Normalize()
	bankAngle := 0.32 * float32(math.Sin(float64(progress*2*float32(math.Pi))))
	up := mgl32.QuatRotate(bankAngle, direction).Rotate(mgl32.Vec3{0, 0, 1})
	target := position.Add(direction.Mul(8))

	camera.SetMode(2)
	camera.SetUp(up.X(), up.Y(), up.Z())
	camera.SetPosition(position.X(), position.Y(), position.Z())
	camera.LookAt(target.X(), target.Y(), target.Z())
}

func flightPosition(progress float32) mgl32.Vec3 {
	x := -15 + 27*progress
	y := -13 + 23*progress + 3*float32(math.Sin(float64(progress*2*float32(math.Pi))))
	diveProgress := smoothStep(0, 0.62, progress)
	z := 17 + (4.2-17)*diveProgress
	return mgl32.Vec3{x, y, z}
}

func smoothStep(edgeStart, edgeEnd, value float32) float32 {
	position := min(max((value-edgeStart)/(edgeEnd-edgeStart), 0), 1)
	return position * position * (3 - 2*position)
}

func renderDemoFrame(win *glfw.Window, state *State) {
	width, height := win.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	projection := mgl32.Perspective(camera.FOV, float32(width)/float32(height), camera.Near, camera.Far)
	RenderFrame(state, camera.ViewMatrix(), projection)
	gl.Finish()
}

func captureDemoFrame(width, height int) *image.Paletted {
	pixels := make([]uint8, width*height*4)
	gl.ReadBuffer(gl.BACK)
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	frame := image.NewRGBA(image.Rect(0, 0, width, height))
	rowLength := width * 4
	for row := 0; row < height; row++ {
		sourceStart := (height - row - 1) * rowLength
		destinationStart := row * frame.Stride
		copy(frame.Pix[destinationStart:destinationStart+rowLength], pixels[sourceStart:sourceStart+rowLength])
	}

	palettedFrame := image.NewPaletted(frame.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(palettedFrame, frame.Bounds(), frame, image.Point{})
	return palettedFrame
}

func writeDemoGIF(outputPath string, animation *gif.GIF) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create demo GIF %q: %w", outputPath, err)
	}
	encodeErr := gif.EncodeAll(outputFile, animation)
	closeErr := outputFile.Close()
	if encodeErr != nil {
		return fmt.Errorf("encode demo GIF %q: %w", outputPath, encodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close demo GIF %q: %w", outputPath, closeErr)
	}
	return nil
}
