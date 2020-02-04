package core

import (
	"errors"
	"math/rand"
	"time"
)

const (
	numPlayers     = 2
	initialPTokens = 5

	initialWondersPerPlayer = 4
	initialWonders          = initialWondersPerPlayer * numPlayers
)

// Game - game state
type Game struct {
	state State

	players       [numPlayers]Player
	currentPlayer PlayerIndex

	military Military

	availablePTokens []PTokenName
	restPTokens      []PTokenName

	availableWonders []WonderName
	restWonders      []WonderName

	deskAgeI deskAgeI

	// log *logrus.Logger

	rnd *rand.Rand
}

type Option interface {
	applyGame(*Game)
}

func WithSeed(s int64) Option { return seed(s) }

type seed int64

func (s seed) applyGame(g *Game) {
	g.rnd = rand.New(rand.NewSource(int64(s)))
}

// NewGame - init new game with options
func NewGame(opts ...Option) *Game {
	var g = Game{}

	for _, opt := range opts {
		opt.applyGame(&g)
	}

	if g.rnd == nil {
		g.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// init PTokens
	ptokens := NewPTokenNames(g.rnd)
	g.availablePTokens = ptokens[:initialPTokens]
	g.restPTokens = ptokens[initialPTokens:]

	wonders := NewWonderNames(g.rnd)
	g.availableWonders = wonders[:initialWonders]
	g.restWonders = wonders[initialWonders:]

	var cards = IDsAgeI
	g.rnd.Shuffle(numAgeI, func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	g.deskAgeI = newDeskAgeI(cards[:])

	for i := 0; i < numPlayers; i++ {
		g.players[i] = NewPlayer()
	}

	g.state = g.state.Next()
	return &g
}

// GetState of the game
func (g *Game) GetState() State {
	return g.state
}

// Errors
var (
	ErrWrongState           = errors.New("wrong state")
	ErrWrongSelectedWonders = errors.New("wrong selected wonders")
)

// Init of a game
func (g *Game) Init() (wonders [initialWonders]WonderName, ptokens [initialPTokens]PTokenName, ok bool) {
	if !g.state.Is(StateInit) {
		return
	}
	copy(wonders[:], g.availableWonders[:initialWonders])
	ok = true
	g.state = g.state.Next()
	return
}

// SelectWonders as part of an initialize of a game
func (g *Game) SelectWonders(fstWonders, sndWonders [initialWondersPerPlayer]WonderName) error {
	if !g.state.Is(StateSelectWonders) {
		return ErrWrongState
	}
	if !WonderNames(g.availableWonders[:4]).IsExistsAll(WonderNames{fstWonders[0], fstWonders[1], sndWonders[0], sndWonders[1]}) {
		return ErrWrongSelectedWonders
	}
	if !WonderNames(g.availableWonders[4:8]).IsExistsAll(WonderNames{fstWonders[2], fstWonders[3], sndWonders[2], sndWonders[3]}) {
		return ErrWrongSelectedWonders
	}

	g.players[0].AvailableWorneds = fstWonders[:]
	g.players[1].AvailableWorneds = sndWonders[:]
	g.state = g.state.Next()
	return nil
}

func (g *Game) AgeI() [SizeAge]DeskCard {
	return g.deskAgeI
}

type DeskCard struct {
	IsVisible   bool
	IsSkipped   bool
	IsAvailable bool
	ID          CardID
}

// Cost card by coins
func (g *Game) Cost(card CardName) Coins {
	return card.card().Cost.ByCoins(g, g.currentPlayer)
}

func (g *Game) Player(i PlayerIndex) Player {
	return g.players[i]
}

func (g *Game) player(i PlayerIndex) *Player {
	return &g.players[i]
}

func (g *Game) apply(card CardName) {
	for _, e := range card.card().Effects {
		e.Apply(g, g.currentPlayer)
	}
}

// State of a game
type State uint8

// Game states
const (
	StateNone State = iota
	StateInit
	StateSelectWonders
	StateAgeI
)

// Is has that current state
func (s State) Is(check State) bool {
	return s == check
}

// Next ...
func (s State) Next() State {
	switch s {
	case StateNone:
		return StateInit
	case StateInit:
		return StateSelectWonders
	case StateSelectWonders:
		return StateAgeI
	default:
		panic("unknown state")
	}
}

// type gameState struct {
// 	g *Game
// }

// func (gs *gameState) Player() *Player {
// 	return gs.g.player(gs.g.currentPlayer)
// }

// func (gs *gameState) Opponent() *Player {
// 	return gs.g.player(gs.g.currentPlayer.Next())
// }
