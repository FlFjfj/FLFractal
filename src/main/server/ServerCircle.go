package server

import (
	"github.com/go-gl/mathgl/mgl32"
	"main/game"
	"math"
)

const (
	acceleration float32 = -7
)

type ServerCircle struct {
	id       int
	owner    int
	position mgl32.Vec2
	velocity mgl32.Vec2
	size     float32
	color    mgl32.Vec3
}

func NewCircle(id int, owner int, size float32, position mgl32.Vec2, velocity mgl32.Vec2, color mgl32.Vec3) *ServerCircle {

	circle := new(ServerCircle)
	*circle = ServerCircle{
		id,
		owner,
		position,
		velocity,
		size,
		color,
	}
	return circle
}

func (circle *ServerCircle) Update(delta float32) {
	circle.position = circle.position.Add(circle.velocity.Mul(delta))
	vLen := circle.velocity.Len()
	if vLen != 0 {
		if vLen < acceleration*delta {
			circle.velocity = mgl32.Vec2{0.0, 0.0}
		} else {
			circle.velocity = circle.velocity.Add(circle.velocity.Normalize().Mul(acceleration * delta))
		}
	}

	if circle.position.Len() > game.SIZE-circle.size {
		circle.position = circle.position.Normalize().Mul(game.SIZE - circle.size)
		a := math.Acos(float64(circle.position.Normalize().Dot(circle.velocity.Normalize())))
		rot := mgl32.Rotate2D(float32(a))
		circle.velocity = rot.Mul2x1(circle.position.Normalize().Mul(-vLen))
	}
}
