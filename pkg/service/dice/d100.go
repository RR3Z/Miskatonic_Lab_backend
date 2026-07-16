package dice

import "fmt"

type D100Roll struct {
	Units      int
	Tens       []int
	Candidates []int
}

func RollD100(tensCount int) (D100Roll, error) {
	if tensCount <= 0 {
		return D100Roll{}, fmt.Errorf("d100 needs at least one tens die")
	}

	unitsFace, err := GenerateDiceValue(10)
	if err != nil {
		return D100Roll{}, err
	}

	tens := make([]int, tensCount)
	for i := range tens {
		face, err := GenerateDiceValue(10)
		if err != nil {
			return D100Roll{}, err
		}
		tens[i] = face % 10
	}

	units := unitsFace % 10
	candidates, err := BuildD100Candidates(units, tens)
	if err != nil {
		return D100Roll{}, err
	}

	return D100Roll{Units: units, Tens: tens, Candidates: candidates}, nil
}

func BuildD100Candidates(units int, tens []int) ([]int, error) {
	if units < 0 || units > 9 {
		return nil, fmt.Errorf("d100 units must be between 0 and 9")
	}
	if len(tens) == 0 {
		return nil, fmt.Errorf("d100 needs at least one tens die")
	}

	candidates := make([]int, len(tens))
	for i, tensValue := range tens {
		if tensValue < 0 || tensValue > 9 {
			return nil, fmt.Errorf("d100 tens must be between 0 and 9")
		}

		candidate := tensValue*10 + units
		if candidate == 0 {
			candidate = 100
		}
		candidates[i] = candidate
	}

	return candidates, nil
}
