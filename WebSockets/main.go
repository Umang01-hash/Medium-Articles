package main

import (
	"gofr.dev/pkg/gofr"
)

type Message struct {
	Content string `json:"content"`
	User    string `json:"user"`
}

func main() {
	app := gofr.New()

	app.WebSocket("/chat", ChatHandler)

	app.Run()
}

func ChatHandler(ctx *gofr.Context) (interface{}, error) {
	// Read the incoming message from the client
	var message Message
	err := ctx.Bind(&message)
	if err != nil {
		return nil, err
	}

	// Process the message (e.g., broadcast to all connected clients)
	ctx.Infof("Received message from %s: %s\n", message.User, message.Content)

	// You can send data back to the client here (optional)
	return nil, nil
}
