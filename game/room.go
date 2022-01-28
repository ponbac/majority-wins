package game

import (
	"encoding/json"
	"fmt"
)

type Room struct {
	ID              string           `json:"id"`
	Players         map[*Player]bool 
	Questions       []Question       `json:"questions"`
	CurrentQuestion int              `json:"current_question"`

	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Player
	// Unregister requests from clients.
	unregister chan *Player
}

type JSONRoom struct {
	ID              string           `json:"id"`
	Players         []*JSONPlayer        `json:"players"`
	Questions       []Question       `json:"questions"`
	CurrentQuestion int              `json:"current_question"`
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
	}
}

func (r *Room) ToJSON() []byte {
	jsonRoom := &JSONRoom{ID: r.ID, Players: []*JSONPlayer{}, Questions: r.Questions, CurrentQuestion: r.CurrentQuestion}

	for player := range r.Players {
		jsonRoom.Players = append(jsonRoom.Players, player.ToJSONPlayer())
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
}

func (r *Room) RemovePlayer(player *Player) {
	if _, ok := r.Players[player]; ok {
		delete(r.Players, player)
		close(player.send)
		fmt.Println("Removed " + player.Name + " from room " + r.ID)
	}
}

func (r *Room) NextQuestion() *Question {
	if r.CurrentQuestion >= len(r.Questions)-1 {
		return nil
	}
	r.CurrentQuestion++

	return &r.Questions[r.CurrentQuestion]
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
