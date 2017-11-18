package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"math/rand"
	"time"
)

type CircleFactory struct {
	wait float32
}

const (
	minDelta = 0.5
	maxDelta = 5

	minSize = 0.01
	maxSize = 0.05
)

func NewFactory() CircleFactory {
	rand.Seed(int64(time.Now().Minute()))
	return CircleFactory{0}
}

func (factory *CircleFactory) Update(delta float32, path *[]Circle) {
	factory.wait -= delta
	if factory.wait < 0 {
		factory.wait += randBetween(minDelta, maxDelta)

		size := randBetween(minSize*SIZE, maxSize*SIZE)

		*path = append(*path, NewCircle(
			 false, size,
			mgl32.Vec2{randBetween(0, SIZE - size), randBetween(0, SIZE - size)},
			mgl32.Vec2{randBetween(0, SIZE * 4), randBetween(0, SIZE * 4)},
			mgl32.Vec3{0.0, 0.6, 0.0}))
	}
}

func randBetween(a float32, b float32) float32 {
	return a + rand.Float32()*(b-a)
}
