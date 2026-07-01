package userDTO

type EmailInput struct {
	ID           string
	EmailAddress string
}

type UpsertUserInput struct {
	ID                    string
	Username              *string
	PrimaryEmailAddressID *string
	EmailAddresses        []EmailInput
	AvatarURL             *string
}

type DeleteUserInput struct {
	ID string
}

type GetUserInput struct {
	ID string
}

func ToUpsertUserInput(data ClerkWebhookUserData) UpsertUserInput {
	emails := make([]EmailInput, len(data.EmailAddresses))
	for i, e := range data.EmailAddresses {
		emails[i] = EmailInput{
			ID:           e.ID,
			EmailAddress: e.EmailAddress,
		}
	}

	return UpsertUserInput{
		ID:                    data.ID,
		Username:              data.Username,
		PrimaryEmailAddressID: data.PrimaryEmailAddressID,
		EmailAddresses:        emails,
		AvatarURL:             data.ImageURL,
	}
}
