package room

import (
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	roomErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room/errors"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/room/helpers"
	roomDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (h *RoomHandler) createRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	var req roomDTO.CreateRoomRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.CreateRoom(r.Context(), roomDTO.CreateRoomInput{
		OwnerID:    userID,
		Name:       req.Name,
		MaxPlayers: req.MaxPlayers,
		Password:   req.Password,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to create room")
	}
	h.presence.AwaitFirstConnection(result.Value.ID, userID)

	utils.WriteJSON(w, http.StatusCreated, result.Value)
	return nil
}

func (h *RoomHandler) listRooms(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())

	result, err := h.service.ListRooms(r.Context(), roomDTO.ListRoomsInput{UserID: userID})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to list rooms")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *RoomHandler) getRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	result, err := h.service.GetRoom(r.Context(), roomDTO.GetRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to get room")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *RoomHandler) updateRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	var req roomDTO.UpdateRoomRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.UpdateRoom(r.Context(), roomDTO.UpdateRoomInput{
		RoomID:     roomID,
		OwnerID:    userID,
		Name:       req.Name,
		MaxPlayers: req.MaxPlayers,
		Password:   roomHelpers.OptionalPassword(req.Password),
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to update room")
	}

	h.broadcastRoomEvents(result.Events)

	utils.WriteJSON(w, http.StatusOK, result.Value)
	return nil
}

func (h *RoomHandler) transferRoomOwnership(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	var req roomDTO.TransferRoomOwnershipRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.TransferOwnership(r.Context(), roomDTO.TransferOwnershipInput{
		RoomID:     roomID,
		OwnerID:    userID,
		NewOwnerID: req.UserID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to transfer room ownership")
	}

	h.broadcastRoomEvents(result.Events)

	utils.WriteJSON(w, http.StatusOK, result.Value)
	return nil
}

func (h *RoomHandler) deleteRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	_, err = h.service.DeleteRoom(r.Context(), roomDTO.DeleteRoomInput{
		RoomID:  roomID,
		OwnerID: userID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to delete room")
	}

	h.presence.ForgetRoom(roomID)
	h.closeRoom(roomID, "room deleted")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *RoomHandler) listSelectedCharacters(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	result, err := h.service.ListSelectedCharacters(r.Context(), roomDTO.ListSelectedCharactersInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to list selected characters")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *RoomHandler) listRoomEvents(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	limit, err := roomHelpers.OptionalInt32Query(r, "limit")
	if err != nil {
		return roomErrors.InvalidInputError("invalid events limit", err)
	}
	afterSequence, err := roomHelpers.OptionalInt64Query(r, "after")
	if err != nil || (afterSequence != nil && *afterSequence < 0) {
		return roomErrors.InvalidInputError("invalid events cursor", err)
	}

	input := roomDTO.ListRoomEventsInput{
		RoomID: roomID,
		UserID: userID,
	}
	if limit != nil {
		input.Limit = *limit
	}
	if afterSequence != nil {
		input.AfterSequence = *afterSequence
	}

	result, err := h.service.ListRoomEvents(r.Context(), input)
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to list room events")
	}

	utils.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *RoomHandler) joinRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	var req roomDTO.JoinRoomRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.JoinRoom(
		r.Context(),
		roomDTO.JoinRoomInput{
			RoomID:      roomID,
			UserID:      userID,
			InviteToken: req.InviteToken,
			Password:    req.Password,
		},
	)
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to join room")
	}

	h.presence.AwaitFirstConnection(roomID, userID)
	h.broadcastRoomEvents(result.Events)

	utils.WriteJSON(w, http.StatusOK, result.Value)
	return nil
}

func (h *RoomHandler) leaveRoom(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	result, err := h.service.LeaveRoom(r.Context(), roomDTO.LeaveRoomInput{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to leave room")
	}

	h.presence.Forget(roomID, userID)
	h.closeUser(roomID, userID, "room left")
	h.broadcastRoomEvents(result.Events)

	if result.Value.DeletedRoomID != nil {
		h.closeRoom(*result.Value.DeletedRoomID, "room deleted")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *RoomHandler) kickMember(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}
	targetUserID := chi.URLParam(r, "userID")

	result, err := h.service.KickMember(
		r.Context(),
		roomDTO.KickMemberInput{
			RoomID:       roomID,
			ActorUserID:  userID,
			TargetUserID: targetUserID,
		},
	)
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to kick member")
	}

	h.presence.Forget(roomID, targetUserID)
	h.closeUser(roomID, targetUserID, "removed from room")
	h.broadcastRoomEvents(result.Events)

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *RoomHandler) selectCharacter(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}

	var req roomDTO.SelectCharacterRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.SelectCharacter(r.Context(), roomDTO.SelectCharacterInput{
		RoomID:      roomID,
		UserID:      userID,
		CharacterID: req.CharacterID,
	})
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to select character")
	}

	h.broadcastRoomEvents(result.Events)

	utils.WriteJSON(w, http.StatusOK, result.Value)
	return nil
}

func (h *RoomHandler) changeRole(w http.ResponseWriter, r *http.Request) *myErrors.AppError {
	userID := utils.GetUserIDFromContext(r.Context())
	roomID, err := roomHelpers.GetRoomIDFromRequest(r)
	if err != nil {
		return roomErrors.InvalidIDError(err)
	}
	targetUserID := chi.URLParam(r, "userID")

	var req roomDTO.ChangeRoleRequest
	if appErr := roomHelpers.DecodeJSON(r, &req); appErr != nil {
		return appErr
	}

	result, err := h.service.ChangeRole(
		r.Context(),
		roomDTO.ChangeRoleInput{
			RoomID:       roomID,
			ActorUserID:  userID,
			TargetUserID: targetUserID,
			Role:         req.Role,
		},
	)
	if err != nil {
		return roomErrors.MapServiceError(err, "failed to change role")
	}

	h.broadcastRoomEvents(result.Events)

	utils.WriteJSON(w, http.StatusOK, result.Value)
	return nil
}
