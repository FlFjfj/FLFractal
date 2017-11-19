package server

import (
	"encoding/gob"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"net"
	"sync"
	"main/Common"
	"time"
)

type Client struct {
	conn    net.Conn
	encoder *gob.Encoder
}

type ServerWorld struct {
	factory CircleFactory
	circles map[int]*ServerCircle
	lastId  int

	actionQueue  chan Common.ActionMessage
	messageQueue chan Common.GeneralMessage
	mutex        *sync.Mutex
	clients      []Client
}

func NewServerWorld() *ServerWorld {
	result := new(ServerWorld)
	*result = ServerWorld{

		NewFactory(),
		make(map[int]*ServerCircle),
		-1,
		make(chan Common.ActionMessage, 10000),
		make(chan Common.GeneralMessage, 10000),
		&sync.Mutex{},
		[]Client{},
	}
	Common.RegisterInterface()
	listener, _ := net.Listen("tcp", ":8090")
	go func() {
		for {
			conn, _ := listener.Accept()
			fmt.Println("Connected")
			result.mutex.Lock()
			//println("lock accept client")
			result.clients = append(result.clients, Client{conn, gob.NewEncoder(conn)})

			owner := len(result.clients) - 1
			circle := NewCircle(result.nextId(), owner, Common.DEFAULT_SIZE,
				mgl32.Vec2{randBetween(0, Common.SIZE-Common.DEFAULT_SIZE) - (Common.SIZE-Common.DEFAULT_SIZE) / 2,
					randBetween(0, Common.SIZE-Common.DEFAULT_SIZE) - (Common.SIZE-Common.DEFAULT_SIZE) / 2},
				mgl32.Vec2{0.0, 0.0},
				mgl32.Vec3{randBetween(0, 1), randBetween(0, 1), randBetween(0, 1)})
			result.circles[circle.id] = circle

			result.clients[owner].encoder.Encode(Common.ConnectionMessage(circle.owner))
			result.messageQueue <- Common.CreationMessage(circle.id, circle.owner, circle.size,
				circle.Position.X(), circle.Position.Y(), circle.Velocity.X(), circle.Velocity.Y(),
				circle.color.X(), circle.color.Y(), circle.color.Z())

			result.mutex.Unlock()
			//println("unlock accept client")
			go func() {
				result.mutex.Lock()
				//println("lock init client")
				for i := 0; i <= result.lastId; i++ {
					if result.circles[i] != nil {
						result.clients[owner].encoder.Encode(Common.CreationMessage(i, result.circles[i].owner, result.circles[i].size,
							result.circles[i].Position.X(), result.circles[i].Position.Y(), result.circles[i].Velocity.X(), result.circles[i].Velocity.Y(),
							result.circles[i].color.X(), result.circles[i].color.Y(), result.circles[i].color.Z()))
					}
				}

				result.mutex.Unlock()
				//println("unlock init client")
				decoder := gob.NewDecoder(conn)
				for {
					action := Common.ActionMessage{}
					decoder.Decode(&action)
					result.actionQueue <- action
				}
			}()
		}
	}()

	go func() {
		message := Common.GeneralMessage{}
		for {
			message = <-result.messageQueue
			//fmt.Printf("Message: %+v\n", message)
			result.mutex.Lock()
			//println("lock broadcast")
			for _, client := range result.clients {
				client.encoder.Encode(message)
			}
			result.mutex.Unlock()
			//println("unlock broadcast")
		}
	}()


	go func() {
		for {
			action := <-result.actionQueue
			//fmt.Printf("Action: %+v\n", action)
			result.mutex.Lock()
			//println("lock Action process")
			position := mgl32.Vec2{action.X, action.Y}
			newSize := result.circles[action.ID].size / 2
			direction := result.circles[action.ID].Position.Sub(position).Normalize()
			power := float32(math.Max(float64(result.circles[action.ID].Position.Sub(position).Len()), 0.5))

			first := NewCircle(
				result.nextId(), result.circles[action.ID].owner, newSize,
				result.circles[action.ID].Position.Add(direction.Mul(newSize)),
				direction.Mul(power*2), result.circles[action.ID].color)
			second := NewCircle(
				result.nextId(), result.circles[action.ID].owner, newSize,
				result.circles[action.ID].Position.Add(direction.Mul(-newSize)),
				direction.Mul(-power*2), result.circles[action.ID].color)

			result.circles[action.ID] = nil
			result.circles[first.id] = first
			result.circles[second.id] = second

			result.messageQueue <- Common.DestructionMessage(action.ID)
			result.messageQueue <- Common.CreationMessage(first.id, first.owner, first.size, first.Position.X(), first.Position.Y(),
				first.Velocity.X(), first.Velocity.Y(), first.color.X(), first.color.Y(), first.color.Z())
			result.messageQueue <- Common.CreationMessage(second.id, second.owner, second.size, second.Position.X(), second.Position.Y(),
				second.Velocity.X(), second.Velocity.Y(), second.color.X(), second.color.Y(), second.color.Z())
			result.mutex.Unlock()
			//println("unlock Action process")
		}
	}()

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			result.mutex.Lock()
			//println("lock sync process")

			data := make([]Common.SynchonizeData, 0)
			for i:= 0; i <= result.lastId; i++ {
				if result.circles[i] != nil {
					data = append(data, Common.SynchonizeData{result.circles[i].id,
						result.circles[i].Position.X(), result.circles[i].Position.Y(),
						result.circles[i].Velocity.X(), result.circles[i].Velocity.Y()})
				}
			}
			result.messageQueue <- Common.SynchronizationMessage(data)

			result.mutex.Unlock()
			//println("unlock sync process")
		}
	}()

	return result
}

func (world *ServerWorld) Update(delta float32) {
	world.mutex.Lock()
	world.factory.Update(delta, world)
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
				dist := first.Position.Sub(second.Position).Len()
				ratio := math.Max(float64(first.size/second.size), float64(second.size/first.size))

			//	fmt.Printf("First: %+v, Second: %+v, Dist: %f, ration: %f\n", first, second, dist, ratio)

				if ratio < 1.2 && dist < first.size+second.size {
					//println("Crashed")
					deltaSF := first.Position.Sub(second.Position).Normalize()
					projF := deltaSF.Mul(deltaSF.Dot(first.Velocity))
					projS := deltaSF.Mul(deltaSF.Dot(second.Velocity))
					projD := projF.Sub(projS)

					first.Velocity = first.Velocity.Sub(projD)
					second.Velocity = second.Velocity.Add(projD)
				}

				if first.size < second.size {
					first, second = second, first
				}

				if ratio >= 1.2 && dist < first.size {
				//	println("Eaten")
					newSize := float32(math.Sqrt(float64((first.size * first.size) + (second.size * second.size))))
					world.messageQueue <- Common.DestructionMessage(second.id)
					world.messageQueue <- Common.UpdationMessage(first.id, newSize)

					first.size = newSize
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
