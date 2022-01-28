package game

import (
	"encoding/json"
	"fmt"
	"time"
)

type Room struct {
	ID              string
	Players         map[*Player]bool
	Questions       []Question
	CurrentQuestion int
	// 0 = not started, 1 = question time, 2 = question results, 3 = game over
	Scene int

	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Player
	// Unregister requests from clients.
	unregister chan *Player
}

type JSONRoom struct {
	ID              string          `json:"id"`
	Players         []*JSONPlayer   `json:"players"`
	Questions       []*JSONQuestion `json:"questions"`
	CurrentQuestion int             `json:"current_question"`
	Scene           int             `json:"scene"`
}

func NewRoom() *Room {
	return &Room{
		broadcast:       make(chan []byte),
		register:        make(chan *Player),
		unregister:      make(chan *Player),
		Players:         make(map[*Player]bool),
		ID:              "1337",
		Questions:       []Question{},
		CurrentQuestion: 0,
		Scene:           0,
	}
}

func (r *Room) ToJSON() []byte {
	jsonRoom := &JSONRoom{ID: r.ID, Players: []*JSONPlayer{}, Questions: []*JSONQuestion{}, CurrentQuestion: r.CurrentQuestion, Scene: r.Scene}

	for player := range r.Players {
		jsonRoom.Players = append(jsonRoom.Players, player.ToJSONPlayer())
	}
	for _, question := range r.Questions {
		jsonRoom.Questions = append(jsonRoom.Questions, question.ToJSONQuestion())
	}

	b, err := json.Marshal(jsonRoom)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func (r *Room) AddPlayer(player *Player) {
	r.Players[player] = true
	fmt.Println("Added " + player.Name + " to room " + r.ID)
	r.BroadcastRoomState()
}

func (r *Room) RemovePlayer(player *Player) {
	if _, ok := r.Players[player]; ok {
		delete(r.Players, player)
		close(player.send)
		fmt.Println("Removed " + player.Name + " from room " + r.ID)
		r.BroadcastRoomState()
	}
}

func (r *Room) NextQuestion() *Question {
	if r.CurrentQuestion >= len(r.Questions)-1 {
		return nil
	}
	r.CurrentQuestion++

	return &r.Questions[r.CurrentQuestion]
}

func (r *Room) BroadcastRoomState() {
	for player := range r.Players {
		select {
		case player.send <- r.ToJSON():
		default:
			close(player.send)
			delete(r.Players, player)
		}
	}
}

func (r *Room) StartGame() {
	r.Scene = 1
	r.BroadcastRoomState()

	prevScene := 0
	for {
		if r.Scene != prevScene {
			switch r.Scene {
			// Question time
			case 1:
				prevScene = 1
				fmt.Println("Starting question " + r.Questions[r.CurrentQuestion].Description)
			// Question results
			case 2:
				prevScene = 2
				fmt.Println("Question results")
				time.Sleep(time.Second * 5)
				if r.NextQuestion() == nil {
					r.Scene = 3
				} else {
					r.Scene = 1
				}
			// Game over
			case 3:
				prevScene = 3
				fmt.Println("Game over")
			}
			r.BroadcastRoomState()
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *Room) Run() {
	for {
		select {
		case player := <-r.register:
			r.AddPlayer(player)
		case player := <-r.unregister:
			r.RemovePlayer(player)
		case message := <-r.broadcast:
			for player := range r.Players {
				select {
				case player.send <- message:
				default:
					close(player.send)
					delete(r.Players, player)
				}
			}
		}
	}
}