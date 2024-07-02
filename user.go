package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

type User struct {
	ID       string
	Name     string
	RoomID   int
	Messages chan string
	conn     *net.Conn
}

func NewUser(name string, roomID int, conn *net.Conn) *User {
	log.Printf("Creando usuario %s", name)
	return &User{
		ID:       uuid.New().String(),
		Name:     name,
		RoomID:   roomID,
		Messages: make(chan string),
		conn:     conn,
	}
}

func (u *User) ReceiveMessage() {
	log.Printf("Activando ReceiveMessage")
	for message := range u.Messages {
		fmt.Fprintln(*u.conn, message)
	}
}
