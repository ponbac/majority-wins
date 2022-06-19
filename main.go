package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ponbac/majority-wins/data"
	"github.com/ponbac/majority-wins/game"
)

// Holds all rooms, key = room ID, value = room pointer
var rooms = map[string]*game.Room{}

func createRoom(c echo.Context) error {
	name := c.QueryParam("name")
	var roomID string
	for ok := true; ok; _, ok = rooms[roomID] {
		roomID = randomString(4)
	}

	room := game.NewRoom(roomID)
	rooms[roomID] = room
	room.Questions = data.FetchQuestions()
	nQuestions := c.QueryParam("questions")
	if nQuestions != "" {
		n, err := strconv.Atoi(nQuestions)
		if err != nil {
			log.Println(err)
		} else {
			room.NQuestions = n
		}
	}
	go room.Run()
	//log.Println("Created room " + roomID)
	err := game.ServeWs(room, true, name, c.Response(), c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Created room "+roomID)
}

func joinRoom(c echo.Context) error {
	// Prevent user to join finished room
	deleteFinishedRooms()

	// Check if room exists
	roomID := c.QueryParam("room")
	if _, ok := rooms[roomID]; !ok {
		return c.String(http.StatusNotFound, "Room "+roomID+" not found")
	}

	name := c.QueryParam("name")
	room := rooms[roomID]
	err := game.ServeWs(room, false, name, c.Response(), c.Request())
	if err != nil {
		// Most probably non unique name used
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Joined room "+roomID)
}

func index(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	return c.String(http.StatusOK, "Hello, Bobba!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}", "uri":"${uri}",` +
			` "status":"${status}", "remote_ip":"${remote_ip}", "latency":"${latency_human}",` +
			` "user_agent":"${user_agent}"}` + "\n",
	}))
	e.Use(middleware.Recover())

	e.Static("/", "./public")

	e.GET("/", index)
	e.GET("/new", createRoom)
	e.GET("/join", joinRoom)

	e.Logger.Fatal(e.Start(":" + port))

	// http.HandleFunc("/", index)
	// http.HandleFunc("/new", createRoom)
	// http.HandleFunc("/join", joinRoom)

	// log.Println("Starting server on http://localhost" + ":" + port)
	// err := http.ListenAndServe(":"+port, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }
}

// TODO: This does not work?
func deleteFinishedRooms() {
	for roomID, room := range rooms {
		if !room.Active {
			log.Println("Deleted room " + roomID)
			delete(rooms, roomID)
		} else {
			//log.Println("Room " + roomID + " is active")
		}
	}
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
