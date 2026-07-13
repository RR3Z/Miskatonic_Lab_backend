package character

type CharacterPortraitReplaceSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterPortraitReplaceFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterPortraitReplaceSucceeded) EventName() string {
	return "character.portrait_replace_succeeded"
}

func (CharacterPortraitReplaceFailed) EventName() string {
	return "character.portrait_replace_failed"
}
