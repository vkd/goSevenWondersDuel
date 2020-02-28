package core

import (
	"fmt"
	"math/rand"
	"time"
)

type Bot interface {
	NextTurn(g *Game, idx PlayerIndex)
}

func SimpleBot(rnd *rand.Rand) Bot {
	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return simpleBot{rnd}
}

type simpleBot struct {
	rnd *rand.Rand
}

var _ Bot = simpleBot{}

func (s simpleBot) NextTurn(g *Game, myIdx PlayerIndex) {
	p := g.Player(myIdx)

	var err error
	switch g.GetState() {
	case StateChooseFirstPlayer:
		err = g.ChooseFirstPlayer(myIdx)
	case StateGameTurn:
		state := g.CardsState()
		var ac = make([]CardID, 0, 8)
		for _, s := range state {
			if !s.Exists || s.Covered {
				continue
			}
			ac = append(ac, s.ID)
		}

		if len(ac) == 0 {
			panic("This situation is unreal: there is no available cards")
		}

		s.rnd.Shuffle(len(ac), func(i, j int) {
			ac[i], ac[j] = ac[j], ac[i]
		})

		var isBuilt bool
		for _, cid := range ac {
			switch cid.Color() {
			case Brown, Grey:
			default:
				continue
			}

			card := cid.card()

			switch e := card.Effects[0].(type) {
			case Resource:
				if p.Resources[e] != 0 {
					continue
				}
			default:
				panic(fmt.Sprintf("Unknown effect: %T", e))
			}

			cost := g.CardCostCoins(cid)
			if cost > p.Coins {
				continue
			}

			_, err = g.ConstructBuilding(cid)
			isBuilt = true
			break
		}
		if isBuilt {
			break
		}

		listWs := g.GetMyAvailableWonders()[myIdx]
		var aws []WonderID
		for _, wid := range listWs {
			cs, tp := g.WonderCost(wid)
			cost := cs + tp
			if cost > p.Coins {
				continue
			}
			aws = append(aws, wid)
		}
		if len(aws) > 0 {
			idx := s.rnd.Intn(len(aws))
			_, err = g.ConstructWonder(ac[s.rnd.Intn(len(ac))], aws[idx])
			break
		}

		var found bool
		var minPrice = p.Coins + 1
		var minCardID CardID
		for _, cid := range ac {
			cost := g.CardCostCoins(cid)
			if cost > p.Coins {
				continue
			}
			if cost < minPrice {
				found = true
				minPrice = cost
				minCardID = cid
			}

		}
		if found {
			_, err = g.ConstructBuilding(minCardID)
			break
		}

		_, err = g.DiscardCard(ac[0])
	case StateChoosePToken:
		ptokens := g.GetAvailablePTokens()
		idx := s.rnd.Intn(len(ptokens))
		err = g.ChoosePToken(ptokens[idx])
	case StateBuildFreeDiscarded:
		dcs := g.DiscardedCards()
		idx := s.rnd.Intn(len(dcs))
		err = g.ConstructDiscardedCard(dcs[idx])
	case StateDiscardOpponentBuild:
		var ds []CardID
		ds, err = g.GetDiscardedOpponentsBuildings()
		if err != nil {
			break
		}
		idx := s.rnd.Intn(len(ds))
		err = g.DiscardOpponentBuild(ds[idx])
	case StateBuildFreePToken:
		var ps []PTokenID
		ps, err = g.GetDiscardedPTokens()
		if err != nil {
			break
		}
		idx := s.rnd.Intn(len(ps))
		err = g.PlayDiscardedPToken(ps[idx])
	case StateVictory:
		return
	case StateSelectWonders:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("Unknown state: %v", g.GetState().String()))
	}
	if err != nil {
		panic(fmt.Sprintf("error on bot's next turn: %v", err))
	}
}
