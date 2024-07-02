package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	clients    []*ClientConnection
	rooms      map[int]*Room
	roomsToRun chan *Room
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make([]*ClientConnection, 0),
		rooms:      make(map[int]*Room),
		roomsToRun: make(chan *Room),
	}
}

func (s *Server) GetRoom(roomID int) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, exists := s.rooms[roomID]; exists {
		log.Printf("Devolviendo Room ya existente no: %d", roomID)
		return room
	}
	newRoom := &Room{
		ID:       roomID,
		join:     make(chan *User),
		leave:    make(chan *User),
		messages: make(chan string),
	}
	log.Printf("Creando Room %d", roomID)
	s.rooms[roomID] = newRoom
	s.roomsToRun <- newRoom
	return newRoom
}

func (s *Server) RunRooms() {
	for {
		select {
		case room := <-s.roomsToRun:
			if !room.state {
				log.Printf("Corriendo Room %d", room.ID)
				s.run(room.ID)
			}
		}
	}
}

func (s *Server) AcceptConnection(conn net.Conn) {
	client := NewClientConnection(conn, s)
	s.mu.Lock()
	s.clients = append(s.clients, client)
	s.mu.Unlock()
	go client.HandleConnection()
}

func (s *Server) BroadcastMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		client.user.Messages <- message
	}
}

func (s *Server) SendMessageToUser(userID, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		if client.user.ID == userID {
			client.user.Messages <- message
			return nil
		}
	}
	return fmt.Errorf("user with ID %s not found", userID)
}

func (s *Server) ListAllUsers() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		fmt.Printf("userID: %s, username: %s, room: %d\n", client.user.ID, client.user.Name, client.user.RoomID)
	}
}

func (s *Server) run(roomID int) {
	r := s.GetRoom(roomID)
	r.state = true
	for {
		select {
		case user := <-r.join:
			for _, client := range s.clients {
				if client.user.RoomID == roomID && client.user.ID != user.ID {
					client.user.Messages <- fmt.Sprintf("New User is here, name %s\n", user.Name)
				}
			}
		case user := <-r.leave:
			for _, client := range s.clients {
				if client.user.RoomID == roomID {
					client.user.Messages <- fmt.Sprintf("%s said goodbye!", user.Name)
				}
			}
		case message := <-r.messages:
			for _, client := range s.clients {
				if client.user.RoomID == roomID {
					client.user.Messages <- message
				}
			}
		}
	}
}
