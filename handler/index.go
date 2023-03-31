package handler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lucsky/cuid"

	"qr-auth/constants"
)

var store struct {
	sync.Mutex
	Connections map[string]*ConnectionStruct
}

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

	connectionId := cuid.New()
	connection.WriteJSON(MessageStruct{
		Data:  connectionId,
		Event: constants.EVENTS.RegisterConnection,
	})

	newConnection := ConnectionStruct{
		Connection:          connection,
		LastMessageReceived: time.Now().UnixMilli(),
		Name:                "",
	}
	if store.Connections == nil {
		store.Connections = make(map[string]*ConnectionStruct)
	}
	store.Connections[connectionId] = &newConnection
	log.Println(
		"Connected",
		connectionId,
		"| Total connections:",
		len(store.Connections),
	)

	for {
		var parsedMessage MessageStruct
		if parsingError := connection.ReadJSON(&parsedMessage); parsingError != nil {
			connection.WriteJSON(MessageStruct{
				Event: constants.EVENTS.ServerDisconnect,
			})
			connection.Close()
			store.Lock()
			delete(store.Connections, connectionId)
			store.Unlock()
			log.Println(
				"Disconnected",
				connectionId,
				"[invalid client message] | Total connections:",
				len(store.Connections),
			)
			break
		}

		if parsedMessage.Event == constants.EVENTS.AuthenticateTarget {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			target := store.Connections[parsedMessage.Data]
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

		if parsedMessage.Event == constants.EVENTS.ClientDisconnect {
			connection.Close()
			store.Lock()
			delete(store.Connections, connectionId)
			store.Unlock()
			log.Println(
				"Disconnected",
				connectionId,
				"[client request] | Total connections:",
				len(store.Connections),
			)
			break
		}

		if parsedMessage.Event == constants.EVENTS.PingResponse {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			continue
		}

		if parsedMessage.Event == constants.EVENTS.RegisterUser {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			newConnection.Name = parsedMessage.Data
			continue
		}

		if parsedMessage.Event == constants.EVENTS.SignOut {
			newConnection.LastMessageReceived = time.Now().UnixMilli()
			newConnection.Name = ""
			continue
		}
	}
}

func PingService() {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			for connectionId := range store.Connections {
				connection := store.Connections[connectionId]
				frame := time.Now().UnixMilli() - constants.CONNECTION_TIMEOUT
				if connection.LastMessageReceived < frame {
					connection.Connection.WriteJSON(MessageStruct{
						Event: constants.EVENTS.ServerDisconnect,
					})
					connection.Connection.Close()
					store.Lock()
					delete(store.Connections, connectionId)
					store.Unlock()
					log.Println(
						"Disconnected",
						connectionId,
						"[client is non-responsive] | Total connections:",
						len(store.Connections),
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
