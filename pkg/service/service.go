package service

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user"
)

type Service struct {
	User       user.IUser
	Character  character.ICharacter
	DiceRoller diceRoller.IDiceRoller
	Room       room.IRoom
}

func NewService(repos *repository.Repository, publisher events.EventPublisher) *Service {
	characterService := character.NewCharacterService(repos, publisher)
	diceRollerService := diceRoller.NewDiceRollerService(repos)

	return &Service{
		User:       user.NewUserService(repos),
		Character:  character.NewEventPublishingCharacterService(characterService, publisher),
		DiceRoller: diceRoller.NewEventPublishingDiceRollerService(diceRollerService, publisher),
		Room:       room.NewRoomService(repos),
	}
}
