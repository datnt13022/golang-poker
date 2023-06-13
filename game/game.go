package game

import (
	"poker/player"
	"sync"
)

// Game represents a game instance
type Game struct {
	players []*player.Player
	mutex   sync.Mutex
}

// NewGame creates a new instance of Game
func NewGame() *Game {
	return &Game{
		players: make([]*player.Player, 0),
	}
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(p *player.Player) {
	g.mutex.Lock()
	g.players = append(g.players, p)
	g.mutex.Unlock()
}

// RemovePlayer removes a player from the game
func (g *Game) RemovePlayer(p *player.Player) {
	g.mutex.Lock()
	for i, player := range g.players {
		if player == p {
			g.players = append(g.players[:i], g.players[i+1:]...)
			break
		}
	}
	g.mutex.Unlock()
}

// GetPlayers returns the list of players in the game
func (g *Game) GetPlayers() []*player.Player {
	g.mutex.Lock()
	players := make([]*player.Player, len(g.players))
	copy(players, g.players)
	g.mutex.Unlock()
	return players
}
