package service

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type IUser interface {
	UpsertUser(ctx context.Context, input db.UpsertUserParams) error
	DeleteUser(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (db.User, error)
}

type UserService struct {
	repos *repository.Repository
}

func NewUserService(repos *repository.Repository) *UserService {
	return &UserService{repos: repos}
}

func (s *UserService) UpsertUser(ctx context.Context, input db.UpsertUserParams) error {
	_, err := s.repos.Queries.UpsertUser(ctx, db.UpsertUserParams{
		ID:        input.ID,
		Username:  input.Username,
		Email:     input.Email,
		AvatarUrl: input.AvatarUrl,
	})

	return err
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.repos.Queries.DeleteUserByClerkID(ctx, userID)
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (db.User, error) {
	return s.repos.Queries.GetUserByClerkID(ctx, userID)
}
