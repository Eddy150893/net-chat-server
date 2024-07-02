package main

type Room struct {
	ID       int
	join     chan *User
	leave    chan *User
	messages chan string
	state    bool
}
