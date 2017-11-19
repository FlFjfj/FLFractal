package game

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"main/graphics"
	"main/utils"
	"math"
)

var circleCounter int = -1

func nextId() int {
	circleCounter++
	return circleCounter
}

type World struct {
	cam utils.OrthographicCamera

	shader  graphics.ShaderProgram
	object  utils.GameObject
	projLoc int32

	factory CircleFactory
	circles map[int]*Circle
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
		make(map[int]*Circle),
	}

	player := NewCircle(nextId(), true, SIZE/4, mgl32.Vec2{0.0, 0.0}, mgl32.Vec2{0.0, 0.0}, mgl32.Vec3{1.0, 0.0, 0.0})
	result.circles[player.id] = player

	return result
}

func (world *World) Update(delta float32, window glfw.Window) {
	world.factory.Update(delta, &world.circles)

	for key := range world.circles {
		if world.circles[key] != nil {
			world.circles[key].Update(delta)
		}
	}

	world.processMouse(window)
	world.processCollision()
}

func (world *World) Draw(worldTrans mgl32.Mat4) {
	world.shader.Begin()
	gl.UniformMatrix4fv(world.projLoc, 1, false, &worldTrans[0])
	world.object.Draw()
	world.shader.End()

	for key := range world.circles {
		if world.circles[key] != nil {
			world.circles[key].Draw(worldTrans)
		}
	}
}

var (
	lastMouseState = false
	choosenId      = -1
)

func (world *World) processMouse(window glfw.Window) {
	state := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
	x64, y64 := window.GetCursorPos()
	position := mgl32.Vec2{float32((x64/WIDTH - 0.5) * 2 * SIZE * WIDTH / HEIGHT), float32(-(y64/HEIGHT - 0.5) * SIZE * 2)}

	if state && !lastMouseState {
		lastMouseState = true

		for _, circle := range world.circles {
			if circle != nil && circle.owns && circle.position.Sub(position).Len() < circle.size {
				choosenId = circle.id
			}
		}

	} else if !state && lastMouseState {
		lastMouseState = false
		if choosenId != -1 && world.circles[choosenId] != nil {
			newSize := world.circles[choosenId].size / 2
			direction := world.circles[choosenId].position.Sub(position).Normalize()
			power := world.circles[choosenId].position.Sub(position).Len()
			first := NewCircle(
				nextId(), true, newSize,
				world.circles[choosenId].position.Add(direction.Mul(newSize)),
				direction.Mul(power*2), world.circles[choosenId].color)
			second := NewCircle(
				nextId(), true, newSize,
				world.circles[choosenId].position.Add(direction.Mul(-newSize)),
				direction.Mul(-power*2), world.circles[choosenId].color)
			world.circles[first.id] = first
			world.circles[second.id] = second

			world.circles[choosenId] = nil
			choosenId = -1
		}
	}
}

func (world *World) processCollision() {
	for i := 0; i < circleCounter; i++ {
		for j := i + 1; j <= circleCounter; j++ {
			if world.circles[i] != nil && world.circles[j] != nil {
				first := world.circles[i]
				second := world.circles[j]
				dist := first.position.Sub(second.position).Len()
				ratio := math.Max(float64(first.size / second.size), float64(second.size / first.size))

				if ratio < 1.2 && dist < first.size + second.size {
					deltaSF := first.position.Sub(second.position).Normalize()
					projF := deltaSF.Mul(deltaSF.Dot(first.velocity))
					projS := deltaSF.Mul(deltaSF.Dot(second.velocity))
					projD := projF.Sub(projS)

					first.velocity = first.velocity.Sub(projD)
					second.velocity = second.velocity.Add(projD)
				}

				if first.size < second.size {
					first, second = second, first
				}

				if ratio >= 1.2 && dist < first.size {
					first.size += second.size
					world.circles[second.id] = nil
				}

			}
		}
	}
}
