package core

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	numPlayers     = 2
	initialPTokens = 5

	numWondersPerPlayer = 4
	initialWonders      = numWondersPerPlayer * numPlayers
)

// Game - game state
type Game struct {
	state State

	players            [numPlayers]Player
	currentPlayerIndex PlayerIndex
	repeatTurn         bool

	wondersPerPlayer [numPlayers][]WonderID
	buildWonders     [numPlayers][]WonderID
	builtCards       [numPlayers][numCardColors][]CardID
	discardedCards   []CardID

	priceMarkets  [numPlayers]PriceMarkets
	oneAnyMarkets [numPlayers]OneAnyMarkets
	endEffects    [numPlayers][]Finaler

	military Military

	vps [numPlayers][numVPTypes]VP

	availablePTokens []PTokenName
	restPTokens      []PTokenName

	availableWonders [initialWonders]WonderID
	restWonders      []WonderID

	ageI    [SizeAge]CardID
	ageII   [SizeAge]CardID
	ageIII  [SizeAge]CardID
	ageDesk ageDesk

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
func NewGame(opts ...Option) (*Game, error) {
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

	wonders := shuffleWonders(g.rnd)
	copy(g.availableWonders[:], wonders[:initialWonders])
	g.restWonders = wonders[initialWonders:]

	var err error
	g.ageDesk, err = newAgeDesk(structureAgeI, shuffleAgeI(g.rnd))
	if err != nil {
		return nil, err
	}

	for i := 0; i < numPlayers; i++ {
		g.players[i] = NewPlayer()
	}

	g.state = g.state.Next()
	return &g, nil
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
func (g *Game) Init() (wonders [initialWonders]WonderID, ptokens [initialPTokens]PTokenName, ok bool) {
	if !g.state.Is(StateInit) {
		return
	}
	wonders = g.availableWonders
	ok = true
	g.state = g.state.Next()
	return
}

// SelectWonders as part of an initialize of a game
func (g *Game) SelectWonders(fstWonders, sndWonders [numWondersPerPlayer]WonderID) error {
	if !g.state.Is(StateSelectWonders) {
		return ErrWrongState
	}

	// TODO: check input arguments

	// if !WonderNames(g.availableWonders[:4]).IsExistsAll(WonderNames{fstWonders[0], fstWonders[1], sndWonders[0], sndWonders[1]}) {
	// 	return ErrWrongSelectedWonders
	// }
	// if !WonderNames(g.availableWonders[4:8]).IsExistsAll(WonderNames{fstWonders[2], fstWonders[3], sndWonders[2], sndWonders[3]}) {
	// 	return ErrWrongSelectedWonders
	// }

	g.wondersPerPlayer[0] = fstWonders[:]
	g.wondersPerPlayer[1] = sndWonders[:]
	g.state = g.state.Next()
	return nil
}

func (g *Game) CardsState() CardsState {
	return g.ageDesk.state
}

func (g *Game) CardCost(id CardID) Coins {
	return CostByCoins(
		id.card().Cost,
		*g.currentPlayer(),
		NewTradingPrice(
			*g.opponent(),
			g.priceMarkets[g.currentPlayerIndex]...,
		),
		g.oneAnyMarkets[g.currentPlayerIndex],
	)
}

func (g *Game) WonderCost(id WonderID) Coins {
	return CostByCoins(
		id.wonder().Cost,
		*g.currentPlayer(),
		NewTradingPrice(
			*g.opponent(),
			g.priceMarkets[g.currentPlayerIndex]...,
		),
		g.oneAnyMarkets[g.currentPlayerIndex],
	)
}

func (g *Game) ConstructBuilding(id CardID) (state CardsState, err error) {
	ok := g.ageDesk.testBuild(id)
	state = g.ageDesk.state
	if !ok {
		return state, fmt.Errorf("card (id = %d) cannot be built", id)
	}

	pay := g.CardCost(id)
	if g.currentPlayer().Coins < pay {
		return state, fmt.Errorf("not enough coins")
	}

	err = g.ageDesk.Build(id)
	state = g.ageDesk.state
	if err != nil {
		return state, err
	}
	g.currentPlayer().Coins -= pay

	g.buildCard(id)
	g.nextState()

	return state, nil
}

func (g *Game) buildCard(id CardID) {
	card := id.card()
	for _, eff := range card.Effects {
		eff.applyEffect(g, g.currentPlayerIndex)
	}
	g.builtCards[g.currentPlayerIndex][card.Color] = append(g.builtCards[g.currentPlayerIndex][card.Color], id)
}

func (g *Game) DiscardCard(id CardID) (state CardsState, _ bool) {
	err := g.ageDesk.Build(id)
	state = g.ageDesk.state
	if err != nil {
		return state, false
	}

	g.discardedCards = append(g.discardedCards, id)

	g.currentPlayer().Coins += Coins(2) + Coins(len(g.builtCards[g.currentPlayerIndex][Yellow]))

	g.nextState()
	return state, true
}

func (g *Game) ConstructWonder(cid CardID, wid WonderID) (state CardsState, err error) {
	state = g.ageDesk.state
	if len(g.buildWonders[0])+len(g.buildWonders[1]) >= 7 {
		return state, fmt.Errorf("wonder (id = %d) cannot be built: max 7 wonders are allowed", wid)
	}

	ok := g.ageDesk.testBuild(cid)
	state = g.ageDesk.state
	if !ok {
		return state, fmt.Errorf("card (id = %d) cannot be taken", cid)
	}

	pay := g.WonderCost(wid)
	if g.currentPlayer().Coins < pay {
		return state, fmt.Errorf("not enough coins")
	}

	err = g.ageDesk.Build(cid)
	state = g.ageDesk.state
	if err != nil {
		return state, err
	}
	g.currentPlayer().Coins -= pay

	wonder := wid.wonder()
	for _, eff := range wonder.Effects {
		eff.applyEffect(g, g.currentPlayerIndex)
	}

	g.buildWonders[g.currentPlayerIndex] = append(g.buildWonders[g.currentPlayerIndex], wid)
	g.nextState()

	return state, nil
}

func (g *Game) ConstructDiscardedCard(id CardID) (err error) {
	// TODO: check current state
	var ok bool
	for _, dID := range g.discardedCards {
		if dID == id {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("card (id=%d) is not discarded", id)
	}

	var newDiscarded []CardID
	for _, did := range g.discardedCards {
		if did != id {
			newDiscarded = append(newDiscarded, did)
		}
	}
	g.discardedCards = newDiscarded

	g.buildCard(id)
	g.nextState()
	return nil
}

func (g *Game) PlayDiscardedPToken(id PTokenID) (err error) {
	// TODO: check game state
	panic("Not implemented")
}

func (g *Game) nextState() {
	if g.repeatTurn {
		g.repeatTurn = false
	} else {
		g.currentPlayerIndex = g.currentPlayerIndex.Next()
	}
}

func (g *Game) Player(i PlayerIndex) Player {
	return g.players[i]
}

func (g *Game) CurrentPlayer() Player {
	return g.Player(g.currentPlayerIndex)
}

func (g *Game) player(i PlayerIndex) *Player {
	return &g.players[i]
}

func (g *Game) currentPlayer() *Player {
	return &g.players[g.currentPlayerIndex]
}

func (g *Game) opponent() *Player {
	return &g.players[g.currentPlayerIndex.Next()]
}

func (g *Game) apply(card CardName) {
	for _, e := range card.card().Effects {
		e.applyEffect(g, g.currentPlayerIndex)
	}
}

func (g *Game) Victory() {
	panic("Not implemented")
}

func (g *Game) GettingPToken(i PlayerIndex) {
	panic("Not implemented")
}

func (g *Game) Military() Military {
	return g.military
}

func (g *Game) CurrentPlayerIndex() PlayerIndex {
	return g.currentPlayerIndex
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

func zeroRand() *rand.Rand {
	return rand.New(rand.NewSource(0))
}

func RepeatTurn() Effect {
	return repeatTurn{}
}

type repeatTurn struct{}

var _ Effect = repeatTurn{}

func (repeatTurn) applyEffect(g *Game, _ PlayerIndex) {
	g.repeatTurn = true
}

func BuildFreeDiscardedCard() Effect {
	return freeDiscardedCard{}
}

type freeDiscardedCard struct{}

var _ Effect = freeDiscardedCard{}

func (freeDiscardedCard) applyEffect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}

func DiscardOpponentBuild(color CardColor) Effect {
	return discardOpponentBuild{color}
}

type discardOpponentBuild struct {
	color CardColor
}

var _ Effect = discardOpponentBuild{}

func (c discardOpponentBuild) applyEffect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}

func DiscardOpponentCoins(c Coins) Effect {
	return discardOpponentCoins{c}
}

type discardOpponentCoins struct {
	coins Coins
}

var _ Effect = discardOpponentCoins{}

func (c discardOpponentCoins) applyEffect(g *Game, i PlayerIndex) {
	g.opponent().Coins.sub(c.coins)
}

func PlayOneOf3DiscardedPToken() Effect {
	return playDiscardedPToken{}
}

type playDiscardedPToken struct{}

var _ Effect = playDiscardedPToken{}

func (playDiscardedPToken) applyEffect(g *Game, i PlayerIndex) {
	panic("Not implemented")
}
