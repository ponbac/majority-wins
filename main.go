package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/ponbac/majority-wins/game"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var testRoom = game.Room{
	ID:      "1337",
	Players: []game.Player{},
	Questions: []game.Question{
		{
			Type: "Music",
		},
		{
			Type: "Trivia",
		},
	},
	CurrentQuestion: 0,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received: %s of type %d", msg, messageType)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	ws.SetCloseHandler(func(code int, text string) error {
		playerRemoved := testRoom.RemovePlayer(ws)
		if !playerRemoved {
			log.Printf("Could not remove player from room!")
		}
		log.Println("Client disconnected!")
		for _, p := range testRoom.Players {
			if err := p.Conn.WriteJSON(testRoom); err != nil {
				log.Println(err)
			}
		}
		return nil
	})

	log.Println("Client connected!")
	testRoom.AddPlayer(game.Player{
		Name:  "Player " + fmt.Sprint(len(testRoom.Players)+1),
		Score: 0,
		Conn:  ws,
	})
	for _, p := range testRoom.Players {
		if err := p.Conn.WriteJSON(testRoom); err != nil {
			log.Println(err)
			return
		}
	}
	go reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	flag.Parse()

	setupRoutes()

	log.Println("Starting server on http://localhost" + *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
