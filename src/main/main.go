package main

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"main/game"
	"main/graphics"
	"main/utils"
	"runtime"
	"main/Common"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	program = graphics.NewGlfwProgram("World", Common.WIDTH, Common.HEIGHT, draw, update)
	defer program.Terminate()

	cam = utils.NewOrthographicCamera(2*Common.SIZE*Common.WIDTH/Common.HEIGHT, 2*Common.SIZE)
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
