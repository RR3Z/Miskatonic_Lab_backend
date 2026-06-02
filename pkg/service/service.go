package service

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
)

type Service struct {
	User      IUser
	Character ICharacter
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:      NewUserService(repos),
		Character: NewCharacterService(repos),
	}
}
