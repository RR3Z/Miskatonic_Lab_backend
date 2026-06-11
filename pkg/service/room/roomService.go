package room

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
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

func (s *RoomService) CreateRoom(ctx context.Context, params db.CreateRoomParams) (model.RoomModel, error) {
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
	params.InviteToken = inviteToken

	room, err := queries.CreateRoom(ctx, params)
	if err != nil {
		return model.RoomModel{}, err
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: room.ID,
		UserID: params.OwnerID,
		Role:   "gm",
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, []db.RoomMember{member}), nil
}

func (s *RoomService) GetRoom(ctx context.Context, params db.GetRoomByIDParams) (model.RoomModel, error) {
	room, err := s.repos.Queries.GetRoomByID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomModel{}, ErrRoomNotFound
		}
		return model.RoomModel{}, err
	}

	members, err := s.repos.Queries.ListMembersByRoomID(ctx, db.ListMembersByRoomIDParams{
		RoomID: params.ID,
		UserID: params.UserID,
	})
	if err != nil {
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, members), nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, params db.UpdateRoomParams) (model.RoomModel, error) {
	room, err := s.repos.Queries.UpdateRoom(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomModel{}, ErrNotOwner
		}
		return model.RoomModel{}, err
	}

	return model.ToRoomModel(room, nil), nil
}

func (s *RoomService) TransferOwnership(ctx context.Context, params db.TransferRoomOwnershipParams) (model.RoomModel, error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	room, err := queries.TransferRoomOwnership(ctx, params)
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

func (s *RoomService) DeleteRoom(ctx context.Context, params db.DeleteRoomParams) error {
	_, err := s.repos.Queries.DeleteRoom(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotOwner
		}
		return err
	}

	return nil
}

func (s *RoomService) JoinRoom(ctx context.Context, metaParams db.GetRoomMetaDataParams, memberParams db.GetMemberParams) (model.RoomMemberModel, error) {
	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		return model.RoomMemberModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	meta, err := queries.GetRoomMetaData(ctx, metaParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMemberModel{}, ErrRoomNotFound
		}
		return model.RoomMemberModel{}, err
	}

	_, err = queries.GetMember(ctx, memberParams)
	if err == nil {
		return model.RoomMemberModel{}, ErrAlreadyMember
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.RoomMemberModel{}, err
	}

	count, err := queries.GetRoomMembersCount(ctx, metaParams.ID)
	if err != nil {
		return model.RoomMemberModel{}, err
	}

	if count >= meta.MaxPlayers {
		return model.RoomMemberModel{}, ErrRoomFull
	}

	member, err := queries.AddMember(ctx, db.AddMemberParams{
		RoomID: memberParams.RoomID,
		UserID: memberParams.UserID,
		Role:   "player",
	})
	if err != nil {
		return model.RoomMemberModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, params db.RemoveMemberParams) error {
	_, err := s.repos.Queries.RemoveMember(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}

		return err
	}

	return nil
}

func (s *RoomService) KickMember(ctx context.Context, actor db.GetRoomByIDParams, target db.RemoveMemberParams) error {
	room, err := s.repos.Queries.GetRoomByID(ctx, actor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	if room.OwnerID != actor.UserID {
		return ErrNotOwner
	}

	if target.UserID == actor.UserID {
		return ErrCannotKickOwner
	}

	_, err = s.repos.Queries.RemoveMember(ctx, target)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotMember
		}
		return err
	}

	return nil
}

func (s *RoomService) SelectCharacter(ctx context.Context, params db.UpdateMemberCharacterParams) (model.RoomMemberModel, error) {
	member, err := s.repos.Queries.UpdateMemberCharacter(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMemberModel{}, ErrCharacterNotOwned
		}

		return model.RoomMemberModel{}, err
	}

	return model.ToRoomMemberModel(member), nil
}

func (s *RoomService) ChangeRole(ctx context.Context, actor db.GetRoomByIDParams, target db.UpdateMemberRoleParams) (model.RoomMemberModel, error) {
	room, err := s.repos.Queries.GetRoomByID(ctx, actor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RoomMemberModel{}, ErrNotMember
		}
		return model.RoomMemberModel{}, err
	}

	if room.OwnerID != actor.UserID {
		return model.RoomMemberModel{}, ErrNotOwner
	}

	member, err := s.repos.Queries.UpdateMemberRole(ctx, target)
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
