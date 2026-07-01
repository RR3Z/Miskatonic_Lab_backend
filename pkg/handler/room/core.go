package room

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	var req model.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.CreateRoom(r.Context(), model.CreateRoomInput{
		OwnerID:    userID,
		MaxPlayers: req.MaxPlayers,
	})
	if err != nil {
		return mapServiceError(err, "failed to create room")
	}

	utils.WriteJSON(w, http.StatusCreated, result)
	return nil
}

func (h *Handler) getRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	result, err := h.rooms.GetRoom(r.Context(), model.GetRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to get room")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) updateRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.UpdateRoom(r.Context(), model.UpdateRoomInput{
		RoomID:     roomID,
		OwnerID:    userID,
		MaxPlayers: req.MaxPlayers,
	})
	if err != nil {
		return mapServiceError(err, "failed to update room")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) transferRoomOwnership(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.TransferRoomOwnershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.TransferOwnership(r.Context(), model.TransferOwnershipInput{
		RoomID:     roomID,
		OwnerID:    userID,
		NewOwnerID: req.UserID,
	})
	if err != nil {
		return mapServiceError(err, "failed to transfer room ownership")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) deleteRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	err = h.rooms.DeleteRoom(r.Context(), model.DeleteRoomInput{
		RoomID:  roomID,
		OwnerID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to delete room")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.JoinRoom(
		r.Context(),
		model.JoinRoomInput{
			RoomID:      roomID,
			UserID:      userID,
			InviteToken: req.InviteToken,
		},
	)
	if err != nil {
		return mapServiceError(err, "failed to join room")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) leaveRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	err = h.rooms.LeaveRoom(r.Context(), model.LeaveRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to leave room")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) kickMember(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}
	targetUserID := chi.URLParam(r, "userID")

	err = h.rooms.KickMember(
		r.Context(),
		model.KickMemberInput{
			RoomID:       roomID,
			ActorUserID:  userID,
			TargetUserID: targetUserID,
		},
	)
	if err != nil {
		return mapServiceError(err, "failed to kick member")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) selectCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.SelectCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.SelectCharacter(r.Context(), model.SelectCharacterInput{
		RoomID:      roomID,
		UserID:      userID,
		CharacterID: req.CharacterID,
	})
	if err != nil {
		return mapServiceError(err, "failed to select character")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) changeRole(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}
	targetUserID := chi.URLParam(r, "userID")

	var req model.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.rooms.ChangeRole(
		r.Context(),
		model.ChangeRoleInput{
			RoomID:       roomID,
			ActorUserID:  userID,
			TargetUserID: targetUserID,
			Role:         req.Role,
		},
	)
	if err != nil {
		return mapServiceError(err, "failed to change role")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}
