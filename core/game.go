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
	InitialWonders      = numWondersPerPlayer * numPlayers
)

// Game - game state
type Game struct {
	GameState

	state       State
	winner      Winner
	victoryType VictoryType

	discardOpponentBuild discardOpponentBuild

	players            [numPlayers]Player
	currentPlayerIndex PlayerIndex
	repeatTurn         bool

	priceMarkets      [numPlayers]PriceMarkets
	oneFreeResMarkets [numPlayers]OneFreeResMarkets
	endEffects        [numPlayers][]Finaler

	military Board

	vps [numPlayers][numVPTypes]VP

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
	// g.PtokensState.setOnBoard(5, g.rnd)
	// g.PtokensState.setForChoose(3, g.rnd)

	// TODO: replace later, it was written by reason of random
	pts := make([]int, len(g.PtokensState.States))
	for i := range pts {
		pts[i] = i
	}
	g.rnd.Shuffle(len(pts), func(i, j int) { pts[i], pts[j] = pts[j], pts[i] })
	for i := 0; i < 5; i++ {
		g.PtokensState.States[pts[i]].PTokenState = PTokenOnBoard
	}
	for i := 5; i < 8; i++ {
		g.PtokensState.States[pts[i]].PTokenState = PTokenChosenFromDiscarded
	}

	for _, id := range TakeNfromM(InitialWonders, len(g.WondersState.States), g.rnd) {
		g.WondersState.States[id].InGame = true
	}

	g.ageDesk = newAgeDesk(structureAgeI, shuffleAgeI(g.rnd))
	for _, cid := range g.ageDesk.cards {
		g.CardsState.Cards[cid].CardStateEnum = CardOnBoard
	}
	g.ageDesk2 = newAgeDesk(structureAgeII, shuffleAgeII(g.rnd))
	g.ageDesk3 = newAgeDesk(structureAgeIII, shuffleAgeIII(g.rnd))

	for i := 0; i < numPlayers; i++ {
		g.players[i] = Player{}
		g.players[i].Coins = 7
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

func (g *Game) GetAvailablePTokens() (ptokens []PTokenID) {
	return g.PtokensState.GetByState(PTokenOnBoard)
}

// SelectWonders as part of an initialize of a game
func (g *Game) SelectWonders(fstWonders, sndWonders [numWondersPerPlayer]WonderID) error {
	if !g.state.Is(StateSelectWonders) {
		return ErrWrongState
	}

	for _, fw := range fstWonders {
		err := g.GameState.WondersState.chooseByPlayer(fw, 0)
		if err != nil {
			return fmt.Errorf("choose wonder by player id = 0: %w", err)
		}
	}
	for _, sw := range sndWonders {
		err := g.GameState.WondersState.chooseByPlayer(sw, 1)
		if err != nil {
			return fmt.Errorf("choose wonder by player id = 1: %w", err)
		}
	}

	g.state = g.state.Next()
	return nil
}

func (g *Game) DeskCardsState() CardsState {
	return g.ageDesk.state
}

func (g *Game) CardCostCoins(id CardID) Coins {
	c, t := g.CardCost(id)
	return c + t
}

func (g *Game) CardCost(id CardID) (coins Coins, tradeRes Coins) {
	card := id.card()

	if card.FreeChain != nil && g.currentPlayer().Chains.Contain(Chain(*card.FreeChain)) {
		return 0, 0
	}

	reduces := g.oneFreeResMarkets[g.currentPlayerIndex]
	if g.currentPlayer().IsMasonry && card.Color == Blue {
		reduces = append(reduces, OneAnyMarket(), OneAnyMarket())
	}
	return CostTrade(
		card.Cost,
		*g.currentPlayer(),
		g.newTradingPrice(),
		reduces,
	)
}

func (g *Game) WonderCost(id WonderID) (coins Coins, tradeRes Coins) {
	reduces := g.oneFreeResMarkets[g.currentPlayerIndex]
	if g.currentPlayer().IsArchitecture {
		reduces = append(reduces, OneAnyMarket(), OneAnyMarket())
	}
	return CostTrade(
		id.wonder().Cost,
		*g.currentPlayer(),
		g.newTradingPrice(),
		reduces,
	)
}

func (g *Game) newTradingPrice() TradingPrice {
	return NewTradingPrice(
		*g.opponent(),
		g.priceMarkets[g.currentPlayerIndex]...,
	)
}

func (g *Game) ConstructBuilding(id CardID) (state CardsState, _ error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

	err := g.ageDesk.testBuild(id)
	if err != nil {
		return state, fmt.Errorf("card (id = %d) cannot be built: %w", id, err)
	}

	coins, trade := g.CardCost(id)
	pay := coins + trade
	if g.currentPlayer().Coins < pay {
		return state, fmt.Errorf("not enough coins")
	}

	err = g.ageDesk.Build(id)
	state = g.ageDesk.state
	if err != nil {
		return state, err
	}
	g.currentPlayer().Coins -= pay
	if g.opponent().IsEconomy {
		g.opponent().Coins += trade
	}
	if g.currentPlayer().IsStrategy && id.card().Color == Red {
		Shields(1).applyEffect(g, g.currentPlayerIndex)
	}
	card := id.card()
	if g.currentPlayer().IsUrbanism && card.FreeChain != nil && g.currentPlayer().Chains.Contain(Chain(*card.FreeChain)) {
		g.currentPlayer().Coins += 4
	}

	err = g.buildCard(id)
	if err != nil {
		return state, fmt.Errorf("build card: %w", err)
	}
	g.nextTurn()

	return state, nil
}

func (g *Game) buildCard(id CardID) error {
	card := id.card()
	for _, eff := range card.Effects {
		eff.applyEffect(g, g.currentPlayerIndex)
	}
	return g.CardsState.built(id, g.currentPlayerIndex)
}

func (g *Game) DiscardCard(id CardID) (state CardsState, _ error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

	err := g.ageDesk.Build(id)
	state = g.ageDesk.state
	if err != nil {
		return state, fmt.Errorf("build on desk: %w", err)
	}

	err = g.CardsState.discard(id)
	if err != nil {
		return state, fmt.Errorf("discard card: %w", err)
	}

	g.currentPlayer().Coins += Coins(2) + Coins(g.CardsState.NumByColor(Yellow, g.currentPlayerIndex))

	g.nextTurn()
	return state, nil
}

func (g *Game) ConstructWonder(cid CardID, wid WonderID) (state CardsState, _ error) {
	state = g.ageDesk.state
	if !g.state.Is(StateGameTurn) {
		return state, ErrWrongState
	}

	err := g.GameState.WondersState.IsBuildable(wid, g.currentPlayerIndex)
	if err != nil {
		return state, fmt.Errorf("wonder (id = %d) cannot be built by %d player: %w", wid, g.currentPlayerIndex, err)
	}

	err = g.ageDesk.testBuild(cid)
	state = g.ageDesk.state
	if err != nil {
		return state, fmt.Errorf("card (id = %d) cannot be taken: %w", cid, err)
	}

	coins, trade := g.WonderCost(wid)
	pay := coins + trade
	if g.currentPlayer().Coins < pay {
		return state, fmt.Errorf("not enough coins")
	}

	err = g.ageDesk.Build(cid)
	state = g.ageDesk.state
	if err != nil {
		return state, err
	}
	g.currentPlayer().Coins -= pay
	if g.opponent().IsEconomy {
		g.opponent().Coins += trade
	}
	if g.currentPlayer().IsTheology {
		g.repeatTurn = true
	}

	wonder := wid.wonder()
	wonder.Effect.applyEffect(g, g.currentPlayerIndex)

	g.GameState.WondersState.built(wid, g.currentPlayerIndex)

	g.nextTurn()

	return state, nil
}

func (g *Game) ChoosePToken(id PTokenID) error {
	if !g.state.Is(StateChoosePToken) {
		return ErrWrongState
	}

	err := g.PtokensState.take(id, g.currentPlayerIndex)
	if err != nil {
		return fmt.Errorf("PToken (id=%d) cannot be taken: %w", id, err)
	}

	id.pToken().Effect.applyEffect(g, g.currentPlayerIndex)

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) DiscardedCards() []CardID {
	return g.CardsState.Get(CardDiscarded)
}

func (g *Game) ConstructDiscardedCard(id CardID) error {
	if !g.state.Is(StateBuildFreeDiscarded) {
		return ErrWrongState
	}

	err := g.CardsState.builtDiscarded(id, g.currentPlayerIndex)
	if err != nil {
		return fmt.Errorf("cannot construct discarded card (ID: %d): %w", id, err)
	}

	g.buildCard(id)
	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) GetDiscardedOpponentsBuildings() ([]CardID, error) {
	if !g.state.Is(StateDiscardOpponentBuild) {
		return nil, ErrWrongState
	}
	return g.CardsState.ByColor(g.discardOpponentBuild.color, g.CurrentPlayerIndex().Next()), nil
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

	err := g.CardsState.discardFromPlayer(id, opponentID)
	if err != nil {
		return fmt.Errorf("discard card from %d player: %w", opponentID, err)
	}

	card.discard(g, opponentID)

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) GetDiscardedPTokens() (_ []PTokenID, err error) {
	if !g.state.Is(StateBuildFreePToken) {
		return nil, ErrWrongState
	}

	return g.PtokensState.GetByState(PTokenChosenFromDiscarded), nil
}

func (g *Game) PlayDiscardedPToken(pid PTokenID) error {
	if !g.state.Is(StateBuildFreePToken) {
		return ErrWrongState
	}

	err := g.PtokensState.takeFromChosen(pid, g.currentPlayerIndex)
	if err != nil {
		return fmt.Errorf("take ptoken from chosen: %w", err)
	}

	pid.pToken().Effect.applyEffect(g, g.currentPlayerIndex)

	g.state = g.state.Next()
	g.nextTurn()
	return nil
}

func (g *Game) ChooseFirstPlayer(i PlayerIndex) error {
	if !g.state.Is(StateChooseFirstPlayer) {
		return ErrWrongState
	}

	g.currentPlayerIndex = i % numPlayers

	g.state = g.state.Next()
	return nil
}

func (g *Game) nextTurn() {
	if !g.state.Is(StateGameTurn) {
		if g.state.Is(StateVictory) {
			g.finalVPs()
		}
		return
	}
	if g.ageDesk.state.anyExists() {
		if g.repeatTurn {
			g.repeatTurn = false
		} else {
			g.currentPlayerIndex = g.currentPlayerIndex.Next()
		}
	} else {
		if g.CurrentAge == AgeIII {
			g.victory(WinnerBoth, CivilianVictory)
			g.finalVPs()
			return
		}
		g.repeatTurn = false
		g.nextAge()
	}
}

func (g *Game) getWinner() Winner {
	switch {
	case g.vps[0][SumVP] > g.vps[1][SumVP]:
		return Winner1Player
	case g.vps[1][SumVP] > g.vps[0][SumVP]:
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

func (g *Game) nextAge() {
	var pos = g.military.ConflictPawn.Position()

	switch {
	case pos > 0:
		g.currentPlayerIndex = 1
		g.state = StateChooseFirstPlayer
	case pos < 0:
		g.currentPlayerIndex = 0
		g.state = StateChooseFirstPlayer
	default:
		// same player starts next age
	}

	g.CurrentAge = g.CurrentAge.Next()
	switch g.CurrentAge {
	case AgeII:
		g.ageDesk = g.ageDesk2
		for _, cid := range g.ageDesk.cards {
			g.CardsState.Cards[cid].CardStateEnum = CardOnBoard
		}
	case AgeIII:
		g.ageDesk = g.ageDesk3
		for _, cid := range g.ageDesk.cards {
			g.CardsState.Cards[cid].CardStateEnum = CardOnBoard
		}
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

func (g *Game) finalVPs() {
	for pi, fs := range g.endEffects {
		for _, f := range fs {
			vp := f.finalVP(g, PlayerIndex(pi))
			g.vps[pi][vp.t] += vp.v
		}
	}

	for i := PlayerIndex(0); i < numPlayers; i++ {
		g.vps[i][CoinsVP] = VP(1).Mul(g.player(i).Coins.Div(3))
	}

	for i := PlayerIndex(0); i < numPlayers; i++ {
		g.vps[i][MilitaryVP] = g.military.VP(i)
	}

	for pi, vs := range g.vps {
		g.vps[pi][SumVP] = 0
		for _, v := range vs {
			g.vps[pi][SumVP] += v
		}
	}

	if g.victoryType == CivilianVictory {
		g.winner = g.getWinner()
	}
}

func (g *Game) VictoryResult() (w Winner, vic VictoryType, vps [2][numVPTypes]VP, _ error) {
	if !g.state.Is(StateVictory) {
		return w, vic, vps, ErrWrongState
	}
	return g.winner, g.victoryType, g.vps, nil
}

func (g *Game) victory(w Winner, vic VictoryType) {
	if g.state.Is(StateVictory) {
		return
	}
	g.state = StateVictory
	g.winner = w
	g.victoryType = vic
}

// VictoryType - one of the way to claim victory.
type VictoryType uint8

// In 7 Wonders Duel, there are 3 ways to claim victory: military supremacy, scientific supremacy, and civilian victory.
const (
	VictoryNone VictoryType = iota
	CivilianVictory
	MilitarySupremacy
	ScientificSupremacy
	VictoryTypeSize = iota

	numVictoryTypes = VictoryTypeSize
)

func (v VictoryType) String() string {
	switch v {
	case CivilianVictory:
		return "CivilianVictory"
	case MilitarySupremacy:
		return "MilitarySupremacy"
	case ScientificSupremacy:
		return "ScientificSupremacy"
	default:
		return "unknown victory type"
	}
}

type Winner uint8

const (
	Winner1Player Winner = iota
	Winner2Player
	WinnerBoth
	numWinners
)

func (w Winner) String() string {
	switch w {
	case Winner1Player:
		return "1 player is winner"
	case Winner2Player:
		return "2 player is winner"
	case WinnerBoth:
		return "No one is winner"
	default:
		return fmt.Sprintf("wrong winner value: %#v", w)
	}
}

var _ = [1]struct{}{}[numWinners-1-numPlayers]

func (g *Game) gettingPToken(_ PlayerIndex) {
	g.state = StateChoosePToken
	ps := g.PtokensState.GetByState(PTokenOnBoard)
	if len(ps) == 0 {
		g.state = g.state.Next()
	}
}

func (g *Game) Military() Board {
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
	case StateChooseFirstPlayer:
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
	ds := g.DiscardedCards()
	if len(ds) == 0 {
		g.state = g.state.Next()
	}
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

	ds, _ := g.GetDiscardedOpponentsBuildings()
	if len(ds) == 0 {
		g.state = g.state.Next()
	}
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

	if len(g.PtokensState.GetByState(PTokenChosenFromDiscarded)) == 0 {
		g.state = g.state.Next()
	}
}
