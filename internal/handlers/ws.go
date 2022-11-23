package handlers

import (
	"net/http"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

type wsJsonResponse struct {
	Message    string `json:"message"`
	Action     string `json:"action"`
	ActionType string `json:"action_type"`
}

var UpgradeConnection = websocket.Upgrader{
	WriteBufferSize: 1024, // 1 MB
	ReadBufferSize:  1024, // 1 MB
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var wsChann = make(chan wsJsonPayload)
var clients = make(map[*WebSocketConnection]string)

type wsJsonPayload struct {
	Message string               `json:"message"`
	Action  string               `json:"action"`
	Conn    *WebSocketConnection `json:"-"`
}

type WebSocketConnection struct {
	*websocket.Conn
}

// Upgrade http connection with socket
func WsUpgrade(w http.ResponseWriter, r *http.Request) {

	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("---Error-> Upgrade--->", err.Error())
	}

	jsonRes := wsJsonResponse{}
	jsonRes.Message = `<div>Connection Successfully!</div>`

	// create socket connection
	conn := &WebSocketConnection{Conn: ws}

	// create client
	clients[conn] = ""

	// send msg to all client about created to connection
	ws.WriteJSON(jsonRes)

	go ListenForWs(conn)
}

// ListenForWs is a goroutine that handles communication between server and client, and
// feeds data into the wsChan
func ListenForWs(conn *WebSocketConnection) {

	var payload wsJsonPayload

	for {
		err := conn.Conn.ReadJSON(&payload)
		if err != nil {
			// not send
		} else {
			payload.Conn = conn
			wsChann <- payload
		}
	}
}

// ListenToWsChannel is a goroutine that waits for an entry on the wsChan, and handles it according to the
// specified action
func ListenToWs(conn *WebSocketConnection) {

	var response wsJsonResponse

	for {
		e := <-wsChann
		response.Message = "Some message with action " + e.Action

		BroadcastToAll(response)
	}
}

// broadcastToAll sends ws response to all connected clients
func BroadcastToAll(resp wsJsonResponse) {
	for client := range clients {
		err := client.Conn.WriteJSON(resp)
		if err != nil {

			logger.Error("---Error-> BroadcastToAll--->", err.Error())
			_ = client.Conn.Close()
			delete(clients, client)
		}
	}
}
