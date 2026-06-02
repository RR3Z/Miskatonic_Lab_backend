package service

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

type IUser interface {
	UpsertUserFromClerk(ctx context.Context, input db.UpsertUserParams) error
	DeleteUserFromClerk(ctx context.Context, clerkUserID string) error
	GetUserByClerkID(ctx context.Context, clerkUserID string) error
}

type UserService struct {
	repos *repository.Repository
}

func NewUserService(repos *repository.Repository) *UserService {
	return &UserService{repos: repos}
}

func (s *UserService) UpsertUserFromClerk(ctx context.Context, input db.UpsertUserParams) error {
	_, err := s.repos.Queries.UpsertUser(ctx, db.UpsertUserParams{
		ClerkUserID: input.ClerkUserID,
		Username:    input.Username,
		Email:       input.Email,
		AvatarUrl:   input.AvatarUrl,
	})

	return err
}

func (s *UserService) DeleteUserFromClerk(ctx context.Context, clerkUserID string) error {
	return s.repos.Queries.DeleteUserByClerkID(ctx, clerkUserID)
}

func (s *UserService) GetUserByClerkID(ctx context.Context, clerkUserID string) error {
	_, err := s.repos.Queries.GetUserByClerkID(ctx, clerkUserID)

	return err
}
