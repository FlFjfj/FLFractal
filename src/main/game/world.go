package game

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"main/game/net"
	"main/graphics"
	"main/utils"
	"math"
	//"os"
	"fmt"
	"os"
	"sync"
	"main/Common"
)

type World struct {
	cam utils.OrthographicCamera

	shader  graphics.ShaderProgram
	object  utils.GameObject
	projLoc int32

	owner_id int

	circles map[int]*Circle
	lastId  int

	actionQueue  chan Common.ActionMessage
	messageQueue chan Common.GeneralMessage
	mutex        sync.Mutex
}

func NewWorld(camera utils.OrthographicCamera) World {
	worldShader := graphics.NewShaderProgram("assets/shader/worldVert.glsl", "assets/shader/worldFrag.glsl")
	transformLoc := gl.GetUniformLocation(uint32(worldShader), gl.Str("u_ObjTrans\x00"))
	projectionLoc := gl.GetUniformLocation(uint32(worldShader), gl.Str("u_ProjTrans\x00"))
	circle := utils.NewMesh(utils.IdentCircle(100))
	object := utils.NewObject(&circle, Common.SIZE, transformLoc)

	result := World{
		camera,
		worldShader,
		object,
		projectionLoc,
		0,
		make(map[int]*Circle),
		-1,
		make(chan Common.ActionMessage, 10000),
		make(chan Common.GeneralMessage, 10000),
		sync.Mutex{},
	}

	net.Connect(os.Args[1], result.actionQueue, result.messageQueue)
	return result
}

func (world *World) Update(delta float32, window glfw.Window) {
	for key := range world.circles {
		if world.circles[key] != nil {
			world.circles[key].Update(delta)
		}
	}

	world.processMouse(window, world.actionQueue)
	world.processCollision()
	world.processMessages()
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

func (world *World) processMouse(window glfw.Window, actionQueue chan Common.ActionMessage) {
	state := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
	x64, y64 := window.GetCursorPos()
	position := mgl32.Vec2{float32((x64/Common.WIDTH - 0.5) * 2 * Common.SIZE * Common.WIDTH / Common.HEIGHT),
	float32(-(y64/Common.HEIGHT - 0.5) * Common.SIZE * 2)}

	if state && !lastMouseState {
		lastMouseState = true

		for _, circle := range world.circles {
			if circle != nil && circle.owner == world.owner_id && circle.position.Sub(position).Len() < circle.size {
				choosenId = circle.id
				break
			}
		}

	} else if !state && lastMouseState {
		lastMouseState = false
		if choosenId != -1 && world.circles[choosenId] != nil {
			actionQueue <- Common.ActionMessage{choosenId, position.X(), position.Y()}
			choosenId = -1
		}
	}
}

func (world *World) processCollision() {
	for i := 0; i < world.lastId; i++ {
		for j := i + 1; j <= world.lastId; j++ {
			if world.circles[i] != nil && world.circles[j] != nil {
				first := world.circles[i]
				second := world.circles[j]
				dist := first.position.Sub(second.position).Len()
				ratio := math.Max(float64(first.size/second.size), float64(second.size/first.size))

				if ratio < 1.2 && dist < first.size+second.size {
					deltaSF := first.position.Sub(second.position).Normalize()
					projF := deltaSF.Mul(deltaSF.Dot(first.velocity))
					projS := deltaSF.Mul(deltaSF.Dot(second.velocity))
					projD := projF.Sub(projS)

					first.velocity = first.velocity.Sub(projD)
					second.velocity = second.velocity.Add(projD)
				}
			}
		}
	}
}

func (world *World) processMessages() {
	select {
	case message, ok := <-world.messageQueue:
		{
			if ok {
				fmt.Printf("Message: %+v\n", message)
				switch message.TYPE {
				case Common.ConnectMessage:
					{
						world.owner_id = message.OWNER
					}
				case Common.CreateMessage:
					{
						if message.ID > world.lastId {
							world.lastId = message.ID
						}

						world.circles[message.ID] = NewCircle(
							message.ID, message.OWNER, message.SIZE,
							mgl32.Vec2{message.X, message.Y},
							mgl32.Vec2{message.VX, message.VY},
							mgl32.Vec3{message.R, message.G, message.B})
					}
				case Common.DestroyMessage:
					{
						world.circles[message.ID] = nil
					}
				case Common.UpdateMessage:
					{
						if world.circles[message.ID] != nil {
							world.circles[message.ID].size = message.SIZE
						}
					}
				case Common.SynchonizeMessage:
					{
						for _, data := range message.DATA {
							if world.circles[data.ID] != nil {
								world.circles[data.ID].position = mgl32.Vec2{data.X, data.Y}
								world.circles[data.ID].velocity = mgl32.Vec2{data.VX, data.VY}
							}
						}
					}
				}
			}
		}
	default: {return}
	}
}
