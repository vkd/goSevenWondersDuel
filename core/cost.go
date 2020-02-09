package core

type Cost interface {
	cost(*Game, PlayerIndex) (Coins, bool)
}

type Price struct {
	Coins     Coins
	Resources Resources
}

type Pricer interface {
	applyPrice(p *Price)
}

func NewCost(ps ...Pricer) Cost {
	var price Price
	for _, p := range ps {
		p.applyPrice(&price)
	}
	return price
}

func (p Price) cost(g *Game, i PlayerIndex) (Coins, bool) {
	var player = g.player(i)
	var toBuy = p.Resources.Reduce(player.Resources)
	var tradingPrice = NewTradingPrice(*i.Next().player(g), player.PriceMarkets)
	var cost = tradingPrice.CostOf(toBuy) + p.Coins
	return cost, true
}

type orCost []Cost

func (cs orCost) cost(g *Game, i PlayerIndex) (Coins, bool) {
	var min Coins
	var foundOne bool
	for _, c := range cs {
		coins, ok := c.cost(g, i)
		if !ok {
			continue
		}
		if !foundOne || coins < min {
			min = coins
		}
		foundOne = true
	}
	if !foundOne {
		return 0, false
	}
	return min, true
}

type FreeChain Chain

func (c FreeChain) cost(g *Game, i PlayerIndex) (Coins, bool) {
	if g.player(i).Chains.Contain(Chain(c)) {
		return 0, true
	}
	return 0, false
}

// ------------------------------

// TradingPrice of resources for one player
type TradingPrice [numResources]Coins

// NewTradingPrice of resources for one player
func NewTradingPrice(opponent Player, markets []PriceMarket) TradingPrice {
	var out TradingPrice
	for i, count := range opponent.Resources {
		out[i] = Coins(2 + count)
	}
	for _, m := range markets {
		if m.Price < out[m.Resource] {
			out[m.Resource] = m.Price
		}
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
