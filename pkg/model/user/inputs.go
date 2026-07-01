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
