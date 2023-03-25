package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lucsky/cuid"

	"qr-auth/constants"
)

var connections = make(map[string]*ConnectionStruct)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(request *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleConnection(writer http.ResponseWriter, request *http.Request) {
	connection, connectionError := upgrader.Upgrade(writer, request, nil)
	if connectionError != nil {
		return
	}
	defer connection.Close()

	// create an ID for this specific connection
	connectionId := cuid.New()

	// send generated ID for this specific connection to the client
	connection.WriteJSON(MessageStruct{
		Data:  connectionId,
		Event: constants.EVENTS.RegisterConnection,
	})

	// store connection
	newConnection := ConnectionStruct{
		Connection: connection,
		Name:       "",
	}
	connections[connectionId] = &newConnection
	log.Println("Connected", connectionId, "| Total connections:", len(connections))

	for {
		// parse incoming message, exit loop if there's a problem & delete connection
		var parsedMessage MessageStruct
		if parsingError := connection.ReadJSON(&parsedMessage); parsingError != nil {
			delete(connections, connectionId)
			break
		}

		// register user: add name
		if parsedMessage.Event == constants.EVENTS.RegisterUser {
			newConnection.Name = parsedMessage.Data
		}

		if parsedMessage.Event == configuration.EVENTS.TransferContacts &&
			parsedMessage.Issuer != "" && parsedMessage.Target != "" &&
			parsedMessage.Data != "" {
			var target *ConnectionStruct
			for i := range connections {
				if connections[i].ConnectionId == parsedMessage.Target {
					target = connections[i]
				}
			}
			if target != nil {
				target.Connection.WriteJSON(MessageStruct{
					Data:   parsedMessage.Data,
					Event:  configuration.EVENTS.TransferContacts,
					Issuer: connectionId,
					Target: target.ConnectionId,
				})
			}
		}

		if parsedMessage.Event == configuration.EVENTS.TransferComplete &&
			parsedMessage.Issuer != "" && parsedMessage.Target != "" {
			var target *ConnectionStruct
			for i := range connections {
				if connections[i].ConnectionId == parsedMessage.Target {
					target = connections[i]
				}
			}
			if target != nil {
				target.Connection.WriteJSON(MessageStruct{
					Event:  configuration.EVENTS.TransferComplete,
					Issuer: connectionId,
					Target: target.ConnectionId,
				})
			}
		}
	}
}
