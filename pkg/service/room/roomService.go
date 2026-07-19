package room

import (
	"context"
	"errors"
	"fmt"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
	"github.com/jackc/pgx/v5"
	"strings"
)

type RoomService struct {
	repos *repository.Repository
}

func NewRoomService(repos *repository.Repository) *RoomService {
	return &RoomService{repos: repos}
}

func (s *RoomService) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomMutationResult[model.RoomModel], error) {
	maxPlayers := DEFAULT_MAX_PLAYERS
	if input.MaxPlayers != nil {
		maxPlayers = *input.MaxPlayers
	}
	if err := validateMaxPlayers(maxPlayers); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	name := input.Name
	if name == "" || strings.TrimSpace(name) == "" {
		owner, err := s.repos.Queries.GetUserByClerkID(ctx, input.OwnerID)
		if err != nil {
			return model.RoomMutationResult[model.RoomModel]{}, err
		}
		name = fmt.Sprintf("Комната %s", owner.Username)
	}
	name, err := normalizeRoomName(name)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	passwordHash, err := roomHelpers.HashPassword(input.Password)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	inviteToken, err := roomHelpers.GenerateInviteToken()
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	room, err := queries.CreateRoom(ctx, db.CreateRoomParams{
		OwnerID:      input.OwnerID,
		Name:         name,
		MaxPlayers:   maxPlayers,
		InviteToken:  inviteToken,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: room.ID,
		UserID: input.OwnerID,
		Role:   ROLE_GM,
	})
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	return model.RoomMutationResult[model.RoomModel]{
		Value: model.ToRoomModel(room, []db.RoomMember{member}, input.OwnerID),
	}, nil
}

func (s *RoomService) ListRooms(ctx context.Context, input model.ListRoomsInput) ([]model.RoomSummaryModel, error) {
	rooms, err := s.repos.Queries.ListRooms(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return model.ToRoomSummaryModels(rooms), nil
}

func (s *RoomService) GetRoom(ctx context.Context, input model.GetRoomInput) (model.RoomModel, error) {
	room, err := s.repos.Queries.GetRoomByID(ctx, db.GetRoomByIDParams{
		ID:     input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomModel{}, ErrRoomNotFound
		}
		return model.RoomModel{}, err
	}

	members, err := s.repos.Queries.ListMembersByRoomID(ctx, db.ListMembersByRoomIDParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModelWithUsernames(room, members, input.UserID), nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomMutationResult[model.RoomModel], error) {
	if err := validateMaxPlayers(input.MaxPlayers); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	var passwordHash *string
	if input.Password != nil {
		if err := validatePassword(*input.Password); err != nil {
			return model.RoomMutationResult[model.RoomModel]{}, err
		}

		hash, err := roomHelpers.HashPassword(*input.Password)
		if err != nil {
			return model.RoomMutationResult[model.RoomModel]{}, err
		}
		passwordHash = &hash
	}

	var name *string
	if input.Name != nil {
		normalizedName, err := normalizeRoomName(*input.Name)
		if err != nil {
			return model.RoomMutationResult[model.RoomModel]{}, err
		}
		name = &normalizedName
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	if input.MaxPlayers < count {
		return model.RoomMutationResult[model.RoomModel]{}, ErrInvalidInput
	}

	room, err := queries.UpdateRoom(ctx, db.UpdateRoomParams{
		Name:         name,
		ID:           input.RoomID,
		OwnerID:      input.OwnerID,
		MaxPlayers:   input.MaxPlayers,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.RoomModel]{}, ErrNotOwner
		}
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.OwnerID, model.EventRoomUpdated, roomHelpers.EmptyEventPayload())
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	return model.RoomMutationResult[model.RoomModel]{
		Value:  model.ToRoomModel(room, nil, input.OwnerID),
		Events: []model.RoomEventModel{event},
	}, nil
}

func (s *RoomService) TransferOwnership(ctx context.Context, input model.TransferOwnershipInput) (model.RoomMutationResult[model.RoomModel], error) {
	if input.NewOwnerID == "" {
		return model.RoomMutationResult[model.RoomModel]{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.TransferRoomOwnership(ctx, db.TransferRoomOwnershipParams{
		ID:         input.RoomID,
		OwnerID:    input.OwnerID,
		NewOwnerID: input.NewOwnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.RoomModel]{}, ErrNotOwner
		}
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	payload, err := roomHelpers.OwnerTransferredPayload(input.OwnerID, input.NewOwnerID)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.OwnerID, model.EventOwnerTransferred, payload)
	if err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomModel]{}, err
	}

	return model.RoomMutationResult[model.RoomModel]{
		Value:  model.ToRoomModel(room, nil, input.OwnerID),
		Events: []model.RoomEventModel{event},
	}, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, input model.DeleteRoomInput) (model.RoomMutationResult[struct{}], error) {
	_, err := s.repos.Queries.DeleteRoom(ctx, db.DeleteRoomParams{
		ID:      input.RoomID,
		OwnerID: input.OwnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[struct{}]{}, ErrNotOwner
		}
		return model.RoomMutationResult[struct{}]{}, err
	}

	return model.RoomMutationResult[struct{}]{}, nil
}

func (s *RoomService) JoinRoom(ctx context.Context, input model.JoinRoomInput) (model.RoomMutationResult[model.RoomMemberModel], error) {
	if !roomHelpers.HasAnyJoinCredential(input) {
		return model.RoomMutationResult[model.RoomMemberModel]{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	meta, err := queries.GetRoomJoinMetaData(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.RoomMemberModel]{}, ErrRoomNotFound
		}
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if !roomHelpers.CanUseJoinCredential(meta.InviteToken, meta.PasswordHash, input) {
		return model.RoomMutationResult[model.RoomMemberModel]{}, ErrRoomNotFound
	}

	_, err = queries.GetMember(ctx, db.GetMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err == nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, ErrAlreadyMember
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if count >= meta.MaxPlayers {
		return model.RoomMutationResult[model.RoomMemberModel]{}, ErrRoomFull
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
		Role:   ROLE_PLAYER,
	})
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	payload, err := roomHelpers.MemberEventPayload(input.UserID, member.Role, "")
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.UserID, model.EventMemberJoined, payload)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	return model.RoomMutationResult[model.RoomMemberModel]{
		Value:  model.ToRoomMemberModel(member),
		Events: []model.RoomEventModel{event},
	}, nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, input model.LeaveRoomInput) (model.RoomMutationResult[model.LeaveRoomResult], error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.GetRoomForUpdate(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, ErrNotMember
		}
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}

	removedMember, err := queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, ErrNotMember
		}
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}

	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}
	if count == 0 {
		if _, err := queries.DeleteRoomByID(ctx, input.RoomID); err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		deletedRoomID := input.RoomID
		return model.RoomMutationResult[model.LeaveRoomResult]{
			Value: model.LeaveRoomResult{DeletedRoomID: &deletedRoomID},
		}, nil
	}

	events := make([]model.RoomEventModel, 0, 2)
	if removedMember.UserID == room.OwnerID {
		nextOwner, err := queries.GetNextRoomOwner(ctx, input.RoomID)
		if err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		if _, err := queries.TransferRoomOwnership(ctx, db.TransferRoomOwnershipParams{
			ID:         input.RoomID,
			OwnerID:    room.OwnerID,
			NewOwnerID: nextOwner.UserID,
		}); err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		payload, err := roomHelpers.OwnerTransferredPayload(room.OwnerID, nextOwner.UserID)
		if err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		event, err := createMutationEvent(ctx, queries, input.RoomID, removedMember.UserID, model.EventOwnerTransferred, payload)
		if err != nil {
			return model.RoomMutationResult[model.LeaveRoomResult]{}, err
		}
		events = append(events, event)
	}

	payload, err := roomHelpers.MemberEventPayload(removedMember.UserID, "", "")
	if err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, removedMember.UserID, model.EventMemberLeft, payload)
	if err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}
	events = append(events, event)

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.LeaveRoomResult]{}, err
	}

	return model.RoomMutationResult[model.LeaveRoomResult]{
		Events: events,
		Value:  model.LeaveRoomResult{},
	}, nil
}

func (s *RoomService) KickMember(ctx context.Context, input model.KickMemberInput) (model.RoomMutationResult[struct{}], error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[struct{}]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.GetRoomForUpdate(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[struct{}]{}, ErrNotMember
		}
		return model.RoomMutationResult[struct{}]{}, err
	}

	if _, err := queries.GetMember(ctx, db.GetMemberParams{
		RoomID: input.RoomID,
		UserID: input.ActorUserID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[struct{}]{}, ErrNotMember
		}
		return model.RoomMutationResult[struct{}]{}, err
	}

	if room.OwnerID != input.ActorUserID {
		return model.RoomMutationResult[struct{}]{}, ErrNotOwner
	}

	if input.TargetUserID == input.ActorUserID {
		return model.RoomMutationResult[struct{}]{}, ErrCannotKickOwner
	}

	_, err = queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.TargetUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[struct{}]{}, ErrNotMember
		}
		return model.RoomMutationResult[struct{}]{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[struct{}]{}, err
	}

	payload, err := roomHelpers.MemberEventPayload(input.TargetUserID, "", "")
	if err != nil {
		return model.RoomMutationResult[struct{}]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.ActorUserID, model.EventMemberKicked, payload)
	if err != nil {
		return model.RoomMutationResult[struct{}]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[struct{}]{}, err
	}
	return model.RoomMutationResult[struct{}]{
		Events: []model.RoomEventModel{event},
	}, nil
}

func (s *RoomService) SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMutationResult[model.RoomMemberModel], error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	member, err := queries.UpdateMemberCharacter(ctx, db.UpdateMemberCharacterParams{
		RoomID:      input.RoomID,
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.RoomMemberModel]{}, ErrCharacterNotOwned
		}

		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	payload, err := roomHelpers.MemberEventPayload(input.UserID, "", member.CharacterID.String())
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.UserID, model.EventMemberCharacterSelected, payload)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	return model.RoomMutationResult[model.RoomMemberModel]{
		Value:  model.ToRoomMemberModel(member),
		Events: []model.RoomEventModel{event},
	}, nil
}

func (s *RoomService) ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMutationResult[model.RoomMemberModel], error) {
	if !IsValidRole(input.Role) {
		return model.RoomMutationResult[model.RoomMemberModel]{}, ErrInvalidInput
	}

	if err := s.EnsureOwner(ctx, input.RoomID, input.ActorUserID); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	member, err := queries.UpdateMemberRole(ctx, db.UpdateMemberRoleParams{
		RoomID: input.RoomID,
		UserID: input.TargetUserID,
		Role:   input.Role,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMutationResult[model.RoomMemberModel]{}, ErrNotMember
		}
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	payload, err := roomHelpers.MemberEventPayload(member.UserID, member.Role, "")
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}
	event, err := createMutationEvent(ctx, queries, input.RoomID, input.ActorUserID, model.EventMemberRoleChanged, payload)
	if err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMutationResult[model.RoomMemberModel]{}, err
	}

	return model.RoomMutationResult[model.RoomMemberModel]{
		Value:  model.ToRoomMemberModel(member),
		Events: []model.RoomEventModel{event},
	}, nil
}
