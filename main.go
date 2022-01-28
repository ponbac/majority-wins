package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ponbac/majority-wins/game"
)

func main() {
	port := os.Getenv("PORT")
	fmt.Println("Fetched from env: ", port)
	if port == "" {
		port = "8080"
	}
	room := game.NewRoom()
	room.Questions = []game.Question{{Type: "Music", Description: "1 or 2, little human?", Answers: make(map[*game.Player]int)},
		{Type: "Trivia", Description: "Is KalleK strong 1=YES, 2=NO?", Answers: make(map[*game.Player]int)}}
	go room.Run()
	go room.StartGame()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeWs(room, w, r)
	})

	log.Println("Starting server on http://localhost" + ":" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
