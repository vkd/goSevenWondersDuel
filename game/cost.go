package game

import (
	"fmt"
	"math"
)

// Cost of construction
type Cost struct {
	Resources Resources
	Money     Money
}

// NewCost - create new cost
func NewCost(args ...interface{}) (out Cost) {
	for _, a := range args {
		switch a := a.(type) {
		case Money:
			out.Money += a
		case Resources:
			out.Resources = out.Resources.Add(a)
		case Resource:
			out.Resources = out.Resources.Change(a, 1)
		// case ChainSymbol:
		default:
			panic(fmt.Sprintf("Not implemented: %T", a))
		}
	}
	return
}

// CostByMoney - calculate cost for current player
func CostByMoney(p *Player, c Cost, tc TradingCosts) Money {
	needToBuyRess := c.Resources.Sub(p.Resources)
	money := minimizeByMarket(tc, needToBuyRess, p.OneOfAnyMarkets)
	return money + c.Money
}

// TradingCosts - costs all resources on market for one player
type TradingCosts [numResources]Money

// NewTradingCosts by opponent player's resources
func NewTradingCosts(rs Resources) TradingCosts {
	var tc TradingCosts
	for _, r := range allResources {
		tc[r] = Money(2 + rs[r])
	}
	return tc
}

// NewTradingCostsDefault - all resources cost by 2
func NewTradingCostsDefault() TradingCosts {
	return NewTradingCosts(Resources{})
}

// Buy resources
func (tc TradingCosts) Buy(rs Resources) Money {
	var res Money
	for i := range rs {
		res += tc[i] * Money(rs[i])
	}
	return res
}

// ApplyMarkets current player's markets
func (tc TradingCosts) ApplyMarkets(mm []OnePriceMarket) TradingCosts {
	for _, m := range mm {
		if m.Price < tc[m.Res] {
			tc[m.Res] = m.Price
		}
	}
	return tc
}

func minimizeByMarket(tc TradingCosts, res Resources, markets []OneOfAnyMarket) Money {
	if len(markets) == 0 {
		return tc.Buy(res)
	}

	debug("Mminimize by markets: %v x %v", res, tc)

	// var out, tOut []MaybeRes
	var minMoney Money = math.MaxInt32
	var tMoney Money
	// var outRes MaybeRes

	for _, r := range markets[0] {
		if res.IsZero(r) {
			continue
		}
		tMoney = minimizeByMarket(tc, res.TakeOne(r), markets[1:])
		// debug("Try to minimize by markets (%v): %v, %v", res.TakeOne(r), tOut, tMoney)
		if tMoney < minMoney {
			minMoney = tMoney
			// out = tOut
			// outRes.Set(r)
		}
	}

	return minMoney
	// return append([]MaybeRes{outRes}, out...), minMoney
}
