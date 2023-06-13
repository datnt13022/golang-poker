package game

import (
	"fmt"
	"poker/player"
	"sync"
)

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
func (g *Game) PlayPoker() {
	// Initialize the game
	deck := logic.createDeck()
	pot := 0
	bigBlind := 10
	smallBlind := bigBlind / 2
	dealerIndex := 0
	currentBet := bigBlind
	communityCards := make([]Card, 0)
	activePlayers := make([]*player.Player, len(g.players))
	copy(activePlayers, g.players)

	// Main game loop
	for len(activePlayers) > 1 {
		// Reset the game state for a new round
		resetGame(&deck, &pot, &communityCards, &activePlayers)
		rotateDealer(&dealerIndex, len(activePlayers))
		postBlinds(&activePlayers, dealerIndex, smallBlind, bigBlind)
		dealCards(&deck, &activePlayers)

		// Pre-flop
		bettingRound(&activePlayers, dealerIndex, currentBet)

		// Flop
		revealCommunityCards(&deck, &communityCards, 3)
		bettingRound(&activePlayers, dealerIndex, currentBet)

		// Turn
		revealCommunityCards(&deck, &communityCards, 1)
		bettingRound(&activePlayers, dealerIndex, currentBet)

		// River
		revealCommunityCards(&deck, &communityCards, 1)
		bettingRound(&activePlayers, dealerIndex, currentBet)

		// Showdown
		winners := determineWinners(&activePlayers, &communityCards)
		distributePot(&pot, winners)

		// Remove players with no chips left
		activePlayers = removeBustedPlayers(&activePlayers)
	}

	// Declare the winner
	if len(activePlayers) == 1 {
		fmt.Printf("Player %s wins the game!\n", activePlayers[0].Name)
	} else {
		fmt.Println("No active players remaining. The game ends in a draw.")
	}
}
