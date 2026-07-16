package skillsDTO

type SkillRequest struct {
	Name      string `json:"name"`
	BaseValue int16  `json:"base_value"`
	Value     int16  `json:"value"`
	Checked   bool   `json:"checked"`
}
