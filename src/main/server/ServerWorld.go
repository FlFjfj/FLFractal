package server

import (
	"encoding/gob"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"main/game"
	n "main/game/net"
	"math"
	"net"
	"sync"
)

type Client struct {
	conn    net.Conn
	encoder *gob.Encoder
}

type ServerWorld struct {
	factory CircleFactory
	circles map[int]*ServerCircle
	lastId  int

	actionQueue  chan n.ActionMessage
	messageQueue chan n.GeneralMessage
	mutex        *sync.Mutex
	clients      []Client
}

func NewServerWorld() ServerWorld {
	result := ServerWorld{

		NewFactory(),
		make(map[int]*ServerCircle),
		-1,
		make(chan n.ActionMessage, 10000),
		make(chan n.GeneralMessage, 10000),
		&sync.Mutex{},
		[]Client{},
	}
	listener, _ := net.Listen("tcp", ":8090")
	go func() {
		for {
			conn, _ := listener.Accept()
			fmt.Println("Connected")
			result.mutex.Lock()
			println("lock accept client")
			result.clients = append(result.clients, Client{conn, gob.NewEncoder(conn)})

			owner := len(result.clients) - 1
			circle := NewCircle(result.nextId(), owner, game.DEFAULT_SIZE,
				mgl32.Vec2{randBetween(0, game.SIZE-game.DEFAULT_SIZE) - (game.SIZE-game.DEFAULT_SIZE) / 2, randBetween(0, game.SIZE-game.DEFAULT_SIZE) - (game.SIZE-game.DEFAULT_SIZE) / 2},
				mgl32.Vec2{0.0, 0.0},
				mgl32.Vec3{randBetween(0, 1), randBetween(0, 1), randBetween(0, 1)})
			result.circles[circle.id] = circle

			result.clients[owner].encoder.Encode(n.ConnectionMessage(circle.owner))
			result.messageQueue <- n.CreationMessage(circle.id, circle.owner, circle.size,
				circle.position.X(), circle.position.Y(), circle.velocity.X(), circle.velocity.Y(),
				circle.color.X(), circle.color.Y(), circle.color.Z())

			result.mutex.Unlock()
			println("unlock accept client")
			go func() {
				result.mutex.Lock()
				println("lock init client")
				for i := 0; i <= result.lastId; i++ {
					if result.circles[i] != nil {
						result.clients[owner].encoder.Encode(n.CreationMessage(i, result.circles[i].owner, result.circles[i].size,
							result.circles[i].position.X(), result.circles[i].position.Y(), result.circles[i].velocity.X(), result.circles[i].velocity.Y(),
							result.circles[i].color.X(), result.circles[i].color.Y(), result.circles[i].color.Z()))
					}
				}

				result.mutex.Unlock()
				println("unlock init client")
				decoder := gob.NewDecoder(conn)
				for {
					action := n.ActionMessage{}
					decoder.Decode(&action)
					result.actionQueue <- action
				}
			}()
		}
	}()

	go func() {
		message := n.GeneralMessage{}
		for {
			message = <-result.messageQueue
			fmt.Printf("Message: %+v\n", message)
			result.mutex.Lock()
			println("lock broadcast")
			for _, client := range result.clients {
				client.encoder.Encode(message)
			}
			result.mutex.Unlock()
			println("unlock broadcast")
		}
	}()


	go func() {
		for {
			action := <-result.actionQueue
			fmt.Printf("Action: %+v\n", action)
			result.mutex.Lock()
			println("lock Action process")
			position := mgl32.Vec2{action.X, action.Y}
			newSize := result.circles[action.ID].size / 2
			direction := result.circles[action.ID].position.Sub(position).Normalize()
			power := result.circles[action.ID].position.Sub(position).Len()
			first := NewCircle(
				result.nextId(), result.circles[action.ID].owner, newSize,
				result.circles[action.ID].position.Add(direction.Mul(newSize)),
				direction.Mul(power*2), result.circles[action.ID].color)
			second := NewCircle(
				result.nextId(), result.circles[action.ID].owner, newSize,
				result.circles[action.ID].position.Add(direction.Mul(-newSize)),
				direction.Mul(-power*2), result.circles[action.ID].color)

			result.circles[first.id] = first
			result.circles[second.id] = second

			result.messageQueue <- n.DestructionMessage(action.ID)
			result.messageQueue <- n.CreationMessage(first.id, first.owner, first.size, first.position.X(), first.position.Y(),
				first.velocity.X(), first.velocity.Y(), first.color.X(), first.color.Y(), first.color.Z())
			result.messageQueue <- n.CreationMessage(second.id, second.owner, second.size, second.position.X(), second.position.Y(),
				second.velocity.X(), second.velocity.Y(), second.color.X(), second.color.Y(), second.color.Z())
			result.mutex.Unlock()
			println("unlock Action process")
		}
	}()

	return result
}

func (world *ServerWorld) Update(delta float32) {
	world.mutex.Lock()
	//println("lock func (world *ServerWorld) Update(delta float32) {")
	for key := range world.circles {
		if world.circles[key] != nil {
			world.circles[key].Update(delta)
		}
	}

	world.processCollision()
	world.mutex.Unlock()
	//println("unlock func (world *ServerWorld) Update(delta float32) {")
}

func (world *ServerWorld) processCollision() {
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

				if first.size < second.size {
					first, second = second, first
				}

				if ratio >= 1.2 && dist < first.size {

					world.messageQueue <- n.DestructionMessage(second.id)
					world.messageQueue <- n.UpdationMessage(first.id, first.size+second.size)

					first.size += second.size
					world.circles[second.id] = nil
				}
			}
		}
	}
}

func (world *ServerWorld) nextId() int {
	world.lastId++
	return world.lastId
}
