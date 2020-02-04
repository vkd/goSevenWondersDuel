package game

// from rulebook: Free Construction
// func TestCost_FreeConstruction(t *testing.T) {
// 	c := CardName("Lumber yard").find().Cost
// 	m := CostByMoney(&Player{}, c, NewTradingCostsDefault())
// 	assert.Equal(t, Money(0), m)
// }

// // from rulebook: Production
// func TestCost_Production(t *testing.T) {
// 	p := &Player{}
// 	p.Resources[Stone] = 1
// 	p.Resources[Clay] = 3
// 	p.Resources[Papyrus] = 1

// 	tc := NewTradingCostsDefault()

// 	c := CardName("Baths").find().Cost
// 	m := CostByMoney(p, c, tc)
// 	assert.True(t, m == 0)

// 	c = CardName("Garrison").find().Cost
// 	assert.True(t, CostByMoney(p, c, tc) == 0)

// 	c = CardName("Apothecary").find().Cost
// 	assert.True(t, CostByMoney(p, c, tc) != 0)
// }

// // from rulebook: Trading
// func TestCost_TradingCosts(t *testing.T) {
// 	var bruno, antoine Player

// 	bruno.Resources[Stone] = 2

// 	// for Antoine
// 	costOneStone := NewTradingCosts(bruno.Resources)[Stone]
// 	assert.Equal(t, Money(4), costOneStone)

// 	// for Bruno
// 	costOneStone = NewTradingCosts(antoine.Resources)[Stone]
// 	assert.Equal(t, Money(2), costOneStone)
// }

// // from rulebook: Trading
// func TestCost_Trading(t *testing.T) {
// 	var bruno, antoine Player

// 	bruno.Resources[Stone] = 2
// 	antoine.Resources[Clay] = 1

// 	c := CardName("Fortifications").find().Cost
// 	tc := NewTradingCosts(antoine.Resources)
// 	m := CostByMoney(&bruno, c, tc)
// 	assert.Equal(t, Money(5), m)

// 	c = CardName("Aqueduct").find().Cost
// 	tc = NewTradingCosts(bruno.Resources)
// 	m = CostByMoney(&antoine, c, tc)
// 	assert.Equal(t, Money(12), m)
// }
