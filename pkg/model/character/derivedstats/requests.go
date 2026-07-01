package derivedStatsDTO

type DerivedStatsRequest struct {
	Speed       *int16  `json:"speed"`
	Physique    *int16  `json:"physique"`
	DamageBonus *string `json:"damage_bonus"`
	DodgeValue  *int16  `json:"dodge_value"`
}
