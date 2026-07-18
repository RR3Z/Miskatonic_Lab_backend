package inventoryDTO

type InventoryItemRequest struct {
	Name        string  `json:"name"`
	Quantity    *int32  `json:"quantity"`
	Category    *string `json:"category"`
	Description *string `json:"description"`
}
