package handler

import "github.com/gorilla/websocket"

type ConnectionStruct struct {
	Connection *websocket.Conn
	Name       string
}

type MessageStruct struct {
	Data  string `json:"data"`
	Event string `json:"event"`
}

type RegisterConnectionDataStruct struct {
	ConnectionId string `json:"connectionId"`
}
