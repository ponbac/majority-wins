package game

import "github.com/gorilla/websocket"

type Room struct {
	ID              string     `json:"id"`
	Players         []Player   `json:"players"`
	Questions       []Question `json:"questions"`
	CurrentQuestion int        `json:"current_question"`
}

func (r *Room) AddPlayer(player Player) {
	r.Players = append(r.Players, player)
}

func (r *Room) RemovePlayer(playerConnection *websocket.Conn) bool {
	for i, other := range r.Players {
		if other.Conn == playerConnection {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return true
		}
	}

	return false
}

func (r *Room) NextQuestion() *Question {
	if r.CurrentQuestion >= len(r.Questions)-1 {
		return nil
	}
	r.CurrentQuestion++

	return &r.Questions[r.CurrentQuestion]
}