package room

import (
	"encoding/json"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (h *RoomHandler) createRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	var req model.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.service.CreateRoom(r.Context(), model.CreateRoomInput{
		OwnerID:    userID,
		MaxPlayers: req.MaxPlayers,
	})
	if err != nil {
		return mapServiceError(err, "failed to create room")
	}

	utils.WriteJSON(w, http.StatusCreated, result)
	return nil
}

func (h *RoomHandler) getRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	result, err := h.service.GetRoom(r.Context(), model.GetRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to get room")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *RoomHandler) updateRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.service.UpdateRoom(r.Context(), model.UpdateRoomInput{
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

func (h *RoomHandler) transferRoomOwnership(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.TransferRoomOwnershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.service.TransferOwnership(r.Context(), model.TransferOwnershipInput{
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

func (h *RoomHandler) deleteRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	err = h.service.DeleteRoom(r.Context(), model.DeleteRoomInput{
		RoomID:  roomID,
		OwnerID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to delete room")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *RoomHandler) joinRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.service.JoinRoom(
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

func (h *RoomHandler) leaveRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	err = h.service.LeaveRoom(r.Context(), model.LeaveRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return mapServiceError(err, "failed to leave room")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *RoomHandler) kickMember(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}
	targetUserID := chi.URLParam(r, "userID")

	err = h.service.KickMember(
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

func (h *RoomHandler) selectCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := getRoomIDFromRequest(r)
	if err != nil {
		return invalidIDError(err)
	}

	var req model.SelectCharacterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return invalidInputError("invalid request body", err)
	}

	result, err := h.service.SelectCharacter(r.Context(), model.SelectCharacterInput{
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

func (h *RoomHandler) changeRole(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
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

	result, err := h.service.ChangeRole(
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
