package user

import (
	"context"
	"errors"
	"strings"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
	userHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/helpers"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repos *repository.Repository
}

func NewUserService(repos *repository.Repository) *UserService {
	return &UserService{repos: repos}
}

func (s *UserService) UpsertUser(ctx context.Context, input model.UpsertUserInput) error {
	if err := validateUserID(input.ID); err != nil {
		return err
	}

	params := db.UpsertUserParams{
		ID:        input.ID,
		Username:  userHelpers.ResolveUsername(input),
		Email:     userHelpers.ResolveEmail(input),
		AvatarUrl: input.AvatarURL,
	}

	_, err := s.repos.Queries.UpsertUser(ctx, params)
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, input model.DeleteUserInput) error {
	if strings.TrimSpace(input.ID) == "" {
		return userErrors.ErrMissingUserID
	}

	return s.repos.Queries.DeleteUserByClerkID(ctx, input.ID)
}

func (s *UserService) GetUserByID(ctx context.Context, input model.GetUserInput) (model.UserModel, error) {
	user, err := s.repos.Queries.GetUserByClerkID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.UserModel{}, userErrors.ErrUserNotFound
		}
		return model.UserModel{}, err
	}

	return model.ToUserModel(user), nil
}
