package net

import (
	"encoding/gob"
	"net"
	"fmt"
)

func Connect(address string, actionQueue chan ActionMessage, messageQueue chan GeneralMessage) {
	RegisterInterface()

	serverAddress, _ := net.ResolveTCPAddr("tcp", address)
	connection, _ := net.DialTCP("tcp", nil, serverAddress)

	go processRead(connection, messageQueue)
	go processWrite(connection, actionQueue)
}

func processRead(conn net.Conn, messageQueue chan GeneralMessage) {
	decoder := gob.NewDecoder(conn)
	for {
		var action = GeneralMessage{}
		decoder.Decode(&action)
		messageQueue <- action
	}
}

func processWrite(conn net.Conn, actionQueue chan ActionMessage) {
	encoder := gob.NewEncoder(conn)
	for {
		action :=<- actionQueue
		fmt.Printf("Action: %+v\n", action)
		encoder.Encode(action)
	}
}
