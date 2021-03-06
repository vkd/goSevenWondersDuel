package core

import "fmt"

// PlayerIndex - index of player
type PlayerIndex int

// Next player index
func (i PlayerIndex) Next() PlayerIndex {
	return (i + 1) % numPlayers
}

func (i PlayerIndex) winner() Winner {
	switch i {
	case 0:
		return Winner1Player
	case 1:
		return Winner2Player
	default:
		panic(fmt.Sprintf("unknown player index %d", i))
	}
}

// Player of a game
type Player struct {
	PlayerState

	Resources Resources

	Chains Chains

	ScientificSymbols ScientificSymbols

	IsArchitecture bool
	IsEconomy      bool
	IsMasonry      bool
	IsStrategy     bool
	IsTheology     bool
	IsUrbanism     bool
}

type PlayerState struct {
	Coins Coins
}
