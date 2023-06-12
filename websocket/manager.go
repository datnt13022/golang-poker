package websocket

import (
	"fmt"
	"log"
	"net/http"
	"poker/game"
	"poker/player"
	"poker/user"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager handles WebSocket connections and game logic
type Manager struct {
	game  *game.Game
	users user.UserMap
	mutex sync.Mutex
}

// NewManager creates a new instance of Manager
func NewManager(game *game.Game, users user.UserMap) *Manager {
	return &Manager{
		game:  game,
		users: users,
	}
}

// HandleWebSocket handles WebSocket connections
func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Perform user authentication
	err = m.authenticate(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Create a new player
	player := player.NewPlayer()

	// Add the player to the game
	m.mutex.Lock()
	m.game.AddPlayer(player)
	m.mutex.Unlock()

	// Game loop
	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Handle player actions based on the received message
		m.handlePlayerAction(player, string(msg))

		// Update game state

		// Send game state updates to all players
		m.mutex.Lock()
		players := m.game.GetPlayers()
		for _, p := range players {
			err := conn.WriteMessage(websocket.TextMessage, []byte(m.getGameState(p)))
			if err != nil {
				log.Println(err)
			}
		}
		m.mutex.Unlock()
	}

	// Remove the player from the game
	m.mutex.Lock()
	m.game.RemovePlayer(player)
	m.mutex.Unlock()
}

// Handle user authentication
func (m *Manager) authenticate(conn *websocket.Conn) error {
	// Read username and password from client
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	credentials := string(msg)
	// Split credentials into username and password
	// You can use a more secure method to store and compare passwords, like bcrypt
	username, password := extractCredentials(credentials)

	// Check if the user exists and the password is correct
	m.mutex.Lock()
	user, ok := m.users[username]
	m.mutex.Unlock()

	if !ok || user.Password != password {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

// Handle player actions
func (m *Manager) handlePlayerAction(player *player.Player, action string) {
	// Implement logic to handle player actions here
}

// Get game state as a string
func (m *Manager) getGameState(player *player.Player) string {
	// Implement logic to generate game state here
	return ""
}
