package graphics

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"time"
)

type Program struct {
	Window *glfw.Window
	draw   func()
	update func(delta float32)
	last   time.Time
}

func NewGlfwProgram(title string, width, height int, draw func(), update func(delta float32)) Program {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	Window, err := glfw.CreateWindow(width, height, title, glfw.GetPrimaryMonitor(), nil)
	if err != nil {
		panic(err)
	}

	Window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	//gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)


	return Program{Window, draw, update, time.Now()}
}

func (program *Program) Terminate() {
	glfw.Terminate()
}

func (program *Program) Update() {
	program.update(float32(time.Now().Sub(program.last).Nanoseconds()) / 1000000000.0)
	program.last = time.Now()
	program.draw()
	program.Window.SwapBuffers()
	glfw.PollEvents()
}
