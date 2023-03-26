package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lucsky/cuid"

	"qr-auth/constants"
)

// TODO: locking
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
		Connection:          connection,
		LastMessageReceived: time.Now().UnixMilli(),
		Name:                "",
	}
	connections[connectionId] = &newConnection
	log.Println("Connected", connectionId, "| Total connections:", len(connections))

	for {
		// parse incoming message, exit loop if there's a problem & delete connection
		var parsedMessage MessageStruct
		if parsingError := connection.ReadJSON(&parsedMessage); parsingError != nil {
			connection.WriteJSON(MessageStruct{
				Event: constants.EVENTS.ServerDisconnect,
			})
			connection.Close()
			delete(connections, connectionId)
			log.Println(
				"Disconnected",
				connectionId,
				" [invalid client message] | Total connections:",
				len(connections),
			)
			break
		}

		// authenticate target
		if parsedMessage.Event == constants.EVENTS.AuthenticateTarget {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			target := connections[parsedMessage.Data]
			if target == nil {
				connection.WriteJSON(MessageStruct{
					Event: constants.EVENTS.InvalidTarget,
				})
				continue
			}
			if newConnection.Name == "" {
				connection.WriteJSON(MessageStruct{
					Event: constants.EVENTS.Unauthorized,
				})
				continue
			}
			target.Connection.WriteJSON(MessageStruct{
				Data:  newConnection.Name,
				Event: constants.EVENTS.AuthenticateTarget,
			})
			continue
		}

		// client disconnects
		if parsedMessage.Event == constants.EVENTS.ClientDisconnect {
			connection.Close()
			delete(connections, connectionId)
			log.Println(
				"Disconnected",
				connectionId,
				" [client request] | Total connections:",
				len(connections),
			)
			break
		}

		// ping response
		if parsedMessage.Event == constants.EVENTS.PingResponse {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			continue
		}

		// register user: add name
		if parsedMessage.Event == constants.EVENTS.RegisterUser {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			newConnection.Name = parsedMessage.Data
			continue
		}
	}
}

func PingService() {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			for connectionId := range connections {
				connection := connections[connectionId]
				frame := time.Now().UnixMilli() - constants.CONNECTION_TIMEOUT
				if connection.LastMessageReceived < frame {
					connection.Connection.WriteJSON(MessageStruct{
						Event: constants.EVENTS.ServerDisconnect,
					})
					connection.Connection.Close()
					delete(connections, connectionId)
					log.Println(
						"Disconnected",
						connectionId,
						" [client is non-responsive] | Total connections:",
						len(connections),
					)
					continue
				}
				if time.Now().UnixMilli()-connection.LastMessageReceived > 30*1000 {
					connection.Connection.WriteJSON(MessageStruct{
						Event: constants.EVENTS.Ping,
					})
				}
			}
		}
	}()
}
