package character

type CharacterBackstoryGetSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterBackstoryGetFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterBackstoryGetSucceeded) EventName() string {
	return "character.backstory.get_succeeded"
}

func (CharacterBackstoryGetFailed) EventName() string {
	return "character.backstory.get_failed"
}

type CharacterBackstoryUpsertSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterBackstoryUpsertFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterBackstoryUpsertSucceeded) EventName() string {
	return "character.backstory.upsert_succeeded"
}

func (CharacterBackstoryUpsertFailed) EventName() string {
	return "character.backstory.upsert_failed"
}

type CharacterBackstoryDeleteSucceeded struct {
	UserID      string
	CharacterID string
}

type CharacterBackstoryDeleteFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterBackstoryDeleteSucceeded) EventName() string {
	return "character.backstory.delete_succeeded"
}

func (CharacterBackstoryDeleteFailed) EventName() string {
	return "character.backstory.delete_failed"
}

type CharacterBackstoryItemsListSucceeded struct {
	UserID      string
	CharacterID string
	Count       int
}

type CharacterBackstoryItemsListFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterBackstoryItemsListSucceeded) EventName() string {
	return "character.backstory_items.list_succeeded"
}

func (CharacterBackstoryItemsListFailed) EventName() string {
	return "character.backstory_items.list_failed"
}

type CharacterBackstoryItemGetSucceeded struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Section         string
	Title           string
}

type CharacterBackstoryItemGetFailed struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Err             error
}

func (CharacterBackstoryItemGetSucceeded) EventName() string {
	return "character.backstory_item.get_succeeded"
}

func (CharacterBackstoryItemGetFailed) EventName() string {
	return "character.backstory_item.get_failed"
}

type CharacterBackstoryItemCreateSucceeded struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Section         string
	Title           string
}

type CharacterBackstoryItemCreateFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterBackstoryItemCreateSucceeded) EventName() string {
	return "character.backstory_item.create_succeeded"
}

func (CharacterBackstoryItemCreateFailed) EventName() string {
	return "character.backstory_item.create_failed"
}

type CharacterBackstoryItemUpdateSucceeded struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Section         string
	Title           string
}

type CharacterBackstoryItemUpdateFailed struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Err             error
}

func (CharacterBackstoryItemUpdateSucceeded) EventName() string {
	return "character.backstory_item.update_succeeded"
}

func (CharacterBackstoryItemUpdateFailed) EventName() string {
	return "character.backstory_item.update_failed"
}

type CharacterBackstoryItemDeleteSucceeded struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
}

type CharacterBackstoryItemDeleteFailed struct {
	UserID          string
	CharacterID     string
	BackstoryItemID string
	Err             error
}

func (CharacterBackstoryItemDeleteSucceeded) EventName() string {
	return "character.backstory_item.delete_succeeded"
}

func (CharacterBackstoryItemDeleteFailed) EventName() string {
	return "character.backstory_item.delete_failed"
}
