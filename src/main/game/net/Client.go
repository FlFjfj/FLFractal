package net

import (
	"encoding/gob"
	"net"
	"fmt"
	"main/Common"
)

func Connect(address string, actionQueue chan Common.ActionMessage, messageQueue chan Common.GeneralMessage) {
	Common.RegisterInterface()

	serverAddress, _ := net.ResolveTCPAddr("tcp", address)
	connection, _ := net.DialTCP("tcp", nil, serverAddress)

	go processRead(connection, messageQueue)
	go processWrite(connection, actionQueue)
}

func processRead(conn net.Conn, messageQueue chan Common.GeneralMessage) {
	decoder := gob.NewDecoder(conn)
	for {
		var action = Common.GeneralMessage{}
		decoder.Decode(&action)
		messageQueue <- action
	}
}

func processWrite(conn net.Conn, actionQueue chan Common.ActionMessage) {
	encoder := gob.NewEncoder(conn)
	for {
		action :=<- actionQueue
		fmt.Printf("Action: %+v\n", action)
		encoder.Encode(action)
	}
}
