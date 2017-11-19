package main

import (
	"main/server"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	serverProgram = server.NewServerProgram(updateServer)
	serverWorld = server.NewServerWorld()

	for {
		serverProgram.Update()
	}
}

var (
	serverProgram server.Program
	serverWorld   *server.ServerWorld
)

func updateServer(delta float32) {
	serverWorld.Update(delta)
}
