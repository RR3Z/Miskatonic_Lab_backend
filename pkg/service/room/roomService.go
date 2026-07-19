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

func (s *RoomService) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomModel, error) {
	maxPlayers := DEFAULT_MAX_PLAYERS
	if input.MaxPlayers != nil {
		maxPlayers = *input.MaxPlayers
	}
	if err := validateMaxPlayers(maxPlayers); err != nil {
		return model.RoomModel{}, err
	}
	if err := validatePassword(input.Password); err != nil {
		return model.RoomModel{}, err
	}

	name := input.Name
	if name == "" || strings.TrimSpace(name) == "" {
		owner, err := s.repos.Queries.GetUserByClerkID(ctx, input.OwnerID)
		if err != nil {
			return model.RoomModel{}, err
		}
		name = fmt.Sprintf("Комната %s", owner.Username)
	}
	name, err := normalizeRoomName(name)
	if err != nil {
		return model.RoomModel{}, err
	}

	passwordHash, err := roomHelpers.HashPassword(input.Password)
	if err != nil {
		return model.RoomModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	inviteToken, err := roomHelpers.GenerateInviteToken()
	if err != nil {
		return model.RoomModel{}, err
	}

	room, err := queries.CreateRoom(ctx, db.CreateRoomParams{
		OwnerID:      input.OwnerID,
		Name:         name,
		MaxPlayers:   maxPlayers,
		InviteToken:  inviteToken,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: room.ID,
		UserID: input.OwnerID,
		Role:   ROLE_GM,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, []db.RoomMember{member}, input.OwnerID), nil
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

func (s *RoomService) UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomModel, error) {
	if err := validateMaxPlayers(input.MaxPlayers); err != nil {
		return model.RoomModel{}, err
	}

	var passwordHash *string
	if input.Password != nil {
		if err := validatePassword(*input.Password); err != nil {
			return model.RoomModel{}, err
		}

		hash, err := roomHelpers.HashPassword(*input.Password)
		if err != nil {
			return model.RoomModel{}, err
		}
		passwordHash = &hash
	}

	var name *string
	if input.Name != nil {
		normalizedName, err := normalizeRoomName(*input.Name)
		if err != nil {
			return model.RoomModel{}, err
		}
		name = &normalizedName
	}

	if err := s.EnsureOwner(ctx, input.RoomID, input.OwnerID); err != nil {
		return model.RoomModel{}, err
	}

	count, err := s.repos.Queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.RoomModel{}, err
	}
	if input.MaxPlayers < count {
		return model.RoomModel{}, ErrInvalidInput
	}

	room, err := s.repos.Queries.UpdateRoom(ctx, db.UpdateRoomParams{
		Name:         name,
		ID:           input.RoomID,
		OwnerID:      input.OwnerID,
		MaxPlayers:   input.MaxPlayers,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomModel{}, ErrNotOwner
		}
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, nil, input.OwnerID), nil
}

func (s *RoomService) TransferOwnership(ctx context.Context, input model.TransferOwnershipInput) (model.RoomModel, error) {
	if input.NewOwnerID == "" {
		return model.RoomModel{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomModel{}, err
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
			return model.RoomModel{}, ErrNotOwner
		}
		return model.RoomModel{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, nil, input.OwnerID), nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, input model.DeleteRoomInput) error {
	_, err := s.repos.Queries.DeleteRoom(ctx, db.DeleteRoomParams{
		ID:      input.RoomID,
		OwnerID: input.OwnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotOwner
		}
		return err
	}

	return nil
}

func (s *RoomService) JoinRoom(ctx context.Context, input model.JoinRoomInput) (model.RoomMemberModel, error) {
	if !roomHelpers.HasAnyJoinCredential(input) {
		return model.RoomMemberModel{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMemberModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	meta, err := queries.GetRoomJoinMetaData(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMemberModel{}, ErrRoomNotFound
		}
		return model.RoomMemberModel{}, err
	}

	if !roomHelpers.CanUseJoinCredential(meta.InviteToken, meta.PasswordHash, input) {
		return model.RoomMemberModel{}, ErrRoomNotFound
	}

	_, err = queries.GetMember(ctx, db.GetMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err == nil {
		return model.RoomMemberModel{}, ErrAlreadyMember
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.RoomMemberModel{}, err
	}

	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.RoomMemberModel{}, err
	}

	if count >= meta.MaxPlayers {
		return model.RoomMemberModel{}, ErrRoomFull
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
		Role:   ROLE_PLAYER,
	})
	if err != nil {
		return model.RoomMemberModel{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMemberModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, input model.LeaveRoomInput) (model.LeaveRoomResult, error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.LeaveRoomResult{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.GetRoomForUpdate(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LeaveRoomResult{}, ErrNotMember
		}
		return model.LeaveRoomResult{}, err
	}

	removedMember, err := queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LeaveRoomResult{}, ErrNotMember
		}

		return model.LeaveRoomResult{}, err
	}

	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return model.LeaveRoomResult{}, err
	}
	if count == 0 {
		if _, err := queries.DeleteRoomByID(ctx, input.RoomID); err != nil {
			return model.LeaveRoomResult{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return model.LeaveRoomResult{}, err
		}
		deletedRoomID := input.RoomID
		return model.LeaveRoomResult{DeletedRoomID: &deletedRoomID}, nil
	}

	if removedMember.UserID == room.OwnerID {
		nextOwner, err := queries.GetNextRoomOwner(ctx, input.RoomID)
		if err != nil {
			return model.LeaveRoomResult{}, err
		}

		if _, err := queries.TransferRoomOwnership(ctx, db.TransferRoomOwnershipParams{
			ID:         input.RoomID,
			OwnerID:    room.OwnerID,
			NewOwnerID: nextOwner.UserID,
		}); err != nil {
			return model.LeaveRoomResult{}, err
		}

		payload, err := roomHelpers.OwnerTransferredPayload(room.OwnerID, nextOwner.UserID)
		if err != nil {
			return model.LeaveRoomResult{}, err
		}
		if _, err := queries.CreateRoomEvent(ctx, db.CreateRoomEventParams{
			RoomID:    input.RoomID,
			ActorID:   removedMember.UserID,
			EventType: string(model.EventOwnerTransferred),
			Payload:   payload,
		}); err != nil {
			return model.LeaveRoomResult{}, err
		}
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.LeaveRoomResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.LeaveRoomResult{}, err
	}

	return model.LeaveRoomResult{}, nil
}

func (s *RoomService) KickMember(ctx context.Context, input model.KickMemberInput) error {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.GetRoomForUpdate(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	if _, err := queries.GetMember(ctx, db.GetMemberParams{
		RoomID: input.RoomID,
		UserID: input.ActorUserID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	if room.OwnerID != input.ActorUserID {
		return ErrNotOwner
	}

	if input.TargetUserID == input.ActorUserID {
		return ErrCannotKickOwner
	}

	_, err = queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.TargetUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return err
	}

	count, err := queries.GetRoomMembersCount(ctx, input.RoomID)
	if err != nil {
		return err
	}
	if count == 0 {
		if _, err := queries.DeleteRoomByID(ctx, input.RoomID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *RoomService) SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMemberModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return model.RoomMemberModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMemberModel{}, err
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
			return model.RoomMemberModel{}, ErrCharacterNotOwned
		}

		return model.RoomMemberModel{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMemberModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMemberModel, error) {
	if !IsValidRole(input.Role) {
		return model.RoomMemberModel{}, ErrInvalidInput
	}

	if err := s.EnsureOwner(ctx, input.RoomID, input.ActorUserID); err != nil {
		return model.RoomMemberModel{}, err
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMemberModel{}, err
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
			return model.RoomMemberModel{}, ErrNotMember
		}
		return model.RoomMemberModel{}, err
	}

	if _, err := queries.TouchRoomActivity(ctx, input.RoomID); err != nil {
		return model.RoomMemberModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}
