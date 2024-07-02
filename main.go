package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func startCommandListener(server *Server) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		if strings.HasPrefix(command, "/broadcast ") {
			message := strings.TrimPrefix(command, "/broadcast ")
			server.BroadcastMessage(message)
		} else if strings.HasPrefix(command, "/msg ") {
			parts := strings.SplitN(strings.TrimPrefix(command, "/msg "), " ", 2)
			if len(parts) == 2 {
				userID, message := parts[0], parts[1]
				err := server.SendMessageToUser(userID, message)
				if err != nil {
					log.Printf("Error sending message to user: %v", err)
				}
			}
		} else if strings.HasPrefix(command, "/listUsers") {
			server.ListAllUsers()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning input: %v", err)
	}
}

func main() {
	port := flag.Int("p", 3090, "port")
	host := flag.String("h", "localhost", "host")
	flag.Parse()

	server := NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server started on %s:%d", *host, *port)

	go startCommandListener(server)
	go server.RunRooms()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		server.AcceptConnection(conn)
	}
}
