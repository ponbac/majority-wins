package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"math/rand"

	"github.com/ponbac/majority-wins/game"
)

var rooms = map[string]*game.Room{}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var roomID string
	for ok := true; ok; _, ok = rooms[roomID] {
		roomID = randomString(4)
	}

	room := game.NewRoom(roomID)
	rooms[roomID] = room
	room.Questions = []*game.Question{{Type: "Music", Description: "1 or 2, little human?", Answers: make(map[*game.Player]int), Reward: 2},
		{Type: "Trivia", Description: "Is KalleK strong 1=YES, 2=NO?", Answers: make(map[*game.Player]int), Reward: 2}}
	go room.Run()
	go room.StartGame()
	log.Println("Created room " + roomID)
	game.ServeWs(room, w, r)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if _, ok := rooms[roomID]; ok {
		room := rooms[roomID]
		game.ServeWs(room, w, r)
	} else {
		w.Write([]byte("Room " + roomID + " does not exist!"))
	}
}

func main() {
	port := os.Getenv("PORT")
	fmt.Println("Fetched from env: ", port)
	if port == "" {
		port = "8080"
	}
	// room := game.NewRoom("1337")
	// rooms["1337"] = room
	// room.Questions = []*game.Question{{Type: "Music", Description: "1 or 2, little human?", Answers: make(map[*game.Player]int), Reward: 2},
	// 	{Type: "Trivia", Description: "Is KalleK strong 1=YES, 2=NO?", Answers: make(map[*game.Player]int), Reward: 2}}
	// go room.Run()
	// go room.StartGame()
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	game.ServeWs(room, w, r)
	// })
	http.HandleFunc("/new", createRoom)
	http.HandleFunc("/join", joinRoom)

	log.Println("Starting server on http://localhost" + ":" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func randomString(n int) string {
    var letters = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ123456789")
 
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}