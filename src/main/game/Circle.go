package game

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"main/graphics"
	"main/utils"
)

var (
	isInit       = false
	circleShader graphics.ShaderProgram
	objLoc       int32
	projLoc      int32
	colorLoc     int32
	mesh         utils.Mesh
	acceleration float32 = -1
)

type Circle struct {
	object   utils.GameObject
	owns bool
	position mgl32.Vec2
	velocity mgl32.Vec2
	size     float32
	color    mgl32.Vec3
}

func NewCircle(owns bool, size float32, position mgl32.Vec2, velocity mgl32.Vec2, color mgl32.Vec3) Circle {
	if !isInit {
		println("Init circle")
		isInit = true
		mesh = utils.NewMesh(utils.IdentCircle(20))
		circleShader = graphics.NewShaderProgram("assets/shader/circleVert.glsl", "assets/shader/circleFrag.glsl")
		objLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_ObjTrans\x00"))
		projLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_ProjTrans\x00"))
		colorLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_Color\x00"))
	}

	return Circle{
		utils.NewObject(&mesh, size, objLoc),
		owns,
		position,
		velocity,
		size,
		color,
	}
}

func (circle *Circle) Draw(worldTrans mgl32.Mat4) {
	circleShader.Begin()
	gl.UniformMatrix4fv(projLoc, 1, false, &worldTrans[0])
	gl.Uniform3fv(colorLoc, 1, &circle.color[0])
	circle.object.Draw()
	circleShader.End()
}

func (circle *Circle) Update(delta float32) {
	circle.position = circle.position.Add(circle.velocity.Mul(delta))
	len := circle.velocity.Len()
	if len != 0 {
		if circle.velocity.Len() < acceleration*delta {
			circle.velocity = mgl32.Vec2{0.0, 0.0}
		} else {
			circle.velocity = circle.velocity.Add(circle.velocity.Normalize().Mul(acceleration * delta))
		}
	}

	if circle.position.Len() > SIZE - circle.size {
		circle.position = circle.position.Normalize().Mul(SIZE - circle.size)
		circle.velocity = circle.velocity.Mul(-1)
	}

	circle.object.Update(circle.position, circle.size)
}
