package main

import (
	"./trace"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// Join into This Room
			r.clients[client] = true
			r.tracer.Trace("A new client enterd")
		case client := <-r.leave:
			// Get Out of This Room
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("A client left this room")
		case msg := <-r.forward:
			r.tracer.Trace("A message received :", string(msg))
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- Send to Client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- Sending Failed. Clean up the client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
