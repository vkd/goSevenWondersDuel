package core

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func tTest_simpleBot(t *testing.T) {
	type k struct {
		Winner
		VictoryType
	}
	type v struct {
		Count int
		Sum   int
	}
	var cardsRating = make(map[CardID]int)

	f, err := os.Open("stats")
	require.NoError(t, err)
	defer f.Close()

	rating, err := LoadBotRating(f)
	require.NoError(t, err)

	res := map[k]v{}
	var count = 10_000
	for i := 0; i < count; i++ { // 100_000 - 10s
		// b1 := SimpleBot(nil)
		// b2 := SimpleBot(nil)
		b1 := ratingBot{rand.New(rand.NewSource(time.Now().UnixNano())), rating}
		b2 := ratingBot{rand.New(rand.NewSource(time.Now().UnixNano())), rating}
		g, err := NewGame()
		require.NoError(t, err)
		ws, err := g.WondersState.AvailableToChoose()
		require.NoError(t, err)
		var w1, w2 [4]WonderID
		copy(w1[:], ws[:4])
		copy(w2[:], ws[4:])
		err = g.SelectWonders(w1, w2)
		require.NoError(t, err)

		for g.GetState() != StateVictory {
			switch g.CurrentPlayerIndex() {
			case 0:
				b1.NextTurn(g, 0)
			case 1:
				b2.NextTurn(g, 1)
			}
		}
		w, vt, vps, err := g.VictoryResult()
		require.NoError(t, err)
		// if vt == CivilianVictory {
		// 	fmt.Printf("%v by %v - %v:%v\n", w.String(), vt.String(), vps[0][SumVP], vps[1][SumVP])
		// } else {
		// 	fmt.Printf("%v by %v\n", w.String(), vt.String())
		// }
		val := res[k{w, vt}]
		val.Count++
		switch w {
		case Winner1Player:
			val.Sum += int(vps[0][SumVP])
		case Winner2Player:
			val.Sum += int(vps[1][SumVP])
		case WinnerBoth:
			val.Sum += int(vps[0][SumVP])
		}
		res[k{w, vt}] = val

		switch w {
		case Winner1Player:
			for _, cs := range g.BuildCards()[0] {
				for _, cid := range cs {
					cardsRating[cid]++
				}
			}
			for _, cs := range g.BuildCards()[1] {
				for _, cid := range cs {
					cardsRating[cid]--
				}
			}
		case Winner2Player:
			for _, cs := range g.BuildCards()[1] {
				for _, cid := range cs {
					cardsRating[cid]++
				}
			}
			for _, cs := range g.BuildCards()[0] {
				for _, cid := range cs {
					cardsRating[cid]--
				}
			}
		}
	}
	for w := Winner(0); w < WinnerBoth; w++ {
		for vt := VictoryType(1); vt < numVictoryTypes; vt++ {
			v := res[k{w, vt}]
			var rate = 0
			if v.Count > 0 {
				rate = v.Sum / v.Count
			}
			log.Printf("%s by %s:\t%d\t%d%%\t(avg score: %v)", w, vt, v.Count, v.Count*100/count, rate)
		}
	}
	{
		v := res[k{WinnerBoth, CivilianVictory}]
		var rate = 0
		if v.Count > 0 {
			rate = v.Sum / v.Count
		}
		log.Printf("%s by %s:\t%d\t%d%%\t(avg score: %v)", WinnerBoth, CivilianVictory, v.Count, v.Count*100/count, rate)
	}
	// for cid, v := range cardsRating {
	// 	fmt.Printf("%5.0f - %s (%s)\n", float64(v), cid.card().Name, cid.Color())
	// }

	// save rating
	// f, err := os.Create("stats")
	// require.NoError(t, err)
	// defer f.Close()

	// err = saveBotRating(f, cardsRating)
	// assert.NoError(t, er)
}
