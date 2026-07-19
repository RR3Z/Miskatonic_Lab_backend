package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	roomModels "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	ws "github.com/RR3Z/Miskatonic_Lab_backend/pkg/ws/room"
	"github.com/stretchr/testify/require"
)

func TestRoomMutationsBroadcastRealtimeEvents(t *testing.T) {
	roomID := testRoomUnitUUID("11111111-1111-1111-1111-111111111111")
	tests := []struct {
		name      string
		method    string
		path      string
		body      string
		service   *fakeRoomHandlerService
		wantType  roomModels.EventType
		wantActor string
	}{
		{
			name:      "room update",
			method:    http.MethodPut,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/",
			body:      `{"max_players":5}`,
			service:   &fakeRoomHandlerService{},
			wantType:  roomModels.EventRoomUpdated,
			wantActor: "user_1",
		},
		{
			name:      "ownership transfer",
			method:    http.MethodPut,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/owner",
			body:      `{"user_id":"user_2"}`,
			service:   &fakeRoomHandlerService{},
			wantType:  roomModels.EventOwnerTransferred,
			wantActor: "user_1",
		},
		{
			name:      "member join",
			method:    http.MethodPost,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/join",
			body:      `{"invite_token":"token"}`,
			service:   &fakeRoomHandlerService{member: roomModels.RoomMemberModel{Role: "player"}},
			wantType:  roomModels.EventMemberJoined,
			wantActor: "user_1",
		},
		{
			name:      "member leave",
			method:    http.MethodDelete,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/leave",
			service:   &fakeRoomHandlerService{},
			wantType:  roomModels.EventMemberLeft,
			wantActor: "user_1",
		},
		{
			name:      "member kick",
			method:    http.MethodDelete,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/kick/user_2",
			service:   &fakeRoomHandlerService{},
			wantType:  roomModels.EventMemberKicked,
			wantActor: "user_1",
		},
		{
			name:      "character selected",
			method:    http.MethodPut,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/character",
			body:      `{"character_id":"22222222-2222-2222-2222-222222222222"}`,
			service:   &fakeRoomHandlerService{member: roomModels.RoomMemberModel{CharacterID: testRoomUnitUUID("22222222-2222-2222-2222-222222222222")}},
			wantType:  roomModels.EventMemberCharacterSelected,
			wantActor: "user_1",
		},
		{
			name:      "role changed",
			method:    http.MethodPut,
			path:      "/api/rooms/11111111-1111-1111-1111-111111111111/members/user_2/role",
			body:      `{"role":"gm"}`,
			service:   &fakeRoomHandlerService{member: roomModels.RoomMemberModel{Role: "gm"}},
			wantType:  roomModels.EventMemberRoleChanged,
			wantActor: "user_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service.mutationEvents = []roomModels.RoomEventModel{{
				ID:      testRoomUnitUUID("33333333-3333-3333-3333-333333333333"),
				RoomID:  roomID,
				ActorID: tt.wantActor,
				Type:    string(tt.wantType),
				Payload: []byte(`{"source":"saved"}`),
			}}
			hub := ws.NewRoomHub()
			events := registerRoomUnitTestClient(t, hub, roomID)
			router := newRoomHandlerTestRouterWithHub(tt.service, hub)

			recorder := performRoomRequest(router, tt.method, tt.path, tt.body)

			require.Less(t, recorder.Code, http.StatusBadRequest)
			select {
			case event := <-events:
				require.Equal(t, string(tt.wantType), event.Type)
				require.Equal(t, roomID.String(), event.RoomID)
				require.Equal(t, tt.wantActor, event.ActorID)
				payload, ok := event.Payload.(json.RawMessage)
				require.True(t, ok)
				require.JSONEq(t, `{"source":"saved"}`, string(payload))
			case <-time.After(time.Second):
				t.Fatal("room mutation did not broadcast a realtime event")
			}
		})
	}
}
