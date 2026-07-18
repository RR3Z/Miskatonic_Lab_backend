package character

type CharacterInventoryItemsListSucceeded struct {
	UserID      string
	CharacterID string
	Count       int
}

type CharacterInventoryItemsListFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterInventoryItemsListSucceeded) EventName() string {
	return "character.inventory_items.list_succeeded"
}

func (CharacterInventoryItemsListFailed) EventName() string {
	return "character.inventory_items.list_failed"
}

type CharacterInventoryItemGetSucceeded struct {
	UserID      string
	CharacterID string
	InventoryID string
	Name        string
}

type CharacterInventoryItemGetFailed struct {
	UserID      string
	CharacterID string
	InventoryID string
	Err         error
}

func (CharacterInventoryItemGetSucceeded) EventName() string {
	return "character.inventory_item.get_succeeded"
}

func (CharacterInventoryItemGetFailed) EventName() string {
	return "character.inventory_item.get_failed"
}

type CharacterInventoryItemCreateSucceeded struct {
	UserID      string
	CharacterID string
	InventoryID string
	Name        string
}

type CharacterInventoryItemCreateFailed struct {
	UserID      string
	CharacterID string
	Err         error
}

func (CharacterInventoryItemCreateSucceeded) EventName() string {
	return "character.inventory_item.create_succeeded"
}

func (CharacterInventoryItemCreateFailed) EventName() string {
	return "character.inventory_item.create_failed"
}

type CharacterInventoryItemUpdateSucceeded struct {
	UserID      string
	CharacterID string
	InventoryID string
	Name        string
}

type CharacterInventoryItemUpdateFailed struct {
	UserID      string
	CharacterID string
	InventoryID string
	Err         error
}

func (CharacterInventoryItemUpdateSucceeded) EventName() string {
	return "character.inventory_item.update_succeeded"
}

func (CharacterInventoryItemUpdateFailed) EventName() string {
	return "character.inventory_item.update_failed"
}

type CharacterInventoryItemDeleteSucceeded struct {
	UserID      string
	CharacterID string
	InventoryID string
}

type CharacterInventoryItemDeleteFailed struct {
	UserID      string
	CharacterID string
	InventoryID string
	Err         error
}

func (CharacterInventoryItemDeleteSucceeded) EventName() string {
	return "character.inventory_item.delete_succeeded"
}

func (CharacterInventoryItemDeleteFailed) EventName() string {
	return "character.inventory_item.delete_failed"
}
