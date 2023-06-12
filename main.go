package main

import (
	"fmt"
	"log"
	"net/http"
	"poker/game"
	"poker/user"
	"poker/websocket"
)

func main() {
	// Initialize game state
	game := game.NewGame()

	// Initialize users map
	users := user.NewUserMap()
	users.AddUser("username", "password")

	// Create WebSocket manager
	manager := websocket.NewManager(game, users)

	// Define WebSocket route
	http.HandleFunc("/ws", manager.HandleWebSocket)

	// Start the server
	fmt.Println("Poker game server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
