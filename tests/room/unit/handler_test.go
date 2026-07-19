package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	characterModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateRoomDefaultsMaxPlayersAndPassesUserID(t *testing.T) {
	roomService := &fakeRoomHandlerService{room: roomModels.RoomModel{OwnerID: "user_1", MaxPlayers: 7}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", `{"password":"keeper-password"}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, roomService.createCalls)
	require.Equal(t, "user_1", roomService.createInput.OwnerID)
	require.Nil(t, roomService.createInput.MaxPlayers)
	require.Equal(t, "keeper-password", roomService.createInput.Password)
}

func TestCreateRoomPassesOptionalName(t *testing.T) {
	roomService := &fakeRoomHandlerService{room: roomModels.RoomModel{OwnerID: "user_1"}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", `{"name":"  Masks of Nyarlathotep  ","password":"keeper-password"}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, "  Masks of Nyarlathotep  ", roomService.createInput.Name)
}

func TestListRoomsPassesCurrentUserAndReturnsSummaries(t *testing.T) {
	roomService := &fakeRoomHandlerService{rooms: []roomModels.RoomSummaryModel{{
		ID:          testRoomUnitUUID("11111111-1111-1111-1111-111111111111"),
		Name:        "Open Table",
		MaxPlayers:  7,
		MemberCount: 2,
		CreatedAt:   pgtype.Timestamptz{Time: time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC), Valid: true},
		IsMember:    true,
	}}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodGet, "/api/rooms/", "")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.listCalls)
	require.Equal(t, "user_1", roomService.listInput.UserID)
	require.JSONEq(t, `[{"id":"11111111-1111-1111-1111-111111111111","name":"Open Table","max_players":7,"member_count":2,"created_at":"2026-01-02T03:04:05Z","is_member":true}]`, recorder.Body.String())
	require.NotContains(t, recorder.Body.String(), "invite_token")
	require.NotContains(t, recorder.Body.String(), "password")
}

func TestCreateRoomRejectsInvalidBodyBeforeService(t *testing.T) {
	roomService := &fakeRoomHandlerService{}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", `{"max_players":`)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, roomService.createCalls)
}

func TestCreateRoomPassesMaxPlayersValidationToService(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		maxPlayers int32
	}{
		{name: "zero max players", body: `{"max_players":0,"password":"keeper-password"}`, maxPlayers: 0},
		{name: "negative max players", body: `{"max_players":-1,"password":"keeper-password"}`, maxPlayers: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{err: room.ErrInvalidInput}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", tt.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Equal(t, 1, roomService.createCalls)
			require.NotNil(t, roomService.createInput.MaxPlayers)
			require.Equal(t, tt.maxPlayers, *roomService.createInput.MaxPlayers)
			require.JSONEq(t, `{"code":"room.invalid_input","message":"invalid room input"}`, recorder.Body.String())
		})
	}
}

func TestCreateRoomPassesMissingPasswordToService(t *testing.T) {
	roomService := &fakeRoomHandlerService{err: room.ErrInvalidPassword}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", `{}`)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 1, roomService.createCalls)
	require.Empty(t, roomService.createInput.Password)
	require.JSONEq(t, `{"code":"room.invalid_input","message":"invalid room input"}`, recorder.Body.String())
}

func TestRoomRoutesRejectInvalidRoomID(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get", method: http.MethodGet, path: "/api/rooms/not-a-uuid/"},
		{name: "update", method: http.MethodPut, path: "/api/rooms/not-a-uuid/", body: `{"max_players":5}`},
		{name: "delete", method: http.MethodDelete, path: "/api/rooms/not-a-uuid/"},
		{name: "characters", method: http.MethodGet, path: "/api/rooms/not-a-uuid/characters"},
		{name: "events", method: http.MethodGet, path: "/api/rooms/not-a-uuid/events"},
		{name: "websocket", method: http.MethodGet, path: "/api/rooms/not-a-uuid/ws"},
		{name: "transfer owner", method: http.MethodPut, path: "/api/rooms/not-a-uuid/owner", body: `{"user_id":"user_2"}`},
		{name: "join", method: http.MethodPost, path: "/api/rooms/not-a-uuid/join", body: `{"invite_token":"token"}`},
		{name: "leave", method: http.MethodDelete, path: "/api/rooms/not-a-uuid/leave"},
		{name: "kick", method: http.MethodDelete, path: "/api/rooms/not-a-uuid/kick/user_2"},
		{name: "select character", method: http.MethodPut, path: "/api/rooms/not-a-uuid/character", body: `{"character_id":"33333333-3333-3333-3333-333333333333"}`},
		{name: "change role", method: http.MethodPut, path: "/api/rooms/not-a-uuid/members/user_2/role", body: `{"role":"gm"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, tt.method, tt.path, tt.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, roomService.totalCalls())
		})
	}
}

func TestListSelectedCharactersPassesParams(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	memberID := testRoomUnitUUID("22222222-2222-2222-2222-222222222222")
	characterID := testRoomUnitUUID("33333333-3333-3333-3333-333333333333")
	roomService := &fakeRoomHandlerService{
		selectedCharacters: []roomModels.SelectedCharacterModel{{
			MemberID: memberID,
			UserID:   "user_1",
			Role:     room.ROLE_PLAYER,
			Character: characterModels.CharacterModel{
				CharacterShortModel: characterModels.CharacterShortModel{
					ID:     characterID,
					UserID: "user_1",
					Name:   "Investigator",
				},
			},
		}},
	}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodGet, "/api/rooms/11111111-1111-1111-1111-111111111111/characters", "")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.listSelectedCharactersCalls)
	require.Equal(t, roomID, roomService.listSelectedCharactersInput.RoomID)
	require.Equal(t, "user_1", roomService.listSelectedCharactersInput.UserID)

	var response []struct {
		MemberID  pgtype.UUID `json:"member_id"`
		UserID    string      `json:"user_id"`
		Role      string      `json:"role"`
		Character struct {
			ID     pgtype.UUID `json:"id"`
			UserID string      `json:"user_id"`
			Name   string      `json:"name"`
		} `json:"character"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Len(t, response, 1)
	require.Equal(t, memberID, response[0].MemberID)
	require.Equal(t, "user_1", response[0].UserID)
	require.Equal(t, room.ROLE_PLAYER, response[0].Role)
	require.Equal(t, characterID, response[0].Character.ID)
	require.Equal(t, "user_1", response[0].Character.UserID)
	require.Equal(t, "Investigator", response[0].Character.Name)
}

func TestListRoomEventsPassesParams(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	eventID := testRoomUnitUUID("22222222-2222-2222-2222-222222222222")
	roomService := &fakeRoomHandlerService{
		events: []roomModels.RoomEventModel{{
			ID:      eventID,
			RoomID:  roomID,
			ActorID: "user_2",
			Type:    "chat.message",
			Payload: []byte(`{"text":"hello"}`),
		}},
	}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodGet, "/api/rooms/11111111-1111-1111-1111-111111111111/events?limit=25", "")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.listEventsCalls)
	require.Equal(t, roomID, roomService.listEventsInput.RoomID)
	require.Equal(t, "user_1", roomService.listEventsInput.UserID)
	require.Equal(t, int32(25), roomService.listEventsInput.Limit)
	require.JSONEq(t, `[{"id":"22222222-2222-2222-2222-222222222222","room_id":"11111111-1111-1111-1111-111111111111","actor_id":"user_2","type":"chat.message","payload":{"text":"hello"},"created_at":null}]`, recorder.Body.String())
}

func TestListRoomEventsRejectsInvalidLimitBeforeService(t *testing.T) {
	roomService := &fakeRoomHandlerService{}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodGet, "/api/rooms/11111111-1111-1111-1111-111111111111/events?limit=nope", "")

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, roomService.listEventsCalls)
}

func TestListRoomEventsPassesOptionalLimitValues(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantLimit int32
	}{
		{name: "missing", target: "/api/rooms/11111111-1111-1111-1111-111111111111/events", wantLimit: 0},
		{name: "negative", target: "/api/rooms/11111111-1111-1111-1111-111111111111/events?limit=-5", wantLimit: -5},
		{name: "large", target: "/api/rooms/11111111-1111-1111-1111-111111111111/events?limit=500", wantLimit: 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, http.MethodGet, tt.target, "")

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, 1, roomService.listEventsCalls)
			require.Equal(t, tt.wantLimit, roomService.listEventsInput.Limit)
		})
	}
}

func TestRoomRoutesRejectInvalidJSONBeforeService(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "update", method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/"},
		{name: "join", method: http.MethodPost, path: "/api/rooms/11111111-1111-1111-1111-111111111111/join"},
		{name: "transfer owner", method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/owner"},
		{name: "select character", method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/character"},
		{name: "change role", method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, tt.method, tt.path, `{"broken":`)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, roomService.totalCalls())
		})
	}
}

func TestRoomWebsocketMembershipFailureDoesNotUpgrade(t *testing.T) {
	roomService := &fakeRoomHandlerService{err: room.ErrNotMember}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodGet, "/api/rooms/11111111-1111-1111-1111-111111111111/ws", "")

	require.Equal(t, http.StatusForbidden, recorder.Code)
	require.Equal(t, 1, roomService.touchActivityCalls)
	require.Equal(t, "user_1", roomService.touchActivityInput.UserID)
	require.Zero(t, roomService.totalCalls())
}

func TestJoinRoomAcceptsInviteTokenAndPassesParams(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	roomService := &fakeRoomHandlerService{member: roomModels.RoomMemberModel{RoomID: roomID, UserID: "user_1"}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/11111111-1111-1111-1111-111111111111/join", `{"invite_token":"token_1"}`)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.joinCalls)
	require.Equal(t, roomID, roomService.joinInput.RoomID)
	require.Equal(t, "token_1", roomService.joinInput.InviteToken)
	require.Empty(t, roomService.joinInput.Password)
	require.Equal(t, "user_1", roomService.joinInput.UserID)
}

func TestJoinRoomAcceptsPasswordAndPassesParams(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	roomService := &fakeRoomHandlerService{member: roomModels.RoomMemberModel{RoomID: roomID, UserID: "user_1"}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/11111111-1111-1111-1111-111111111111/join", `{"password":"keeper-password"}`)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.joinCalls)
	require.Equal(t, roomID, roomService.joinInput.RoomID)
	require.Empty(t, roomService.joinInput.InviteToken)
	require.Equal(t, "keeper-password", roomService.joinInput.Password)
	require.Equal(t, "user_1", roomService.joinInput.UserID)
}

func TestJoinRoomPassesMissingCredentialsToService(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "missing", body: `{}`},
		{name: "blank", body: `{"invite_token":"","password":""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{err: room.ErrInvalidInput}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/11111111-1111-1111-1111-111111111111/join", tt.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Equal(t, 1, roomService.joinCalls)
			require.Empty(t, roomService.joinInput.InviteToken)
			require.Empty(t, roomService.joinInput.Password)
			require.JSONEq(t, `{"code":"room.invalid_input","message":"invalid room input"}`, recorder.Body.String())
		})
	}
}

func TestUpdateRoomPassesOptionalPassword(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	roomService := &fakeRoomHandlerService{room: roomModels.RoomModel{ID: roomID, OwnerID: "user_1", MaxPlayers: 5}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/", `{"max_players":5,"password":"new-password"}`)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.updateCalls)
	require.Equal(t, roomID, roomService.updateInput.RoomID)
	require.Equal(t, "user_1", roomService.updateInput.OwnerID)
	require.Equal(t, int32(5), roomService.updateInput.MaxPlayers)
	require.NotNil(t, roomService.updateInput.Password)
	require.Equal(t, "new-password", *roomService.updateInput.Password)
}

func TestTransferOwnershipPassesMissingUserIDToServiceAndMapsNotOwner(t *testing.T) {
	roomService := &fakeRoomHandlerService{err: room.ErrInvalidInput}
	router := newRoomHandlerTestRouter(roomService)
	recorder := performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/owner", `{}`)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 1, roomService.transferCalls)
	require.Empty(t, roomService.transferInput.NewOwnerID)
	require.JSONEq(t, `{"code":"room.invalid_input","message":"invalid room input"}`, recorder.Body.String())

	roomService = &fakeRoomHandlerService{err: room.ErrNotOwner}
	router = newRoomHandlerTestRouter(roomService)
	recorder = performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/owner", `{"user_id":"user_2"}`)
	require.Equal(t, http.StatusForbidden, recorder.Code)
	require.Equal(t, 1, roomService.transferCalls)
	require.Equal(t, "user_1", roomService.transferInput.OwnerID)
	require.Equal(t, "user_2", roomService.transferInput.NewOwnerID)
	require.JSONEq(t, `{"code":"room.not_owner","message":"only the room owner can perform this action"}`, recorder.Body.String())
}

func TestChangeRolePassesInvalidRoleToService(t *testing.T) {
	roomService := &fakeRoomHandlerService{err: room.ErrInvalidInput}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role", `{"role":"keeper"}`)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 1, roomService.changeRoleCalls)
	require.Equal(t, "keeper", roomService.changeRoleInput.Role)
	require.JSONEq(t, `{"code":"room.invalid_input","message":"invalid room input"}`, recorder.Body.String())
}

func TestRoomHandlerErrorMappings(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
		method     string
		path       string
		body       string
	}{
		{name: "get room not found", err: room.ErrRoomNotFound, wantStatus: http.StatusNotFound, wantBody: `{"code":"room.not_found","message":"room not found"}`, method: http.MethodGet, path: "/api/rooms/11111111-1111-1111-1111-111111111111/"},
		{name: "join room full", err: room.ErrRoomFull, wantStatus: http.StatusConflict, wantBody: `{"code":"room.full","message":"room is full"}`, method: http.MethodPost, path: "/api/rooms/11111111-1111-1111-1111-111111111111/join", body: `{"invite_token":"token"}`},
		{name: "join already member", err: room.ErrAlreadyMember, wantStatus: http.StatusConflict, wantBody: `{"code":"room.already_member","message":"already a member of this room"}`, method: http.MethodPost, path: "/api/rooms/11111111-1111-1111-1111-111111111111/join", body: `{"invite_token":"token"}`},
		{name: "leave not member", err: room.ErrNotMember, wantStatus: http.StatusNotFound, wantBody: `{"code":"room.not_member","message":"not a member of this room"}`, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/leave"},
		{name: "kick not owner", err: room.ErrNotOwner, wantStatus: http.StatusForbidden, wantBody: `{"code":"room.not_owner","message":"only the room owner can perform this action"}`, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/kick/user_2"},
		{name: "kick owner self", err: room.ErrCannotKickOwner, wantStatus: http.StatusForbidden, wantBody: `{"code":"room.cannot_kick_owner","message":"cannot kick the room owner"}`, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/kick/user_2"},
		{name: "select character not owned", err: room.ErrCharacterNotOwned, wantStatus: http.StatusForbidden, wantBody: `{"code":"room.character_not_owned","message":"character does not belong to you"}`, method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/character", body: `{"character_id":"33333333-3333-3333-3333-333333333333"}`},
		{name: "change role generic", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantBody: `{"code":"common.internal_error","message":"failed to change role"}`, method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role", body: `{"role":"gm"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{err: tt.err}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, tt.method, tt.path, tt.body)

			require.Equal(t, tt.wantStatus, recorder.Code)
			require.JSONEq(t, tt.wantBody, recorder.Body.String())
		})
	}
}
