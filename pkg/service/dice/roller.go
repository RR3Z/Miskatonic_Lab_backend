package dice

func RollDice(components []DiceRollFormulaComponent) ([]DiceRollDetail, int, error) {
	result := 0
	var details []DiceRollDetail
	for _, c := range components {
		if !c.IsDice {
			details = append(details, DiceRollDetail{
				Type:  DiceDetailTypeModifier,
				Value: c.Count,
			})
			result += c.Count
			continue
		}

		rolls := make([]int, 0, c.Count)
		for i := 0; i < c.Count; i++ {
			roll, err := GenerateDiceValue(c.Sides)
			if err != nil {
				return nil, 0, err
			}
			rolls = append(rolls, roll)
			result += roll
		}
		details = append(details, DiceRollDetail{
			Type:  DiceDetailTypeDice,
			Sides: c.Sides,
			Rolls: rolls,
		})
	}
	return details, result, nil
}
