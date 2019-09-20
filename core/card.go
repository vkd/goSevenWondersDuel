package core

// Card - In 7 Wonders Duel, all of the Age and Guild cards represent Buildings.
// The Building cards all consist of a name, an effect and a construction cost.
type Card struct {
	Name    CardName
	Color   CardColor
	Effects []Effect
	Cost    Cost
}

// CardName - name of card
type CardName string

// CardColor - There are 7 different types of Buildings, easily identifiable by their colored border.
type CardColor uint8

// Color of card
const (
	Brown           CardColor = iota // Raw materials
	Grey                             // Manufactured goods
	Blue                             // Civilian Buildings
	Green                            // Scientific Buildings
	Yellow                           // Commercial Buildings
	Red                              // Military Buildings
	Purple                           // Guilds
	numOfCardColors = iota
)

var (
	nameCardColor = map[CardColor]string{
		Brown:  "Brown",
		Grey:   "Grey",
		Blue:   "Blue",
		Green:  "Green",
		Yellow: "Yellow",
		Red:    "Red",
		Purple: "Purple",
	}
	_ = [1]struct{}{}[len(nameCardColor)-numOfCardColors]
)

// String representation of card color
func (c CardColor) String() string { return nameCardColor[c] }
