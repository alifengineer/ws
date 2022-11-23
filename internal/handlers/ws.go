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

func WsUpgrade(w http.ResponseWriter, r *http.Request) {

	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("---Error-> Upgrade--->", err.Error())
	}

	jsonRes := wsJsonResponse{}
	jsonRes.Message = `<div>Connection Successfully!</div>`

	ws.WriteJSON(jsonRes)
}
