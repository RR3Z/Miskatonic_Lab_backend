package tests

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
)

type FakeUserService struct {
	UpsertUserCalls int
	DeleteUserCalls int
	GetUserCalls    int

	LastUpsertUserInput model.UpsertUserInput
	LastDeleteUserInput model.DeleteUserInput
	LastGetUserInput    model.GetUserInput

	UpsertUserErr error
	DeleteUserErr error
	GetUserResult model.UserModel
	GetUserErr    error
}

func (f *FakeUserService) UpsertUser(_ context.Context, input model.UpsertUserInput) error {
	f.UpsertUserCalls++
	f.LastUpsertUserInput = input

	return f.UpsertUserErr
}

func (f *FakeUserService) DeleteUser(_ context.Context, input model.DeleteUserInput) error {
	f.DeleteUserCalls++
	f.LastDeleteUserInput = input

	return f.DeleteUserErr
}

func (f *FakeUserService) GetUserByID(_ context.Context, input model.GetUserInput) (model.UserModel, error) {
	f.GetUserCalls++
	f.LastGetUserInput = input

	return f.GetUserResult, f.GetUserErr
}
