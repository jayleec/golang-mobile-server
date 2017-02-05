package main

import (
	"time"
	"github.com/gorilla/websocket"
	"github.com/jayleec/GolangTraining/27_code-in-process/66_authentication_OAUTH/03_oauth-github/06-complete"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second
	// Time allowed to read next pong message from the peer
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var(
	newline = []byte{'\n'}
	space	= []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:	1024,
	WriteBufferSize:1024,
}

type Client struct {
	hub *githubexample.CommitStats

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}