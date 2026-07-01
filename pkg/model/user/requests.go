package userDTO

type ClerkWebhookUserEvent struct {
	Type string               `json:"type"`
	Data ClerkWebhookUserData `json:"data"`
}

type ClerkWebhookUserData struct {
	ID                    string                  `json:"id"`
	Username              *string                 `json:"username"`
	ImageURL              *string                 `json:"image_url"`
	PrimaryEmailAddressID *string                 `json:"primary_email_address_id"`
	EmailAddresses        []ClerkWebhookUserEmail `json:"email_addresses"`
}

type ClerkWebhookUserEmail struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}
