package handler

import (
	"net/http"

	"github.com/farhanramadhan/app-chat/manager"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// WsPage is handler for web socket page
func WsPage(res http.ResponseWriter, req *http.Request,
	cm *manager.ClientManager) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &manager.Client{
		ID:     uuid.NewV4().String(),
		Socket: conn,
		Send:   make(chan []byte),
	}

	cm.Register <- client

	go client.Read(cm)
	go client.Write()
}
