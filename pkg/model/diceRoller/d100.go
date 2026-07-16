package diceRollerDTO

type D100Mode string

const (
	D100ModeNormal  D100Mode = "normal"
	D100ModeBonus   D100Mode = "bonus"
	D100ModePenalty D100Mode = "penalty"
)

func (m D100Mode) IsValid() bool {
	return m == D100ModeNormal || m == D100ModeBonus || m == D100ModePenalty
}
