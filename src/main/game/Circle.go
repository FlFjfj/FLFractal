package game

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"main/graphics"
	"main/utils"
	"math"
)

var (
	isInit       = false
	circleShader graphics.ShaderProgram
	objLoc       int32
	projLoc      int32
	colorLoc     int32
	mesh         utils.Mesh
	acceleration float32 = -7
)

type Circle struct {
	id       int
	object   utils.GameObject
	owner    int
	position mgl32.Vec2
	velocity mgl32.Vec2
	size     float32
	color    mgl32.Vec3
}

func NewCircle(id int, owner int, size float32, position mgl32.Vec2, velocity mgl32.Vec2, color mgl32.Vec3) *Circle {
	if !isInit {
		isInit = true
		mesh = utils.NewMesh(utils.IdentCircle(20))
		circleShader = graphics.NewShaderProgram("assets/shader/circleVert.glsl", "assets/shader/circleFrag.glsl")
		objLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_ObjTrans\x00"))
		projLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_ProjTrans\x00"))
		colorLoc = gl.GetUniformLocation(uint32(circleShader), gl.Str("u_Color\x00"))
	}

	circle := new(Circle)
	*circle = Circle{
		id,
		utils.NewObject(&mesh, size, objLoc),
		owner,
		position,
		velocity,
		size,
		color,
	}

	return circle
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
	vLen := circle.velocity.Len()
	if vLen != 0 {
		if vLen < acceleration*delta {
			circle.velocity = mgl32.Vec2{0.0, 0.0}
		} else {
			circle.velocity = circle.velocity.Add(circle.velocity.Normalize().Mul(acceleration * delta))
		}
	}

	if circle.position.Len() > SIZE-circle.size {
		circle.position = circle.position.Normalize().Mul(SIZE - circle.size)
		a := math.Acos(float64(circle.position.Normalize().Dot(circle.velocity.Normalize())))
		rot := mgl32.Rotate2D(float32(a))
		circle.velocity = rot.Mul2x1(circle.position.Normalize().Mul(-vLen))
	}

	circle.object.Update(circle.position, circle.size)
}
