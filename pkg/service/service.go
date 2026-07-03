package service

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
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

	roomMaintenance room.IRoomMaintenance
}

type BackgroundWorkerHooks struct {
	RoomCleanup func(roomModel.CleanupRoomsResult)
}

func NewService(repos *repository.Repository, publisher events.EventPublisher) *Service {
	characterService := character.NewCharacterService(repos, publisher)
	diceRollerService := diceRoller.NewDiceRollerService(repos)
	roomService := room.NewRoomService(repos)

	return &Service{
		User:            user.NewUserService(repos),
		Character:       character.NewEventPublishingCharacterService(characterService, publisher),
		DiceRoller:      diceRoller.NewEventPublishingDiceRollerService(diceRollerService, publisher),
		Room:            roomService,
		roomMaintenance: roomService,
	}
}

func (s *Service) StartBackgroundWorkers(ctx context.Context, hooks BackgroundWorkerHooks) {
	s.roomMaintenance.StartCleanupWorker(ctx, room.DEFAULT_ROOM_CLEANUP_INTERVAL, hooks.RoomCleanup)
}
