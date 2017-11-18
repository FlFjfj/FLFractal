package main

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	//"github.com/go-gl/mathgl/mgl32"
	"main/game"
	"main/graphics"
	"main/utils"
	"runtime"
)

const (
	WIDTH  = 800.0
	HEIGHT = 600.0
)

func init() {
	runtime.LockOSThread()
}

func main() {
	program = graphics.NewGlfwProgram("World", WIDTH, HEIGHT, draw, update)
	defer program.Terminate()

	cam = utils.NewOrthographicCamera(game.SIZE*WIDTH/HEIGHT, game.SIZE)
	cam.ZOOM = 2
	world = game.NewWorld(cam)

	for !program.Window.ShouldClose() {
		program.Update()
	}
}

var (
	program graphics.Program
	cam   utils.OrthographicCamera
	world game.World
)

func update(delta float32) {
	cam.Update()
	world.Update(delta, *program.Window)
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	world.Draw(cam.Combined())

}
