package server

import (
	//"github.com/go-gl/mathgl/mgl32"
	"math/rand"
	"time"
	"github.com/go-gl/mathgl/mgl32"
	"main/Common"
)

type CircleFactory struct {
	wait float32
}

const (
	minDelta = 0.5
	maxDelta = 5

	minSize = 0.05
	maxSize = 0.025
)

func NewFactory() CircleFactory {
	rand.Seed(int64(time.Now().Minute()))
	return CircleFactory{0}
}

func (factory *CircleFactory) Update(delta float32, world *ServerWorld) {
		factory.wait -= delta
		if factory.wait < 0 {
			factory.wait += randBetween(minDelta, maxDelta)

			size := randBetween(minSize*Common.SIZE, maxSize*Common.SIZE)

			circle := NewCircle(
				world.nextId(), -1, size,
				mgl32.Vec2{randBetween(0, Common.SIZE-size), randBetween(0, Common.SIZE-size)},
				mgl32.Vec2{0.0, 0.0},
				mgl32.Vec3{0.0, 0.6, 0.0})
			world.circles[circle.id] = circle
			world.messageQueue <- Common.CreationMessage(circle.id, circle.owner, circle.size, circle.Position.X(), circle.Position.Y(),
				circle.Velocity.X(), circle.Velocity.Y(), circle.color.X(), circle.color.Y(), circle.color.Z())
		}
}

func randBetween(a float32, b float32) float32 {
	return a + rand.Float32()*(b-a)
}
