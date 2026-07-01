package user

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
)

type IUser interface {
	UpsertUser(ctx context.Context, input model.UpsertUserInput) error
	DeleteUser(ctx context.Context, input model.DeleteUserInput) error
	GetUserByID(ctx context.Context, input model.GetUserInput) (model.UserModel, error)
}
