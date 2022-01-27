package game

import "github.com/gorilla/websocket"

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Conn  *websocket.Conn
}