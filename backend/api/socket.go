package api

import (
	"encoding/json"
	"quandat10/htttdl/backend/utils"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Username  string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type WClient struct {
	username string
	radius   int16
	ws       *websocket.Conn
}

func (s *Server) handleWebSocket(c echo.Context) error {
	username := c.QueryParam("username")
	radiusParam := c.QueryParam("radius")

	radius, err := strconv.ParseInt(radiusParam, 10, 16)
	if err != nil {
		return err
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Register connection with map
	room := s.getMapRoom()
	room.registerRoom(conn, username, int16(radius))
	var msg Message

	// Response list coordinates in the first connection
	room.broadcastRoom(s)

	// Read messages from WebSocket and broadcast to other clients
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			room.unregisterRoom(conn, username)
			return err
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Err(err)
		}

		user := User{
			Name:      msg.Username,
			Longitude: msg.Longitude,
			Latitude:  msg.Latitude,
		}

		err = s.updateStore(user)
		if err != nil {
			log.Err(err)
		}

		room.broadcastRoom(s)
	}
}

// MapRoom represents a Map room with multiple clients
type MapRoom struct {
	clients    map[WClient]bool
	broadcast  chan []byte
	register   chan WClient
	unregister chan WClient
}

// Create a new Map room
func newMapRoom() *MapRoom {
	return &MapRoom{
		clients:    make(map[WClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan WClient),
		unregister: make(chan WClient),
	}
}

// Register a new client with the map room
func (c *MapRoom) registerRoom(conn *websocket.Conn, username string, radius int16) {
	cl := WClient{
		username: username,
		radius:   radius,
		ws:       conn,
	}
	c.clients[cl] = true
}

// Unregister a client with the map room
func (c *MapRoom) unregisterRoom(conn *websocket.Conn, username string) {
	cl := WClient{
		username: username,
		ws:       conn,
	}
	if _, ok := c.clients[cl]; ok {
		delete(c.clients, cl)
		close(connCloseChan(conn))
	}
}

// Get a channel to signal the WebSocket connection's close
func connCloseChan(conn *websocket.Conn) chan bool {
	ch := make(chan bool, 1)
	go func() {
		ch <- true
	}()
	return ch
}

// Broadcast a message to all clients in the map room
func (c *MapRoom) broadcastRoom(s *Server) {

	for conn := range c.clients {
		users, _ := s.listUsersInsideRadius(conn.username, conn.radius)
		rs := utils.AIToAB(users)
		if err := conn.ws.WriteMessage(websocket.TextMessage, rs); err != nil {
			c.unregisterRoom(conn.ws, conn.username)
		}
	}
}

var mapRoom *MapRoom

// Get the map room singleton
func (s *Server) getMapRoom() *MapRoom {
	if mapRoom == nil {
		mapRoom = newMapRoom()
		go mapRoom.run(s)
	}
	return mapRoom
}

// Run the map room event loop
func (c *MapRoom) run(s *Server) {
	for {
		select {
		case conn := <-c.register:
			c.registerRoom(conn.ws, conn.username, conn.radius)
		case conn := <-c.unregister:
			c.unregisterRoom(conn.ws, conn.username)
		case _ = <-c.broadcast:
			c.broadcastRoom(s)
		}
	}
}
