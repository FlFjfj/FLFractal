package main

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	//"github.com/go-gl/mathgl/mgl32"
	"main/game"
	"main/graphics"
	"main/utils"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	program = graphics.NewGlfwProgram("World", game.WIDTH, game.HEIGHT, draw, update)
	defer program.Terminate()

	cam = utils.NewOrthographicCamera(2*game.SIZE*game.WIDTH/game.HEIGHT, 2*game.SIZE)
	world = game.NewWorld(cam)

	for !program.Window.ShouldClose() {
		program.Update()
	}
}

var (
	program graphics.Program
	cam     utils.OrthographicCamera
	world   game.World
)

func update(delta float32) {
	cam.Update()
	world.Update(delta, *program.Window)
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	world.Draw(cam.Combined())

}
