package main

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"main/graphics"
	"main/utils"
	"runtime"
)

const (
	WIDTH    = 800.0
	HEIGHT   = 600.0
)

func init() {
	runtime.LockOSThread()
}

func main() {
	program := graphics.NewGlfwProgram("World", WIDTH, HEIGHT, draw, update)
	defer program.Terminate()

	sh = graphics.NewShaderProgram("assets/shader/vert.glsl", "assets/shader/frag.glsl")

	cube = utils.NewMesh(utils.GenCube(1))
	tex = graphics.GetTexture("assets/texture/cube.png")

	cam = utils.NewCamera(WIDTH/HEIGHT, 0.1, 1000.0)
	cam.SetPosition(mgl32.Vec3{0, 3, -3})
	cam.SetTarget(mgl32.Vec3{0, 0, 0})
	cam.SetTop(mgl32.Vec3{0, 1, 0})
	cam.Update()

	for !program.Window.ShouldClose() {
		program.Update()
	}
}

var(
	sh graphics.ShaderProgram
	cam utils.Camera
	tex graphics.Texture
	cube utils.Mesh
	angle float32 = 0
)

func update(delta float32){
	angle += delta
}

func draw(){
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		sh.Begin()

		combined := cam.Combined()
		rotation := rotate(angle)
		gl.UniformMatrix4fv(gl.GetUniformLocation(uint32(sh), gl.Str("u_ProjTrans\x00")), 1, false, &combined[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(uint32(sh), gl.Str("u_ObjTrans\x00")), 1, false, &rotation[0])
		gl.Uniform3f(gl.GetUniformLocation(uint32(sh), gl.Str("u_LightPos\x00")), 0.0, 2.0, -1.0)
		tex.Bind(0)

		cube.Draw()

		tex.Unbind(0)

		sh.End()
}

func rotate(angle float32) mgl32.Mat4 {

	rotate := mgl32.HomogRotate3D(angle * 3.14, mgl32.Vec3{0, 1, 0})
	return rotate
	//return mgl32.HomogRotate3DX(0.75).Mul4(rotate)
}
