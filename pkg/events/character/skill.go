package character

type CharacterSkillsListSucceeded struct {
	UserID      string
	CharacterID string
	Count       int
}

type CharacterSkillsListFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterSkillsListSucceeded) EventName() string {
	return "character.skills.list_succeeded"
}

func (CharacterSkillsListFailed) EventName() string {
	return "character.skills.list_failed"
}

type CharacterSkillGetSucceeded struct {
	UserID      string
	CharacterID string
	SkillID     string
	Name        string
}

type CharacterSkillGetFailed struct {
	UserID      string
	CharacterID string
	SkillID     string
	Err         error
}

func (CharacterSkillGetSucceeded) EventName() string {
	return "character.skill.get_succeeded"
}

func (CharacterSkillGetFailed) EventName() string {
	return "character.skill.get_failed"
}

type CharacterSkillCreateSucceeded struct {
	UserID      string
	CharacterID string
	SkillID     string
	Name        string
}

type CharacterSkillCreateFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterSkillCreateSucceeded) EventName() string {
	return "character.skill.create_succeeded"
}

func (CharacterSkillCreateFailed) EventName() string {
	return "character.skill.create_failed"
}

type CharacterSkillUpdateSucceeded struct {
	UserID      string
	CharacterID string
	SkillID     string
	Name        string
}

type CharacterSkillUpdateFailed struct {
	UserID      string
	CharacterID string
	SkillID     string
	Err         error
}

func (CharacterSkillUpdateSucceeded) EventName() string {
	return "character.skill.update_succeeded"
}

func (CharacterSkillUpdateFailed) EventName() string {
	return "character.skill.update_failed"
}

type CharacterSkillDeleteSucceeded struct {
	UserID      string
	CharacterID string
	SkillID     string
}

type CharacterSkillDeleteFailed struct {
	UserID      string
	CharacterID string
	SkillID     string
	Err         error
}

func (CharacterSkillDeleteSucceeded) EventName() string {
	return "character.skill.delete_succeeded"
}

func (CharacterSkillDeleteFailed) EventName() string {
	return "character.skill.delete_failed"
}
