package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ponbac/majority-wins/data"
	"github.com/ponbac/majority-wins/game"
)

// Holds all rooms, key = room ID, value = room pointer
var rooms = map[string]*game.Room{}

func createRoom(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	var roomID string
	for ok := true; ok; _, ok = rooms[roomID] {
		roomID = randomString(4)
	}

	room := game.NewRoom(roomID)
	rooms[roomID] = room
	room.Questions = data.FetchQuestions()
	nQuestions := r.URL.Query().Get("questions")
	if nQuestions != "" {
		n, err := strconv.Atoi(nQuestions)
		if err != nil {
			log.Println(err)
		} else {
			room.NQuestions = n
		}
	}
	go room.Run()
	log.Println("Created room " + roomID)
	game.ServeWs(room, true, name, w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Hello Bobba!"))
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	deleteFinishedRooms() // Prevent user to join finished room

	name := r.URL.Query().Get("name")
	roomID := r.URL.Query().Get("room")
	if _, ok := rooms[roomID]; ok {
		room := rooms[roomID]
		err := game.ServeWs(room, false, name, w, r)
		if err != nil {
			w.Write([]byte("Name already taken"))
		}
	} else {
		w.Write([]byte("Room " + roomID + " does not exist!"))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No PORT in env, using default port 8080")
		port = "8080"
	} else {
		fmt.Println("Using PORT from env: ", port)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/new", createRoom)
	http.HandleFunc("/join", joinRoom)

	log.Println("Starting server on http://localhost" + ":" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// TODO: This does not work?
func deleteFinishedRooms() {
	for roomID, room := range rooms {
		if !room.Active {
			log.Println("Deleted room " + roomID)
			delete(rooms, roomID)
		} else {
			log.Println("Room " + roomID + " is active")
		}
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
