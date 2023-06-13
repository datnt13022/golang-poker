package player

// Card represents a card in the game
type Card struct {
	Rank string
	Suit string
}

// Player represents a player in the game
type Player struct {
	Name   string
	Chips  int
	Hand   []Card
	Folded bool
}

// NewPlayer creates a new instance of Player
func NewPlayer(name string, chips int) *Player {
	return &Player{
		Name:   name,
		Chips:  chips,
		Hand:   make([]Card, 0),
		Folded: false,
	}
}

// ReceiveCard receives a card and adds it to the player's hand
func (p *Player) ReceiveCard(card Card) {
	p.Hand = append(p.Hand, card)
}

// ResetHand resets the player's hand
func (p *Player) ResetHand() {
	p.Hand = make([]Card, 0)
	p.Folded = false
}

// Bet makes a bet and deducts the bet amount from the player's chips
func (p *Player) Bet(amount int) int {
	if amount > p.Chips {
		amount = p.Chips
	}
	p.Chips -= amount
	return amount
}

// WinPot adds the winnings to the player's chips
func (p *Player) WinPot(amount int) {
	p.Chips += amount
}

// Fold marks the player as folded
func (p *Player) Fold() {
	p.Folded = true
}

// Folded returns true if the player has folded, false otherwise
func (p *Player) isFolded() bool {
	return p.Folded
}

// Chips returns the number of chips the player has
func (p *Player) GetChips() int {
	return p.Chips
}

// GetHand returns the player's hand
func (p *Player) GetHand() []Card {
	return p.Hand
}
