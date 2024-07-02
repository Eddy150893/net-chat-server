package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type ClientConnection struct {
	user   *User
	server *Server
}

func NewClientConnection(conn net.Conn, server *Server) *ClientConnection {
	return &ClientConnection{
		user:   &User{conn: &conn},
		server: server,
	}
}

func (c *ClientConnection) HandleConnection() {
	defer (*c.user.conn).Close()

	inputMessage := bufio.NewScanner(*c.user.conn)
	for inputMessage.Scan() {
		message := inputMessage.Text()
		if strings.HasPrefix(message, "MESSAGE:") {
			content := strings.TrimPrefix(message, "MESSAGE: ")
			c.handleMessage(content)
		} else if strings.HasPrefix(message, "USERNAME:") && strings.Contains(message, " ROOM:") {
			c.handleUserAndRoom(message)
		} else {
			log.Printf("Unknown message format: %s", message)
		}
	}
	room := c.server.GetRoom(c.user.RoomID)
	room.leave <- c.user
}

func (c *ClientConnection) handleMessage(content string) {
	room := c.server.GetRoom(c.user.RoomID)
	room.messages <- fmt.Sprintf("%s: %s\n", c.user.Name, content)
}

func (c *ClientConnection) handleUserAndRoom(message string) {
	parts := strings.Split(message, " ROOM:")
	if len(parts) != 2 {
		log.Printf("Invalid USERNAME/ROOM message format: %s", message)
		return
	}

	username := strings.TrimPrefix(parts[0], "USERNAME: ")
	roomIDStr := parts[1]
	roomID, err := strconv.Atoi(strings.TrimSpace(roomIDStr))
	if err != nil {
		log.Printf("Invalid room ID: %s", roomIDStr)
		return
	}

	c.user = NewUser(username, roomID, c.user.conn)
	room := c.server.GetRoom(c.user.RoomID)

	room.join <- c.user

	go c.user.ReceiveMessage()
	c.user.Messages <- fmt.Sprintf("Welcome %s to room %d", c.user.Name, c.user.RoomID)
}
