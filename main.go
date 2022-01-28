package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ponbac/majority-wins/game"
)

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	flag.Parse()
	room := game.NewRoom()
	room.Questions = []game.Question{{Type: "Music", Description: "1 or 2, little human?", Answers: make(map[*game.Player]int)}, 
		{Type: "Trivia", Description: "Is KalleK strong 1=YES, 2=NO?", Answers: make(map[*game.Player]int)}}
	go room.Run()
	go room.StartGame()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeWs(room, w, r)
	})

	log.Println("Starting server on http://localhost" + *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
