package healthDTO

type HealthRequest struct {
	MaxHp       *int16 `json:"max_hp"`
	CurrentHp   *int16 `json:"current_hp"`
	MajorWound  *bool  `json:"major_wound"`
	Unconscious *bool  `json:"unconscious"`
	Dying       *bool  `json:"dying"`
	Dead        *bool  `json:"dead"`
}
