package tests

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateRoomDefaultsMaxPlayersAndPassesUserID(t *testing.T) {
	roomService := &fakeRoomHandlerService{room: model.RoomModel{OwnerID: "user_1", MaxPlayers: 7}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", `{}`)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, 1, roomService.createCalls)
	require.Equal(t, "user_1", roomService.createParams.OwnerID)
	require.Equal(t, int32(7), roomService.createParams.MaxPlayers)
}

func TestCreateRoomRejectsInvalidBodyAndMaxPlayers(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{"max_players":`},
		{name: "zero max players", body: `{"max_players":0}`},
		{name: "negative max players", body: `{"max_players":-1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/", tt.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, roomService.createCalls)
		})
	}
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

func TestJoinRoomRequiresInviteTokenAndPassesParams(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	roomService := &fakeRoomHandlerService{member: model.RoomMemberModel{RoomID: roomID, UserID: "user_1"}}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/11111111-1111-1111-1111-111111111111/join", `{"invite_token":"token_1"}`)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, roomService.joinCalls)
	require.Equal(t, roomID, roomService.joinMeta.ID)
	require.Equal(t, "token_1", roomService.joinMeta.InviteToken)
	require.Equal(t, roomID, roomService.joinMember.RoomID)
	require.Equal(t, "user_1", roomService.joinMember.UserID)
}

func TestJoinRoomRejectsMissingOrBlankInviteToken(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "missing", body: `{}`},
		{name: "blank", body: `{"invite_token":""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, http.MethodPost, "/api/rooms/11111111-1111-1111-1111-111111111111/join", tt.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, roomService.joinCalls)
		})
	}
}

func TestTransferOwnershipRequiresUserIDAndMapsNotOwner(t *testing.T) {
	router := newRoomHandlerTestRouter(&fakeRoomHandlerService{})
	recorder := performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/owner", `{}`)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	roomService := &fakeRoomHandlerService{err: room.ErrNotOwner}
	router = newRoomHandlerTestRouter(roomService)
	recorder = performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/owner", `{"user_id":"user_2"}`)
	require.Equal(t, http.StatusForbidden, recorder.Code)
	require.Equal(t, 1, roomService.transferCalls)
	require.Equal(t, "user_1", roomService.transferParams.OwnerID)
	require.Equal(t, "user_2", roomService.transferParams.NewOwnerID)
}

func TestChangeRoleRejectsInvalidRole(t *testing.T) {
	roomService := &fakeRoomHandlerService{}
	router := newRoomHandlerTestRouter(roomService)

	recorder := performRoomRequest(router, http.MethodPut, "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role", `{"role":"keeper"}`)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, roomService.changeRoleCalls)
}

func TestRoomHandlerErrorMappings(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		method     string
		path       string
		body       string
	}{
		{name: "get room not found", err: room.ErrRoomNotFound, wantStatus: http.StatusNotFound, method: http.MethodGet, path: "/api/rooms/11111111-1111-1111-1111-111111111111/"},
		{name: "join room full", err: room.ErrRoomFull, wantStatus: http.StatusConflict, method: http.MethodPost, path: "/api/rooms/11111111-1111-1111-1111-111111111111/join", body: `{"invite_token":"token"}`},
		{name: "join already member", err: room.ErrAlreadyMember, wantStatus: http.StatusConflict, method: http.MethodPost, path: "/api/rooms/11111111-1111-1111-1111-111111111111/join", body: `{"invite_token":"token"}`},
		{name: "leave not member", err: room.ErrNotMember, wantStatus: http.StatusNotFound, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/leave"},
		{name: "kick not owner", err: room.ErrNotOwner, wantStatus: http.StatusForbidden, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/kick/user_2"},
		{name: "kick owner self", err: room.ErrCannotKickOwner, wantStatus: http.StatusForbidden, method: http.MethodDelete, path: "/api/rooms/11111111-1111-1111-1111-111111111111/kick/user_2"},
		{name: "select character not owned", err: room.ErrCharacterNotOwned, wantStatus: http.StatusForbidden, method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/character", body: `{"character_id":"33333333-3333-3333-3333-333333333333"}`},
		{name: "change role generic", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, method: http.MethodPut, path: "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role", body: `{"role":"gm"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomService := &fakeRoomHandlerService{err: tt.err}
			router := newRoomHandlerTestRouter(roomService)

			recorder := performRoomRequest(router, tt.method, tt.path, tt.body)

			require.Equal(t, tt.wantStatus, recorder.Code)
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

	room   model.RoomModel
	member model.RoomMemberModel

	createCalls  int
	createParams db.CreateRoomParams

	getCalls  int
	getParams db.GetRoomByIDParams

	updateCalls  int
	updateParams db.UpdateRoomParams

	transferCalls  int
	transferParams db.TransferRoomOwnershipParams

	deleteCalls  int
	deleteParams db.DeleteRoomParams

	joinCalls  int
	joinMeta   db.GetRoomMetaDataParams
	joinMember db.GetMemberParams

	leaveCalls  int
	leaveParams db.RemoveMemberParams

	kickCalls  int
	kickActor  db.GetRoomByIDParams
	kickTarget db.RemoveMemberParams

	selectCharacterCalls  int
	selectCharacterParams db.UpdateMemberCharacterParams

	changeRoleCalls  int
	changeRoleActor  db.GetRoomByIDParams
	changeRoleTarget db.UpdateMemberRoleParams
}

func (f *fakeRoomHandlerService) totalCalls() int {
	return f.createCalls + f.getCalls + f.updateCalls + f.transferCalls + f.deleteCalls + f.joinCalls + f.leaveCalls + f.kickCalls + f.selectCharacterCalls + f.changeRoleCalls
}

func (f *fakeRoomHandlerService) CreateRoom(_ context.Context, params db.CreateRoomParams) (model.RoomModel, error) {
	f.createCalls++
	f.createParams = params
	return f.room, f.err
}

func (f *fakeRoomHandlerService) GetRoom(_ context.Context, params db.GetRoomByIDParams) (model.RoomModel, error) {
	f.getCalls++
	f.getParams = params
	return f.room, f.err
}

func (f *fakeRoomHandlerService) UpdateRoom(_ context.Context, params db.UpdateRoomParams) (model.RoomModel, error) {
	f.updateCalls++
	f.updateParams = params
	return f.room, f.err
}

func (f *fakeRoomHandlerService) TransferOwnership(_ context.Context, params db.TransferRoomOwnershipParams) (model.RoomModel, error) {
	f.transferCalls++
	f.transferParams = params
	return f.room, f.err
}

func (f *fakeRoomHandlerService) DeleteRoom(_ context.Context, params db.DeleteRoomParams) error {
	f.deleteCalls++
	f.deleteParams = params
	return f.err
}

func (f *fakeRoomHandlerService) JoinRoom(_ context.Context, meta db.GetRoomMetaDataParams, member db.GetMemberParams) (model.RoomMemberModel, error) {
	f.joinCalls++
	f.joinMeta = meta
	f.joinMember = member
	return f.member, f.err
}

func (f *fakeRoomHandlerService) LeaveRoom(_ context.Context, params db.RemoveMemberParams) error {
	f.leaveCalls++
	f.leaveParams = params
	return f.err
}

func (f *fakeRoomHandlerService) KickMember(_ context.Context, actor db.GetRoomByIDParams, target db.RemoveMemberParams) error {
	f.kickCalls++
	f.kickActor = actor
	f.kickTarget = target
	return f.err
}

func (f *fakeRoomHandlerService) SelectCharacter(_ context.Context, params db.UpdateMemberCharacterParams) (model.RoomMemberModel, error) {
	f.selectCharacterCalls++
	f.selectCharacterParams = params
	return f.member, f.err
}

func (f *fakeRoomHandlerService) ChangeRole(_ context.Context, actor db.GetRoomByIDParams, target db.UpdateMemberRoleParams) (model.RoomMemberModel, error) {
	f.changeRoleCalls++
	f.changeRoleActor = actor
	f.changeRoleTarget = target
	return f.member, f.err
}

func testRoomUnitUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
