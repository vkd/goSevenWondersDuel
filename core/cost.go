package core

func CostTrade(cost Cost, p Player, tp TradingPrice, rr ...ResourceReducer) (coins Coins, tradeRes Coins) {
	// TODO: make nil interface is valuable
	if cost == nil {
		return 0, 0
	}
	return cost.cost(p, tp, rr...)
}

func CostByCoins(cost Cost, p Player, tp TradingPrice, rr ...ResourceReducer) Coins {
	c, t := CostTrade(cost, p, tp, rr...)
	return c + t
}

type ResourceReducer interface {
	reduceResources(Resources, TradingPrice) Resources
}

type Cost interface {
	cost(Player, TradingPrice, ...ResourceReducer) (coins Coins, tradeRes Coins)
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

type Price struct {
	Coins     Coins
	Resources Resources
}

func (p Price) cost(player Player, tp TradingPrice, rr ...ResourceReducer) (coins Coins, tradeRes Coins) {
	var toBuy = p.Resources.Reduce(player.Resources)
	for _, reducer := range rr {
		toBuy = reducer.reduceResources(toBuy, tp)
	}
	return p.Coins, tp.CostOf(toBuy)
}

// type orCost []Cost

// func (cs orCost) cost(p Player, tp TradingPrice, rr ...ResourceReducer) (coins Coins, tradeRes Coins) {
// 	var min, minTrade Coins = maxCoins, maxCoins
// 	for _, c := range cs {
// 		coins, trade := c.cost(p, tp, rr...)
// 		if coins+trade < min+minTrade {
// 			min = coins
// 			minTrade = trade
// 		}
// 	}
// 	return min, minTrade
// }

type FreeChain Chain

// func (c FreeChain) cost(p Player, _ TradingPrice, _ ...ResourceReducer) (Coins, Coins) {
// 	if p.Chains.Contain(Chain(c)) {
// 		return 0, 0
// 	}
// 	return maxCoins, maxCoins
// }

// TradingPrice of resources for one player
type TradingPrice [numResources]Coins

// NewTradingPrice of resources for one player
func NewTradingPrice(opponent Player, markets ...PriceMarket) TradingPrice {
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

func (tp TradingPrice) CostOne(r Resource) Coins {
	return tp[r]
}

// ------------------------------

// PriceMarket - a market that allows you to buy a resource at one price, independent of the opponent’s resources
type PriceMarket struct {
	Resource Resource
	Price    Coins
}

var _ Effect = PriceMarket{}

func (p PriceMarket) applyEffect(g *Game, i PlayerIndex) {
	g.priceMarkets[i].Append(p)
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

// OneFreeResMarket - one of these resources by every round
type OneFreeResMarket []Resource

var _ Effect = OneFreeResMarket{}

// OneRawMarket by raw materials
func OneRawMarket() OneFreeResMarket { return OneFreeResMarket(rawMaterials) }

// OneManufacturedMarket by manufactured goods
func OneManufacturedMarket() OneFreeResMarket { return OneFreeResMarket(manufacturedGoods) }

func OneAnyMarket() OneFreeResMarket {
	return OneFreeResMarket(allResources)
}

func (m OneFreeResMarket) applyEffect(g *Game, i PlayerIndex) {
	g.oneFreeResMarkets[i] = append(g.oneFreeResMarkets[i], m)
}

type OneFreeResMarkets []OneFreeResMarket

func (ms OneFreeResMarkets) reduceResources(rs Resources, tp TradingPrice) Resources {
	_, out := reduceCosts(ms, rs, tp)
	return out
}

func reduceCosts(ms OneFreeResMarkets, rs Resources, tp TradingPrice) (Coins, Resources) {
	if len(ms) == 0 {
		return 0, rs
	}

	var max Coins
	var maxRs Resources = rs
	for _, r := range ms[0] {
		if rs[r] == 0 {
			continue
		}
		costOne := tp.CostOne(r)
		reduceCoins, reduceRs := reduceCosts(ms[1:], rs.ReduceOne(r), tp)
		if reduceCoins+costOne > max {
			max = reduceCoins + costOne
			maxRs = reduceRs
		}
	}
	return max, maxRs
}
