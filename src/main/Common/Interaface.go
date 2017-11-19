package Common

import (
	"encoding/gob"
)

const (
	CreateMessage     byte = 0
	DestroyMessage    byte = 1
	UpdateMessage     byte = 2
	ConnectMessage    byte = 3
	SynchonizeMessage byte = 4
)

func CreationMessage(ID int, OWNER int, SIZE float32,
	X float32, Y float32, VX float32, VY float32,
	R float32, G float32, B float32) GeneralMessage {
	return GeneralMessage{CreateMessage, ID, OWNER, SIZE, X, Y,
		VX, VY, R, G, B, make([]SynchonizeData, 0)}
}

func DestructionMessage(ID int) GeneralMessage {
	return GeneralMessage{DestroyMessage, ID, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, make([]SynchonizeData, 0)}
}

func UpdationMessage(ID int, SIZE float32) GeneralMessage {
	return GeneralMessage{UpdateMessage, ID, -1, SIZE, -1, -1,
		-1, -1, -1, -1, -1, make([]SynchonizeData, 0)}
}

func ConnectionMessage(OWNER int) GeneralMessage {
	return GeneralMessage{ConnectMessage, -1, OWNER, -1, -1, -1,
		-1, -1, -1, -1, -1, make([]SynchonizeData, 0)}
}

type SynchonizeData struct {
	ID     int
	X, Y   float32
	VX, VY float32
}

func SynchronizationMessage(data []SynchonizeData) GeneralMessage {
	return GeneralMessage{SynchonizeMessage, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, data}
}

type ActionMessage struct {
	ID   int
	X, Y float32
}

type GeneralMessage struct {
	TYPE    byte
	ID      int
	OWNER   int
	SIZE    float32
	X, Y    float32
	VX, VY  float32
	R, G, B float32
	DATA    []SynchonizeData
}

func RegisterInterface() {
	gob.Register(ActionMessage{})
	gob.Register(GeneralMessage{})
}
