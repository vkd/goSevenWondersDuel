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
	state      State
	currentAge uint8
	winner     Winner
	winReason  WinReason

	discardOpponentBuild discardOpponentBuild

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

	availablePTokens []PTokenID
	restPTokens      []PTokenID

	availableWonders [initialWonders]WonderID
	restWonders      []WonderID

	ageI     [SizeAge]CardID
	ageII    [SizeAge]CardID
	ageIII   [SizeAge]CardID
	ageDesk  ageDesk
	ageDesk2 ageDesk
	ageDesk3 ageDesk

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
	ptokens := shufflePTokens(g.rnd)
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
	g.ageDesk2, err = newAgeDesk(structureAgeII, shuffleAgeII(g.rnd))
	if err != nil {
		return nil, err
	}
	g.ageDesk3, err = newAgeDesk(structureAgeIII, shuffleAgeIII(g.rnd))
	if err != nil {
		return nil, err
	}

	for i := 0; i < numPlayers; i++ {
		g.players[i] = NewPlayer()
	}

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
	if !g.state.Is(StateNone) {
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

	var available [numWonders]uint8
	for _, w := range g.availableWonders {
		available[w]++
	}
	for _, fst := range fstWonders {
		if available[fst] > 0 {
			available[fst]--
		} else {
			return ErrWrongSelectedWonders
		}
	}
	for _, snd := range sndWonders {
		if available[snd] > 0 {
			available[snd]--
		} else {
			return ErrWrongSelectedWonders
		}
	}

	for _, a := range available {
		if a != 0 {
			return ErrWrongSelectedWonders
		}
	}

	g.wondersPerPlayer[0] = fstWonders[:]
	g.wondersPerPlayer[1] = sndWonders[:]
	g.state = g.state.Next()
	return nil
}

func (g *Game) CardsState() CardsState {
	return g.ageDesk.state
}

func (g *Game) CardCost(id CardID) Coins {
	return g.costForCurrentPlayer(id.card().Cost)
}

func (g *Game) WonderCost(id WonderID) Coins {
	return g.costForCurrentPlayer(id.wonder().Cost)
}

func (g *Game) costForCurrentPlayer(cost Cost) Coins {
	return CostByCoins(
		cost,
		*g.currentPlayer(),
		NewTradingPrice(
			*g.opponent(),
			g.priceMarkets[g.currentPlayerIndex]...,
		),
		g.oneAnyMarkets[g.currentPlayerIndex],
	)
}

func (g *Game) ConstructBuilding(id CardID) (state CardsState, err error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

	ok := g.ageDesk.testBuild(id)
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
	g.nextTurn()

	return state, nil
}

func (g *Game) buildCard(id CardID) {
	card := id.card()
	for _, eff := range card.Effects {
		eff.applyEffect(g, g.currentPlayerIndex)
	}
	g.builtCards[g.currentPlayerIndex][card.Color] = append(g.builtCards[g.currentPlayerIndex][card.Color], id)
}

func (g *Game) DiscardCard(id CardID) (state CardsState, _ error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

	err := g.ageDesk.Build(id)
	state = g.ageDesk.state
	if err != nil {
		return state, err
	}

	g.discardedCards = append(g.discardedCards, id)

	g.currentPlayer().Coins += Coins(2) + Coins(len(g.builtCards[g.currentPlayerIndex][Yellow]))

	g.nextTurn()
	return state, nil
}

func (g *Game) ConstructWonder(cid CardID, wid WonderID) (state CardsState, err error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

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
	g.nextTurn()

	return state, nil
}

func (g *Game) ChoosePToken(id PTokenID) error {
	if !g.state.Is(StateChoosePToken) {
		return ErrWrongState
	}

	panic("Not implemented")

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) ConstructDiscardedCard(id CardID) (err error) {
	if !g.state.Is(StateBuildFreeDiscarded) {
		return ErrWrongState
	}

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
	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) DiscardOpponentBuild(id CardID) error {
	if !g.state.Is(StateDiscardOpponentBuild) {
		return ErrWrongState
	}

	card := id.card()

	if card.Color != g.discardOpponentBuild.color {
		return fmt.Errorf("wrong card id: wrong color")
	}

	opponentID := g.CurrentPlayerIndex().Next()

	var ok bool
	for _, b := range g.builtCards[opponentID][card.Color] {
		if b == id {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("wrong card id: %d is not built on opponent side", id)
	}

	var newBuiltCards []CardID
	for _, b := range g.builtCards[opponentID][card.Color] {
		if b != id {
			newBuiltCards = append(newBuiltCards, b)
		}
	}
	g.builtCards[opponentID][card.Color] = newBuiltCards

	card.discard(g, opponentID)

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) PlayDiscardedPToken(id PTokenID) (err error) {
	if !g.state.Is(StateBuildFreePToken) {
		return ErrWrongState
	}

	panic("Not implemented")

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) nextTurn() {
	if !g.state.Is(StateGameTurn) {
		return
	}
	if g.ageDesk.state.anyExists() {
		if g.repeatTurn {
			g.repeatTurn = false
		} else {
			g.currentPlayerIndex = g.currentPlayerIndex.Next()
		}
	} else {
		if g.currentAge == 2 {
			winner := g.finalVPs()
			g.victory(winner, WinCivilian)
			return
		}
		g.repeatTurn = false
		g.nextAge()
	}
}

func (g *Game) nextAge() {
	switch {
	case g.military.Shields[0] > g.military.Shields[1]:
		g.currentPlayerIndex = 1
		g.state = StateChooseFirstPlayer
	case g.military.Shields[0] < g.military.Shields[1]:
		g.currentPlayerIndex = 0
		g.state = StateChooseFirstPlayer
	default:
		// same player starts next age
	}

	switch g.currentAge {
	case 0:
	case 1:
		g.ageDesk = g.ageDesk2
	case 2:
		g.ageDesk = g.ageDesk3
	}
	g.currentAge++
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

func (g *Game) finalVPs() Winner {
	for pi, fs := range g.endEffects {
		for _, f := range fs {
			vp := f.finalVP(g, PlayerIndex(pi))
			g.vps[pi][vp.t] += vp.v
		}
	}

	var score [numPlayers]VP
	for pi, vs := range g.vps {
		for _, v := range vs {
			score[pi] += v
		}
	}

	switch {
	case score[0] > score[1]:
		return Winner1Player
	case score[1] > score[0]:
		return Winner2Player
	}

	switch {
	case g.vps[0][BlueVP] > g.vps[1][BlueVP]:
		return Winner1Player
	case g.vps[1][BlueVP] > g.vps[0][BlueVP]:
		return Winner2Player
	}

	return WinnerBoth
}

func (g *Game) victory(w Winner, reason WinReason) {
	g.state = StateVictory
	g.winner = w
	g.winReason = reason
}

type WinReason uint8

const (
	WinNone WinReason = iota
	WinCivilian
	WinMilitary
	WinScience
)

type Winner uint8

const (
	Winner1Player Winner = iota
	Winner2Player
	WinnerBoth
	numWinners
)

var _ = [1]struct{}{}[numWinners-1-numPlayers]

func (g *Game) gettingPToken(i PlayerIndex) {
	g.state = StateChoosePToken
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
	StateSelectWonders
	StateGameTurn
	StateBuildFreeDiscarded
	StateDiscardOpponentBuild
	StateBuildFreePToken
	StateChoosePToken
	StateChooseFirstPlayer
	StateVictory
	numStates int = iota
)

var (
	stateNames = map[State]string{
		StateNone:                 "None",
		StateSelectWonders:        "SelectWonders",
		StateGameTurn:             "GameTurn",
		StateBuildFreeDiscarded:   "BuildFreeDiscarded",
		StateDiscardOpponentBuild: "DiscardOpponentBuild",
		StateBuildFreePToken:      "BuildFreePToken",
		StateChoosePToken:         "ChoosePToken",
		StateChooseFirstPlayer:    "ChooseFirstPlayer",
		StateVictory:              "Victory",
	}
	_ = [1]struct{}{}[len(stateNames)-numStates]
)

// Is has that current state
func (s State) Is(check State) bool {
	return s == check
}

func (s State) String() string {
	return stateNames[s]
}

// Next ...
func (s State) Next() State {
	switch s {
	case StateNone:
		return StateSelectWonders
	case StateSelectWonders:
		return StateGameTurn
	case StateGameTurn:
	case StateBuildFreeDiscarded:
		return StateGameTurn
	case StateDiscardOpponentBuild:
		return StateGameTurn
	case StateChoosePToken:
		return StateGameTurn
	case StateBuildFreePToken:
		return StateGameTurn
	case StateVictory:
	default:
		panic("unknown state")
	}
	return s
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
	g.state = StateBuildFreeDiscarded
}

func DiscardOpponentBuild(color CardColor) Effect {
	return discardOpponentBuild{color}
}

type discardOpponentBuild struct {
	color CardColor
}

var _ Effect = discardOpponentBuild{}

func (c discardOpponentBuild) applyEffect(g *Game, i PlayerIndex) {
	g.discardOpponentBuild = c
	g.state = StateDiscardOpponentBuild
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
	g.state = StateBuildFreePToken
}
