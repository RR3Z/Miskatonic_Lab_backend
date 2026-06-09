package calculators

import (
	"fmt"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func CalculateDerivedStats(userID string, characterID pgtype.UUID, age int16, characteristics db.Characteristic) db.UpsertDerivedStatsParams {
	return db.UpsertDerivedStatsParams{
		UserID:      userID,
		CharacterID: characterID,
		Speed:       calculateSpeed(characteristics, age),
		Physique:    calculatePhysique(characteristics),
		DamageBonus: calculateDamageBonus(characteristics),
		DodgeValue:  calculateDodgeValue(characteristics),
	}
}

func calculateSpeed(characteristics db.Characteristic, age int16) *int16 {
	var speed int16

	if *characteristics.Strength > *characteristics.Size && *characteristics.Dexterity > *characteristics.Size {
		speed = 9
	} else if *characteristics.Strength < *characteristics.Size && *characteristics.Dexterity < *characteristics.Size {
		speed = 7
	} else {
		speed = 8
	}

	speed -= ageSpeedPenalty(age)

	return &speed
}

func ageSpeedPenalty(age int16) int16 {
	var penaltyValue int16 = 0

	if age >= 40 && age <= 49 {
		penaltyValue = 1
	} else if age >= 50 && age <= 59 {
		penaltyValue = 2
	} else if age >= 60 && age <= 69 {
		penaltyValue = 3
	} else if age >= 70 && age <= 79 {
		penaltyValue = 4
	} else if age >= 80 {
		penaltyValue = 5
	}

	return penaltyValue
}

func calculatePhysique(characteristics db.Characteristic) *int16 {
	var physiqueValue int16
	characteristicsValue := *characteristics.Strength + *characteristics.Size

	if characteristicsValue >= 2 && characteristicsValue <= 64 {
		physiqueValue = -2
	} else if characteristicsValue >= 65 && characteristicsValue <= 84 {
		physiqueValue = -1
	} else if characteristicsValue >= 85 && characteristicsValue <= 124 {
		physiqueValue = 0
	} else if characteristicsValue >= 125 && characteristicsValue <= 164 {
		physiqueValue = 1
	} else if characteristicsValue >= 165 && characteristicsValue <= 204 {
		physiqueValue = 2
	} else {
		physiqueValue = 3 + (characteristicsValue-205)/80
	}

	return &physiqueValue
}

func calculateDamageBonus(characteristics db.Characteristic) *string {
	characteristicsValue := *characteristics.Strength + *characteristics.Size
	damageBonusValue := "0"

	if characteristicsValue >= 2 && characteristicsValue <= 64 {
		damageBonusValue = "-2"
	} else if characteristicsValue >= 65 && characteristicsValue <= 84 {
		damageBonusValue = "-1"
	} else if characteristicsValue >= 85 && characteristicsValue <= 124 {
		damageBonusValue = "0"
	} else if characteristicsValue >= 125 && characteristicsValue <= 164 {
		damageBonusValue = "+1d4"
	} else if characteristicsValue >= 165 && characteristicsValue <= 204 {
		damageBonusValue = "+1d6"
	} else {
		diceCount := int16(2)
		if characteristicsValue >= 285 {
			diceCount += (characteristicsValue - 205) / 80
		}
		damageBonusValue = fmt.Sprintf("+%dd6", diceCount)
	}

	return &damageBonusValue
}

func calculateDodgeValue(characteristics db.Characteristic) *int16 {
	newVal := *characteristics.Dexterity / 2

	return &newVal
}
