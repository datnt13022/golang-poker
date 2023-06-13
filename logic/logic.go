package logic

import (
	"fmt"
	"math/rand"
	"poker/player"
	"sync"
	"time"
)

// Game represents a multiplayer poker game
type Game struct {
	players        []*player.Player
	communityCards []player.Card
	mutex          sync.Mutex
}

// NewGame creates a new instance of Game with the given players
func NewGame(playerNames []string, startingChips int) *Game {
	players := make([]*player.Player, len(playerNames))
	for i, name := range playerNames {
		players[i] = player.NewPlayer(name, startingChips)
	}

	return &Game{
		players:        players,
		communityCards: make([]player.Card, 0),
	}
}

// PlayGame plays a round of poker game
func (g *Game) PlayGame() {
	deck := CreateDeck()
	dealHands(g.players, deck)
	dealFlop(deck, &g.communityCards)
	dealTurnOrRiver(deck, &g.communityCards)
	dealTurnOrRiver(deck, &g.communityCards)

	for _, p := range g.players {
		fmt.Printf("%s's hand: %v\n", p.Name, p.Hand)
	}
	fmt.Printf("Community cards: %v\n", g.communityCards)

	winners := determineWinner(g.players, g.communityCards)
	if len(winners) == 1 {
		fmt.Printf("Winner: %s\n", winners[0].Name)
	} else {
		fmt.Printf("Winners: ")
		for i, winner := range winners {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(winner.Name)
		}
		fmt.Println()
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

// CreateDeck creates a standard deck of 52 playing cards
func CreateDeck() []player.Card {
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}

	deck := make([]player.Card, 0)
	for _, rank := range ranks {
		for _, suit := range suits {
			card := player.Card{Rank: rank, Suit: suit}
			deck = append(deck, card)
		}
	}
	rand.NewSource(time.Now().UnixNano())
	// rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

// dealHands deals two cards to each player
func dealHands(players []*player.Player, deck []player.Card) {
	for i := 0; i < 2; i++ {
		for _, p := range players {
			card := deck[0]
			p.ReceiveCard(card)
			deck = deck[1:]
		}
	}
}

// dealFlop deals the flop (three community cards)
func dealFlop(deck []player.Card, communityCards *[]player.Card) {
	for i := 0; i < 3; i++ {
		card := deck[0]
		*communityCards = append(*communityCards, card)
		deck = deck[1:]
	}
}

// dealTurnOrRiver deals the turn or river (one community card)
func dealTurnOrRiver(deck []player.Card, communityCards *[]player.Card) {
	card := deck[0]
	*communityCards = append(*communityCards, card)
	deck = deck[1:]
}

// determineWinner determines the winner(s) based on the players' hands and community cards
func determineWinner(players []*player.Player, communityCards []player.Card) []*player.Player {
	winners := []*player.Player{}

	maxRank := 0
	for _, p := range players {
		hand := append(p.Hand, communityCards...)
		rank := evaluateHand(hand)
		if rank > maxRank {
			winners = []*player.Player{p}
			maxRank = rank
		} else if rank == maxRank {
			winners = append(winners, p)
		}
	}

	return winners
}

// evaluateHand evaluates the rank of a hand
func evaluateHand(hand []player.Card) int {
	if isRoyalFlush(hand) {
		return 10
	}
	if isStraightFlush(hand) {
		return 9
	}
	if isFourOfAKind(hand) {
		return 8
	}
	if isFullHouse(hand) {
		return 7
	}
	if isFlush(hand) {
		return 6
	}
	if isStraight(hand) {
		return 5
	}
	if isThreeOfAKind(hand) {
		return 4
	}
	if isTwoPair(hand) {
		return 3
	}
	if isOnePair(hand) {
		return 2
	}
	// High Card
	return 1
}

// isRoyalFlush checks if the hand contains a royal flush
func isRoyalFlush(hand []player.Card) bool {
	// Royal flush: A, K, Q, J, 10 of the same suit
	suits := make(map[string]bool)
	for _, card := range hand {
		suits[card.Suit] = true
	}
	if len(suits) != 1 {
		return false
	}

	ranks := make(map[string]bool)
	for _, card := range hand {
		ranks[card.Rank] = true
	}
	if ranks["A"] && ranks["K"] && ranks["Q"] && ranks["J"] && ranks["10"] {
		return true
	}

	return false
}

// isStraightFlush checks if the hand contains a straight flush
func isStraightFlush(hand []player.Card) bool {
	if !isFlush(hand) {
		return false
	}

	return isStraight(hand)
}

// isFourOfAKind checks if the hand contains four cards of the same rank
func isFourOfAKind(hand []player.Card) bool {
	ranks := make(map[string]int)
	for _, card := range hand {
		ranks[card.Rank]++
	}
	for _, count := range ranks {
		if count == 4 {
			return true
		}
	}
	return false
}

// isFullHouse checks if the hand contains a full house
func isFullHouse(hand []player.Card) bool {
	return isThreeOfAKind(hand) && isOnePair(hand)
}

// isFlush checks if the hand contains five cards of the same suit
func isFlush(hand []player.Card) bool {
	suits := make(map[string]int)
	for _, card := range hand {
		suits[card.Suit]++
	}
	for _, count := range suits {
		if count == 5 {
			return true
		}
	}
	return false
}

// isStraight checks if the hand contains five cards in consecutive ranks
func isStraight(hand []player.Card) bool {
	ranks := make(map[string]bool)
	for _, card := range hand {
		ranks[card.Rank] = true
	}
	if len(ranks) != 5 {
		return false
	}

	rankValues := map[string]int{
		"A":  14,
		"K":  13,
		"Q":  12,
		"J":  11,
		"10": 10,
		"9":  9,
		"8":  8,
		"7":  7,
		"6":  6,
		"5":  5,
		"4":  4,
		"3":  3,
		"2":  2,
	}

	var maxRank, minRank int
	for _, card := range hand {
		rankValue := rankValues[card.Rank]
		if maxRank == 0 || rankValue > maxRank {
			maxRank = rankValue
		}
		if minRank == 0 || rankValue < minRank {
			minRank = rankValue
		}
	}

	return maxRank-minRank == 4
}

// isThreeOfAKind checks if the hand contains three cards of the same rank
func isThreeOfAKind(hand []player.Card) bool {
	ranks := make(map[string]int)
	for _, card := range hand {
		ranks[card.Rank]++
	}
	for _, count := range ranks {
		if count == 3 {
			return true
		}
	}
	return false
}

// isTwoPair checks if the hand contains two pairs
func isTwoPair(hand []player.Card) bool {
	pairCount := 0
	ranks := make(map[string]int)
	for _, card := range hand {
		ranks[card.Rank]++
		if ranks[card.Rank] == 2 {
			pairCount++
		}
	}
	return pairCount == 2
}

// isOnePair checks if the hand contains a single pair
func isOnePair(hand []player.Card) bool {
	ranks := make(map[string]int)
	for _, card := range hand {
		ranks[card.Rank]++
	}
	for _, count := range ranks {
		if count == 2 {
			return true
		}
	}
	return false
}

// func main() {
// 	playerNames := []string{"Alice", "Bob", "Charlie"}
// 	startingChips := 100

// 	game := NewGame(playerNames, startingChips)
// 	game.PlayGame()
// }
