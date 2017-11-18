package game

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"main/graphics"
	"main/utils"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type World struct {
	cam utils.OrthographicCamera

	shader  graphics.ShaderProgram
	object  utils.GameObject
	projLoc int32

	factory CircleFactory
	circles []Circle
}

func NewWorld(camera utils.OrthographicCamera) World {
	worldShader := graphics.NewShaderProgram("assets/shader/worldVert.glsl", "assets/shader/worldFrag.glsl")
	transformLoc := gl.GetUniformLocation(uint32(worldShader), gl.Str("u_ObjTrans\x00"))
	projectionLoc := gl.GetUniformLocation(uint32(worldShader), gl.Str("u_ProjTrans\x00"))
	circle := utils.NewMesh(utils.IdentCircle(100))
	object := utils.NewObject(&circle, SIZE, transformLoc)

	result := World{
		camera,
		worldShader,
		object,
		projectionLoc,
		NewFactory(),
		make([]Circle, 0),
	}
	result.circles = append(result.circles, NewCircle(true, SIZE / 4, mgl32.Vec2{0.0 ,0.0}, mgl32.Vec2{0.0, 0.0}, mgl32.Vec3{1.0, 0.0, 0.0}))

	return result
}

func (world *World) Update(delta float32, window glfw.Window) {
	world.factory.Update(delta, &world.circles)

	for i := 0; i < len(world.circles); i++ {
		world.circles[i].Update(delta)
	}

	world.processMouse(window)
}

func (world *World) Draw(worldTrans mgl32.Mat4) {
	world.shader.Begin()
	gl.UniformMatrix4fv(world.projLoc, 1, false, &worldTrans[0])
	world.object.Draw()
	world.shader.End()

	for _, circle := range world.circles {
		circle.Draw(worldTrans)
	}
}


var(
	lastMouseState = false
	circle *Circle = nil
)
func (world *World) processMouse(window glfw.Window) {
	state := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
	if state && !lastMouseState {
		lastMouseState = true

	} else if !state && lastMouseState {
		lastMouseState = false
	}
}
