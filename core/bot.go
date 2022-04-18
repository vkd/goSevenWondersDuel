package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strconv"
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
		state := g.DeskCardsState()
		ac := GetUncovered(state, CoverageByForAge(g.CurrentAge))

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

		listWs := g.WondersState.AvailableToBuild(myIdx)
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
			if err == nil {
				break
			}
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

func RatingBot(rating map[CardID]int) Bot {
	return ratingBot{rand.New(rand.NewSource(time.Now().UnixNano())), rating}
}

type ratingBot struct {
	rnd    *rand.Rand
	rating map[CardID]int
}

var _ Bot = ratingBot{}

func (r ratingBot) NextTurn(g *Game, myIdx PlayerIndex) {
	var err error
	switch g.GetState() {
	case StateGameTurn:
		p := g.Player(myIdx)

		ac := GetUncovered(g.DeskCardsState(), CoverageByForAge(g.CurrentAge))

		if len(ac) == 0 {
			panic("This situation is unreal: there is no available cards")
		}

		r.rnd.Shuffle(len(ac), func(i, j int) {
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

		var maxCardID = ac[0]
		var maxRating = r.rating[maxCardID]

		var found bool
		var maxAvailableRating int
		var maxAvailableCardID CardID
		for _, cid := range ac {
			rt := r.rating[cid]
			if rt > maxRating {
				maxRating = rt
				maxCardID = cid
			}

			cost := g.CardCostCoins(cid)
			if cost > p.Coins {
				continue
			}
			if !found || rt > maxAvailableRating {
				found = true
				maxAvailableRating = rt
				maxAvailableCardID = cid
			}
		}

		listWs := g.WondersState.AvailableToBuild(myIdx)
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
			idx := r.rnd.Intn(len(aws))
			_, err = g.ConstructWonder(maxCardID, aws[idx])
			if err == nil {
				break
			}
		}

		if found {
			_, err = g.ConstructBuilding(maxAvailableCardID)
			break
		}

		_, err = g.DiscardCard(maxCardID)
	default:
		simpleBot{r.rnd}.NextTurn(g, myIdx)
		return
	}
	if err != nil {
		panic(fmt.Sprintf("error on bot's next turn: %v", err))
	}
}

func LoadBotRating(r io.Reader) (map[CardID]int, error) {
	out := make(map[CardID]int)
	var rd = csv.NewReader(r)
	for {
		row, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(row) < 2 {
			continue
		}
		cid, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}
		v, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}
		out[CardID(cid)] = v
	}
	return out, nil
}

func SaveBotRating(w io.Writer, m map[CardID]int) error {
	wr := csv.NewWriter(w)
	defer wr.Flush()

	for cid, v := range m {
		err := wr.Write([]string{strconv.Itoa(int(cid)), strconv.Itoa(v)})
		if err != nil {
			return err
		}
	}
	return nil
}
