package dice

type DiceRollFormulaComponent struct {
	IsDice bool
	Sides  int
	Count  int
}

type DiceDetailType string

const (
	DiceDetailTypeDice     DiceDetailType = "dice"
	DiceDetailTypeModifier DiceDetailType = "modifier"
)

type DiceRollDetail struct {
	Type  DiceDetailType `json:"type"`
	Sides int            `json:"sides,omitempty"`
	Rolls []int          `json:"rolls,omitempty"`
	Value int            `json:"value,omitempty"`
}
