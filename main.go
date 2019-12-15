package main

import (
	"fmt"
	"net/http"

	"github.com/farhanramadhan/app-chat/manager"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var clientManager = manager.CreateClientManager()

func main() {
	port := 12345
	fmt.Printf("Starting application at port %d \n", port)
	go clientManager.Start()
	http.HandleFunc("/ws", WsPage)
	http.ListenAndServe(":12345", nil)
}

// WsPage is handler for web socket page
func WsPage(res http.ResponseWriter, req *http.Request) {
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

	clientManager.Register <- client

	go client.Read(&clientManager)
	go client.Write()
}
