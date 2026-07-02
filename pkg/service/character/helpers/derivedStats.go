package characterHelpers

import "github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"

func DerivedStatsRecalculationReadiness(age *int16, characteristics db.Characteristic) (string, bool) {
	if age == nil {
		return "age_missing", false
	}

	if characteristics.Strength == nil || characteristics.Size == nil || characteristics.Dexterity == nil {
		return "required_characteristics_missing", false
	}

	return "", true
}
