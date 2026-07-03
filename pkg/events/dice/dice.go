package dice

type DiceRollMakeSucceeded struct {
	UserID      string
	CharacterID string
	RollID      string
	Expression  string
	Result      int32
	Details     []byte
	RoomID      *string
}

type DiceRollMakeFailed struct {
	UserID      string
	CharacterID string
	Err         error
	RoomID      *string
}

func (DiceRollMakeSucceeded) EventName() string {
	return "dice_roll.make_succeeded"
}

func (DiceRollMakeFailed) EventName() string {
	return "dice_roll.make_failed"
}

type DiceRollsListSucceeded struct {
	UserID      string
	CharacterID string
	Count       int
}

type DiceRollsListFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (DiceRollsListSucceeded) EventName() string {
	return "dice_rolls.list_succeeded"
}

func (DiceRollsListFailed) EventName() string {
	return "dice_rolls.list_failed"
}
