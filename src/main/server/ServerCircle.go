package server

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"main/Common"
)

const (
	acceleration float32 = -7
)

type ServerCircle struct {
	id       int
	owner    int
	Position mgl32.Vec2
	Velocity mgl32.Vec2
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
	circle.Position = circle.Position.Add(circle.Velocity.Mul(delta))
	vLen := circle.Velocity.Len()
	if vLen != 0 {
		if vLen < acceleration*delta {
			circle.Velocity = mgl32.Vec2{0.0, 0.0}
		} else {
			circle.Velocity = circle.Velocity.Add(circle.Velocity.Normalize().Mul(acceleration * delta))
		}
	}

	if circle.Position.Len() > Common.SIZE-circle.size {
		circle.Position = circle.Position.Normalize().Mul(Common.SIZE - circle.size)
		a := math.Acos(float64(circle.Position.Normalize().Dot(circle.Velocity.Normalize())))
		var minus float32
		if math.Signbit(float64(circle.Position.X() * circle.Velocity.Y() - circle.Velocity.X()+circle.Position.Y())) {
			minus = 1;
		} else {
			minus = -1;
		}

		rot := mgl32.Rotate2D(minus*float32(a))
		circle.Velocity = rot.Mul2x1(circle.Position.Normalize().Mul(-vLen))
	}
}
