package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// define upgrader - requires Read and Wrie buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//Checking the origin of connection to make requests from React server to here
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader to listen for new messages being sent to WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// define WebSocket endpoint

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade connection to WebSocket conn
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
	}
	// listen for new messages
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	// map '/ws' endpoint to the 'serveWs' func
	http.HandleFunc("/ws", serveWs)
}

func main() {
	fmt.Println("Chat App")
	setupRoutes()
	http.ListenAndServe(":8081", nil)
}
