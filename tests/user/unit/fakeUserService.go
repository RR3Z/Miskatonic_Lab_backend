package tests

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type FakeUserService struct {
	UpsertUserCalls int
	DeleteUserCalls int
	GetUserCalls    int

	LastUpsertUserInput db.UpsertUserParams
	LastDeleteUserID    string

	UpsertUserErr error
	DeleteUserErr error
	GetUserResult db.User
	GetUserErr    error
}

func (f *FakeUserService) UpsertUser(_ context.Context, input db.UpsertUserParams) error {
	f.UpsertUserCalls++
	f.LastUpsertUserInput = input

	return f.UpsertUserErr
}

func (f *FakeUserService) DeleteUser(_ context.Context, userID string) error {
	f.DeleteUserCalls++
	f.LastDeleteUserID = userID

	return f.DeleteUserErr
}

func (f *FakeUserService) GetUserByID(_ context.Context, _ string) (db.User, error) {
	f.GetUserCalls++

	return f.GetUserResult, f.GetUserErr
}
