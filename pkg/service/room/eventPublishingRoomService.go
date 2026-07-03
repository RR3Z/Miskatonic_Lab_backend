package room

import (
	"context"
	"log/slog"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventPublishingRoomService struct {
	next        IRoom
	maintenance IRoomMaintenance
	publisher   events.EventPublisher
}

func NewEventPublishingRoomService(next IRoom, maintenance IRoomMaintenance, publisher events.EventPublisher) *EventPublishingRoomService {
	return &EventPublishingRoomService{
		next:        next,
		maintenance: maintenance,
		publisher:   publisher,
	}
}

func (s *EventPublishingRoomService) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomModel, error) {
	room, err := s.next.CreateRoom(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomCreateFailed{OwnerID: input.OwnerID, Err: err})
		return model.RoomModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomCreateSucceeded{RoomID: room.ID.String(), OwnerID: room.OwnerID})
	return room, nil
}

func (s *EventPublishingRoomService) GetRoom(ctx context.Context, input model.GetRoomInput) (model.RoomModel, error) {
	room, err := s.next.GetRoom(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomGetFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return model.RoomModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomGetSucceeded{RoomID: room.ID.String(), UserID: input.UserID})
	return room, nil
}

func (s *EventPublishingRoomService) UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomModel, error) {
	room, err := s.next.UpdateRoom(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomUpdateFailed{RoomID: input.RoomID.String(), OwnerID: input.OwnerID, Err: err})
		return model.RoomModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomUpdateSucceeded{RoomID: room.ID.String(), OwnerID: input.OwnerID})
	return room, nil
}

func (s *EventPublishingRoomService) TransferOwnership(ctx context.Context, input model.TransferOwnershipInput) (model.RoomModel, error) {
	room, err := s.next.TransferOwnership(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomTransferOwnershipFailed{
			RoomID:     input.RoomID.String(),
			OwnerID:    input.OwnerID,
			NewOwnerID: input.NewOwnerID,
			Err:        err,
		})
		return model.RoomModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomTransferOwnershipSucceeded{
		RoomID:     room.ID.String(),
		OwnerID:    input.OwnerID,
		NewOwnerID: input.NewOwnerID,
	})
	return room, nil
}

func (s *EventPublishingRoomService) DeleteRoom(ctx context.Context, input model.DeleteRoomInput) error {
	if err := s.next.DeleteRoom(ctx, input); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomDeleteFailed{RoomID: input.RoomID.String(), OwnerID: input.OwnerID, Err: err})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomDeleteSucceeded{RoomID: input.RoomID.String(), OwnerID: input.OwnerID})
	return nil
}

func (s *EventPublishingRoomService) JoinRoom(ctx context.Context, input model.JoinRoomInput) (model.RoomMemberModel, error) {
	member, err := s.next.JoinRoom(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomMemberJoinFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return model.RoomMemberModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomMemberJoinSucceeded{
		RoomID:   member.RoomID.String(),
		UserID:   member.UserID,
		MemberID: member.ID.String(),
	})
	return member, nil
}

func (s *EventPublishingRoomService) LeaveRoom(ctx context.Context, input model.LeaveRoomInput) (model.LeaveRoomResult, error) {
	result, err := s.next.LeaveRoom(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomMemberLeaveFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return model.LeaveRoomResult{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomMemberLeaveSucceeded{
		RoomID:        input.RoomID.String(),
		UserID:        input.UserID,
		DeletedRoomID: uuidPtrString(result.DeletedRoomID),
	})
	return result, nil
}

func (s *EventPublishingRoomService) KickMember(ctx context.Context, input model.KickMemberInput) error {
	if err := s.next.KickMember(ctx, input); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomMemberKickFailed{
			RoomID:       input.RoomID.String(),
			ActorUserID:  input.ActorUserID,
			TargetUserID: input.TargetUserID,
			Err:          err,
		})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomMemberKickSucceeded{
		RoomID:       input.RoomID.String(),
		ActorUserID:  input.ActorUserID,
		TargetUserID: input.TargetUserID,
	})
	return nil
}

func (s *EventPublishingRoomService) SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMemberModel, error) {
	member, err := s.next.SelectCharacter(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomMemberSelectCharacterFailed{
			RoomID:      input.RoomID.String(),
			UserID:      input.UserID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return model.RoomMemberModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomMemberSelectCharacterSucceeded{
		RoomID:      member.RoomID.String(),
		UserID:      member.UserID,
		CharacterID: member.CharacterID.String(),
	})
	return member, nil
}

func (s *EventPublishingRoomService) ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMemberModel, error) {
	member, err := s.next.ChangeRole(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomMemberChangeRoleFailed{
			RoomID:       input.RoomID.String(),
			ActorUserID:  input.ActorUserID,
			TargetUserID: input.TargetUserID,
			Role:         input.Role,
			Err:          err,
		})
		return model.RoomMemberModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomMemberChangeRoleSucceeded{
		RoomID:       member.RoomID.String(),
		ActorUserID:  input.ActorUserID,
		TargetUserID: member.UserID,
		Role:         member.Role,
	})
	return member, nil
}

func (s *EventPublishingRoomService) ListSelectedCharacters(ctx context.Context, input model.ListSelectedCharactersInput) ([]model.SelectedCharacterModel, error) {
	characters, err := s.next.ListSelectedCharacters(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomSelectedCharactersListFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return nil, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomSelectedCharactersListSucceeded{
		RoomID: input.RoomID.String(),
		UserID: input.UserID,
		Count:  len(characters),
	})
	return characters, nil
}

func (s *EventPublishingRoomService) TouchRoomActivity(ctx context.Context, input model.TouchRoomActivityInput) error {
	if err := s.next.TouchRoomActivity(ctx, input); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomActivityTouchFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomActivityTouchSucceeded{RoomID: input.RoomID.String(), UserID: input.UserID})
	return nil
}

func (s *EventPublishingRoomService) EnsureMember(ctx context.Context, roomID pgtype.UUID, userID string) error {
	if err := s.next.EnsureMember(ctx, roomID, userID); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomEnsureMemberFailed{RoomID: roomID.String(), UserID: userID, Err: err})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomEnsureMemberSucceeded{RoomID: roomID.String(), UserID: userID})
	return nil
}

func (s *EventPublishingRoomService) EnsureOwner(ctx context.Context, roomID pgtype.UUID, userID string) error {
	if err := s.next.EnsureOwner(ctx, roomID, userID); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomEnsureOwnerFailed{RoomID: roomID.String(), UserID: userID, Err: err})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomEnsureOwnerSucceeded{RoomID: roomID.String(), UserID: userID})
	return nil
}

func (s *EventPublishingRoomService) EnsureCanPublishRoomEvent(ctx context.Context, roomID pgtype.UUID, userID string) error {
	if err := s.next.EnsureCanPublishRoomEvent(ctx, roomID, userID); err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomEnsureCanPublishEventFailed{RoomID: roomID.String(), UserID: userID, Err: err})
		return err
	}

	s.publisher.Publish(ctx, roomEvents.RoomEnsureCanPublishEventSucceeded{RoomID: roomID.String(), UserID: userID})
	return nil
}

func (s *EventPublishingRoomService) ListRoomEvents(ctx context.Context, input model.ListRoomEventsInput) ([]model.RoomEventModel, error) {
	roomEventsModels, err := s.next.ListRoomEvents(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomEventsListFailed{RoomID: input.RoomID.String(), UserID: input.UserID, Err: err})
		return nil, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomEventsListSucceeded{
		RoomID: input.RoomID.String(),
		UserID: input.UserID,
		Count:  len(roomEventsModels),
	})
	return roomEventsModels, nil
}

func (s *EventPublishingRoomService) CreateChatMessage(ctx context.Context, input model.CreateChatMessageInput) (model.RoomEventModel, error) {
	event, err := s.next.CreateChatMessage(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomChatMessageCreateFailed{RoomID: input.RoomID.String(), ActorID: input.ActorID, Err: err})
		return model.RoomEventModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomChatMessageCreateSucceeded{
		RoomID:  event.RoomID.String(),
		ActorID: event.ActorID,
		EventID: event.ID.String(),
	})
	return event, nil
}

func (s *EventPublishingRoomService) CreateDiceRollRoomEvent(ctx context.Context, input model.CreateDiceRollRoomEventInput) (model.RoomEventModel, error) {
	event, err := s.next.CreateDiceRollRoomEvent(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomDiceRollEventCreateFailed{
			RoomID:      input.RoomID.String(),
			ActorID:     input.ActorID,
			RollID:      input.RollID,
			CharacterID: input.CharacterID,
			Err:         err,
		})
		return model.RoomEventModel{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomDiceRollEventCreateSucceeded{
		RoomID:      event.RoomID.String(),
		ActorID:     event.ActorID,
		EventID:     event.ID.String(),
		RollID:      input.RollID,
		CharacterID: input.CharacterID,
	})
	return event, nil
}

func (s *EventPublishingRoomService) CreateCharacterChangedRoomEvents(ctx context.Context, input model.CreateCharacterChangedRoomEventsInput) ([]model.RoomEventModel, error) {
	createdEvents, err := s.next.CreateCharacterChangedRoomEvents(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomCharacterChangedEventsCreateFailed{
			ActorID:     input.ActorID,
			CharacterID: input.CharacterID.String(),
			Err:         err,
		})
		return nil, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomCharacterChangedEventsCreateSucceeded{
		ActorID:     input.ActorID,
		CharacterID: input.CharacterID.String(),
		Count:       len(createdEvents),
	})
	return createdEvents, nil
}

func (s *EventPublishingRoomService) CleanupRooms(ctx context.Context, input model.CleanupRoomsInput) (model.CleanupRoomsResult, error) {
	result, err := s.maintenance.CleanupRooms(ctx, input)
	if err != nil {
		s.publisher.Publish(ctx, roomEvents.RoomCleanupFailed{Err: err})
		return model.CleanupRoomsResult{}, err
	}

	s.publisher.Publish(ctx, roomEvents.RoomCleanupSucceeded{
		InactiveDeleted: result.InactiveDeleted,
		InvalidDeleted:  result.InvalidDeleted,
		DeletedCount:    len(result.DeletedRoomIDs),
	})
	return result, nil
}

func (s *EventPublishingRoomService) StartCleanupWorker(ctx context.Context, interval time.Duration, afterCleanup func(model.CleanupRoomsResult)) {
	if interval <= 0 {
		interval = DEFAULT_ROOM_CLEANUP_INTERVAL
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				result, err := s.CleanupRooms(ctx, model.CleanupRoomsInput{})
				if err != nil {
					slog.Warn("room cleanup failed", "component", "room_cleanup", "error", err)
					continue
				}
				if result.InactiveDeleted > 0 || result.InvalidDeleted > 0 {
					slog.Info(
						"room cleanup deleted rooms",
						"component", "room_cleanup",
						"inactive_deleted", result.InactiveDeleted,
						"invalid_deleted", result.InvalidDeleted,
					)
				}
				if afterCleanup != nil {
					afterCleanup(result)
				}
			}
		}
	}()
}

func uuidPtrString(uuid *pgtype.UUID) *string {
	if uuid == nil {
		return nil
	}

	value := uuid.String()
	return &value
}
