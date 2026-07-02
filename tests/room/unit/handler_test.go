package tests

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/clerk/clerk-sdk-go/v2"
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

func newRoomHandlerTestRouter(roomService room.IRoom) http.Handler {
	h := handler.NewHandler(&service.Service{Room: roomService})
	return h.InitRoutesWithAuth(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
				RegisteredClaims: clerk.RegisteredClaims{Subject: "user_1"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func performRoomRequest(router http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

type fakeRoomHandlerService struct {
	err error

	room   roomModels.RoomModel
	member roomModels.RoomMemberModel

	createCalls int
	createInput roomModels.CreateRoomInput

	getCalls int
	getInput roomModels.GetRoomInput

	updateCalls int
	updateInput roomModels.UpdateRoomInput

	transferCalls int
	transferInput roomModels.TransferOwnershipInput

	deleteCalls int
	deleteInput roomModels.DeleteRoomInput

	joinCalls int
	joinInput roomModels.JoinRoomInput

	leaveCalls int
	leaveInput roomModels.LeaveRoomInput

	kickCalls int
	kickInput roomModels.KickMemberInput

	selectCharacterCalls int
	selectCharacterInput roomModels.SelectCharacterInput

	changeRoleCalls int
	changeRoleInput roomModels.ChangeRoleInput
}

func (f *fakeRoomHandlerService) totalCalls() int {
	return f.createCalls + f.getCalls + f.updateCalls + f.transferCalls + f.deleteCalls + f.joinCalls + f.leaveCalls + f.kickCalls + f.selectCharacterCalls + f.changeRoleCalls
}

func (f *fakeRoomHandlerService) CreateRoom(_ context.Context, input roomModels.CreateRoomInput) (roomModels.RoomModel, error) {
	f.createCalls++
	f.createInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) GetRoom(_ context.Context, input roomModels.GetRoomInput) (roomModels.RoomModel, error) {
	f.getCalls++
	f.getInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) UpdateRoom(_ context.Context, input roomModels.UpdateRoomInput) (roomModels.RoomModel, error) {
	f.updateCalls++
	f.updateInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) TransferOwnership(_ context.Context, input roomModels.TransferOwnershipInput) (roomModels.RoomModel, error) {
	f.transferCalls++
	f.transferInput = input
	return f.room, f.err
}

func (f *fakeRoomHandlerService) DeleteRoom(_ context.Context, input roomModels.DeleteRoomInput) error {
	f.deleteCalls++
	f.deleteInput = input
	return f.err
}

func (f *fakeRoomHandlerService) JoinRoom(_ context.Context, input roomModels.JoinRoomInput) (roomModels.RoomMemberModel, error) {
	f.joinCalls++
	f.joinInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) LeaveRoom(_ context.Context, input roomModels.LeaveRoomInput) error {
	f.leaveCalls++
	f.leaveInput = input
	return f.err
}

func (f *fakeRoomHandlerService) KickMember(_ context.Context, input roomModels.KickMemberInput) error {
	f.kickCalls++
	f.kickInput = input
	return f.err
}

func (f *fakeRoomHandlerService) SelectCharacter(_ context.Context, input roomModels.SelectCharacterInput) (roomModels.RoomMemberModel, error) {
	f.selectCharacterCalls++
	f.selectCharacterInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) ChangeRole(_ context.Context, input roomModels.ChangeRoleInput) (roomModels.RoomMemberModel, error) {
	f.changeRoleCalls++
	f.changeRoleInput = input
	return f.member, f.err
}

func (f *fakeRoomHandlerService) EnsureMember(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeRoomHandlerService) EnsureOwner(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func (f *fakeRoomHandlerService) EnsureCanPublishRoomEvent(_ context.Context, _ pgtype.UUID, _ string) error {
	return f.err
}

func testRoomUnitUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
