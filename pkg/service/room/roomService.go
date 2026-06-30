package room

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
)

type RoomService struct {
	repos *repository.Repository
}

func NewRoomService(repos *repository.Repository) *RoomService {
	return &RoomService{repos: repos}
}

func (s *RoomService) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.RoomModel, error) {
	maxPlayers := DefaultMaxPlayers
	if input.MaxPlayers != nil {
		maxPlayers = *input.MaxPlayers
	}
	if maxPlayers < 1 {
		return model.RoomModel{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	inviteToken, err := generateInviteToken()
	if err != nil {
		return model.RoomModel{}, err
	}

	room, err := queries.CreateRoom(ctx, db.CreateRoomParams{
		OwnerID:     input.OwnerID,
		MaxPlayers:  maxPlayers,
		InviteToken: inviteToken,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: room.ID,
		UserID: input.OwnerID,
		Role:   RoleGM,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, []db.RoomMember{member}), nil
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

	return model.ToRoomModel(room, members), nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, input model.UpdateRoomInput) (model.RoomModel, error) {
	if input.MaxPlayers < 1 {
		return model.RoomModel{}, ErrInvalidInput
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
		ID:         input.RoomID,
		OwnerID:    input.OwnerID,
		MaxPlayers: input.MaxPlayers,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomModel{}, ErrNotOwner
		}
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, nil), nil
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

	if err := tx.Commit(ctx); err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, nil), nil
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
	if input.InviteToken == "" {
		return model.RoomMemberModel{}, ErrInvalidInput
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMemberModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	meta, err := queries.GetRoomMetaData(ctx, db.GetRoomMetaDataParams{
		ID:          input.RoomID,
		InviteToken: input.InviteToken,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMemberModel{}, ErrRoomNotFound
		}
		return model.RoomMemberModel{}, err
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
		Role:   RolePlayer,
	})
	if err != nil {
		return model.RoomMemberModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, input model.LeaveRoomInput) error {
	_, err := s.repos.Queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}

		return err
	}

	return nil
}

func (s *RoomService) KickMember(ctx context.Context, input model.KickMemberInput) error {
	room, err := s.repos.Queries.GetRoomByID(ctx, db.GetRoomByIDParams{
		ID:     input.RoomID,
		UserID: input.ActorUserID,
	})
	if err != nil {
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

	_, err = s.repos.Queries.RemoveMember(ctx, db.RemoveMemberParams{
		RoomID: input.RoomID,
		UserID: input.TargetUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	return nil
}

func (s *RoomService) SelectCharacter(ctx context.Context, input model.SelectCharacterInput) (model.RoomMemberModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return model.RoomMemberModel{}, err
	}

	member, err := s.repos.Queries.UpdateMemberCharacter(ctx, db.UpdateMemberCharacterParams{
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

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) ChangeRole(ctx context.Context, input model.ChangeRoleInput) (model.RoomMemberModel, error) {
	if !IsValidRole(input.Role) {
		return model.RoomMemberModel{}, ErrInvalidInput
	}

	if err := s.EnsureOwner(ctx, input.RoomID, input.ActorUserID); err != nil {
		return model.RoomMemberModel{}, err
	}

	member, err := s.repos.Queries.UpdateMemberRole(ctx, db.UpdateMemberRoleParams{
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

	return model.ToRoomMemberModel(member), nil
}

func generateInviteToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
