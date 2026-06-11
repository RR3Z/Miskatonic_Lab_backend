package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	var req model.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	maxPlayers := int32(7)
	if req.MaxPlayers != nil {
		maxPlayers = *req.MaxPlayers
	}
	if maxPlayers < 1 {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "max_players must be greater than 0",
			Err:     nil,
		}
	}

	result, err := h.services.Room.CreateRoom(r.Context(), db.CreateRoomParams{
		OwnerID:    userID,
		MaxPlayers: maxPlayers,
	})
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create room",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusCreated, result)
	return nil
}

func (h *Handler) getRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	result, err := h.services.Room.GetRoom(r.Context(), db.GetRoomByIDParams{
		ID:     roomID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, room.ErrRoomNotFound) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "room not found",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get room",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) updateRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	var req model.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	if req.MaxPlayers < 1 {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "max_players must be greater than 0",
			Err:     nil,
		}
	}

	result, err := h.services.Room.UpdateRoom(r.Context(), db.UpdateRoomParams{
		ID:         roomID,
		OwnerID:    userID,
		MaxPlayers: req.MaxPlayers,
	})
	if err != nil {
		if errors.Is(err, room.ErrNotOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "only the room owner can update the room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update room",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) transferRoomOwnership(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	var req model.TransferRoomOwnershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	if req.UserID == "" {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "user_id is required",
			Err:     nil,
		}
	}

	result, err := h.services.Room.TransferOwnership(r.Context(), db.TransferRoomOwnershipParams{
		ID:         roomID,
		OwnerID:    userID,
		NewOwnerID: req.UserID,
	})
	if err != nil {
		if errors.Is(err, room.ErrNotOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "only the room owner can transfer ownership",
				Err:     err,
			}
		}
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to transfer room ownership",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) deleteRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	err = h.services.Room.DeleteRoom(r.Context(), db.DeleteRoomParams{
		ID:      roomID,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, room.ErrNotOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "only the room owner can delete the room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete room",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	var req model.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	if req.InviteToken == "" {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invite_token is required",
			Err:     nil,
		}
	}

	result, err := h.services.Room.JoinRoom(
		r.Context(),
		db.GetRoomMetaDataParams{
			ID:          roomID,
			InviteToken: req.InviteToken,
		},
		db.GetMemberParams{
			RoomID: roomID,
			UserID: userID,
		},
	)
	if err != nil {
		if errors.Is(err, room.ErrRoomNotFound) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "room not found",
				Err:     err,
			}
		}

		if errors.Is(err, room.ErrRoomFull) {
			return &myErrors.AppError{
				Status:  http.StatusConflict,
				Message: "room is full",
				Err:     err,
			}
		}

		if errors.Is(err, room.ErrAlreadyMember) {
			return &myErrors.AppError{
				Status:  http.StatusConflict,
				Message: "already a member of this room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to join room",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) leaveRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	err = h.services.Room.LeaveRoom(r.Context(), db.RemoveMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, room.ErrNotMember) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "not a member of this room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to leave room",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) kickMember(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}
	targetUserID := chi.URLParam(r, "userID")

	err = h.services.Room.KickMember(
		r.Context(),
		db.GetRoomByIDParams{
			ID:     roomID,
			UserID: userID,
		},
		db.RemoveMemberParams{
			RoomID: roomID,
			UserID: targetUserID,
		},
	)
	if err != nil {
		if errors.Is(err, room.ErrNotOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "only the room owner can kick members",
				Err:     err,
			}
		}

		if errors.Is(err, room.ErrCannotKickOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "cannot kick the room owner",
				Err:     err,
			}
		}

		if errors.Is(err, room.ErrNotMember) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "user is not a member of this room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to kick member",
			Err:     err,
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) selectCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}

	var req model.SelectCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	result, err := h.services.Room.SelectCharacter(r.Context(), db.UpdateMemberCharacterParams{
		RoomID:      roomID,
		UserID:      userID,
		CharacterID: req.CharacterID,
	})
	if err != nil {
		if errors.Is(err, room.ErrCharacterNotOwned) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "character does not belong to you",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to select character",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) changeRole(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid room id",
			Err:     err,
		}
	}
	targetUserID := chi.URLParam(r, "userID")

	var req model.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Err:     err,
		}
	}

	if req.Role != "player" && req.Role != "gm" {
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Message: "role must be 'player' or 'gm'",
			Err:     nil,
		}
	}

	result, err := h.services.Room.ChangeRole(
		r.Context(),
		db.GetRoomByIDParams{
			ID:     roomID,
			UserID: userID,
		},
		db.UpdateMemberRoleParams{
			RoomID: roomID,
			UserID: targetUserID,
			Role:   req.Role,
		},
	)
	if err != nil {
		if errors.Is(err, room.ErrNotOwner) {
			return &myErrors.AppError{
				Status:  http.StatusForbidden,
				Message: "only the room owner can change roles",
				Err:     err,
			}
		}

		if errors.Is(err, room.ErrNotMember) {
			return &myErrors.AppError{
				Status:  http.StatusNotFound,
				Message: "user is not a member of this room",
				Err:     err,
			}
		}

		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to change role",
			Err:     err,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func getRoomIDFromRequest(r *http.Request) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(chi.URLParam(r, "roomID")); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}
