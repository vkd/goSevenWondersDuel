package game

import (
	"math"
)

func trade(tc TradingCosts, cost CostOfCard, player *Player) (out []MaybeRes, money Money) {
	debug("Want to buy: %v (prices: %v)", cost, tc)
	needToBuy := cost.Resources.Sub(player.Resources)
	out, money = minimizeByMarket(tc, needToBuy, player.OneOfAnyMarkets)
	return out, money + cost.Money
}

func minimizeByMarket(tc TradingCosts, res Resources, markets []OneOfAnyMarket) ([]MaybeRes, Money) {
	if len(markets) == 0 {
		return nil, res.Money(tc)
	}

	debug("Mminimize by markets: %v x %v", res, tc)

	var out, tOut []MaybeRes
	var minMoney Money = math.MaxInt32
	var tMoney Money
	var outRes MaybeRes

	for _, r := range markets[0] {
		if res.IsZero(r) {
			continue
		}
		tOut, tMoney = minimizeByMarket(tc, res.TakeOne(r), markets[1:])
		debug("Try to minimize by markets (%v): %v, %v", res.TakeOne(r), tOut, tMoney)
		if tMoney < minMoney {
			minMoney = tMoney
			out = tOut
			outRes.Set(r)
		}
	}

	return append([]MaybeRes{outRes}, out...), minMoney
}
