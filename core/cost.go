package core

import "fmt"

// Cost of a card or a wonder
type Cost struct {
	Coins     Coins
	Resources Resources
	FreeChain MaybeChain
}

// ByCoins - cost of a card by coins
func (c Cost) ByCoins(g *Game, i PlayerIndex) Coins {
	if c.FreeChain.OK && g.players[i].Chains.Contain(c.FreeChain.Chain) {
		return 0
	}

	missingRes := c.Resources.Reduce(g.players[i].Resources)
	tp := NewTradingPrice(g.players[i.Next()], g.players[i].PriceMarkets)
	return c.Coins + tp.CostOf(missingRes)
}

// NewCost by different goods
func NewCost(args ...interface{}) Cost {
	var c Cost

	for _, arg := range args {
		switch arg := arg.(type) {
		case Coins:
			c.Coins += arg
		case Resource:
			c.Resources[arg]++
		case Resources:
			c.Resources = c.Resources.Add(arg)
		case Chain:
			c.FreeChain.Set(arg)
		default:
			panic(fmt.Sprintf("unknown type %T for cost", arg))
		}
	}
	return c
}

// TradingPrice of resources for one player
type TradingPrice [numResources]Coins

// NewTradingPrice of resources for one player
func NewTradingPrice(opponent Player, markets []PriceMarket) TradingPrice {
	var out TradingPrice
	for i, count := range opponent.Resources {
		out[i] = Coins(2 + count)
	}
	for _, m := range markets {
		out[m.Resource] = m.Price
	}
	return out
}

// CostOf resources
func (tp TradingPrice) CostOf(rs Resources) Coins {
	var out Coins
	for i, price := range tp {
		out += price.Mul(rs[i])
	}
	return out
}

// PriceMarket - a market that allows you to buy a resource at one price, independent of the opponent’s resources
type PriceMarket struct {
	Resource Resource
	Price    Coins
}

// Apply ...
func (p PriceMarket) Apply(g *Game, i PlayerIndex) {
	g.players[i].PriceMarkets.Append(p)
}

// PriceMarkets ...
type PriceMarkets []PriceMarket

// Append ...
func (pm *PriceMarkets) Append(p PriceMarket) {
	*pm = append(*pm, p)
}

// OneCoinPrice market
func OneCoinPrice(r Resource) PriceMarket {
	return PriceMarket{
		Resource: r,
		Price:    Coins(1),
	}
}

// // OneAnyMarket - one of these resources by every round
// type OneAnyMarket []Resource

// // OneRawMarket by raw materials
// func OneRawMarket() OneAnyMarket { return OneAnyMarket(rawMaterials) }

// // OneManufacturedMarket by manufactured goods
// func OneManufacturedMarket() OneAnyMarket { return OneAnyMarket(manufacturedGoods) }

// CostOfCard by coins
// func CostOfCard(c *Card, tp TradingPrice, p Player) Coins {
// 	for _, cnd := range c.FreeConstructionConditions {
// 		if cnd.IsFree(p) {
// 			return 0
// 		}
// 	}
// 	return c.Cost.ReduceBy(p.Resources).Cost(tp)
// }

// AnyOneOfCosts with minimum price
// type AnyOneOfCosts []Cost

// // ByCoins - return minimum price by coins
// func (c AnyOneOfCosts) ByCoins(g *Game, i PlayerIndex) (Coins, bool) {
// 	var outCoins Coins
// 	var outOk bool
// 	for _, cost := range c {
// 		cn, ok := cost.ByCoins(g, i)
// 		if !ok {
// 			continue
// 		}
// 		// first time OR less price
// 		if !outOk || cn < outCoins {
// 			outCoins = cn
// 		}
// 		outOk = true
// 	}
// 	return outCoins, outOk
// }
