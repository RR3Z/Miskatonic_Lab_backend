package characterHelpers

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"

func DerivedStatsRecalculationReadiness(characteristics db.Characteristic) (string, bool) {
	if characteristics.Strength == nil || characteristics.Size == nil || characteristics.Dexterity == nil {
		return "required_characteristics_missing", false
	}

	return "", true
}
