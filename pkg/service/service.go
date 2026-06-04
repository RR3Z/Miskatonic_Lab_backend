package service

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
)

type Service struct {
	User      IUser
	Character characterServices.ICharacter
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:      NewUserService(repos),
		Character: characterServices.NewCharacterService(repos),
	}
}
