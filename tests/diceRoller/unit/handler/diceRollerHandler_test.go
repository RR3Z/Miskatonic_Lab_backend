package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/stretchr/testify/require"
)

func TestDiceRollerRoutesAreMounted(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d20"}`)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, 1, svc.makeCalls)
}

func TestDiceRollerReturnsDetailsAsJSONObject(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{
		roll: diceRollerDTO.DiceRollModel{
			Expression: "1d100",
			Result:     24,
			Details:    json.RawMessage(`{"mode":"bonus","units":4,"tens":[2,4],"candidates":[24,44],"selected":24}`),
		},
	}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d100","d100_mode":"bonus"}`)

	require.Equal(t, http.StatusCreated, rec.Code)
	var response map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.NotContains(t, response, "d100")
	require.Equal(t, "bonus", response["details"].(map[string]any)["mode"])
}

func TestDiceRollerRejectsInvalidCharacterID(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/not-a-uuid/",
		`{"expression":"1d20"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Zero(t, svc.makeCalls)
	require.Contains(t, rec.Body.String(), "dice.invalid_character_id")
}

func TestDiceRollerRejectsInvalidBody(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "dice.invalid_input")
	require.Zero(t, svc.makeCalls)
}

func TestDiceRollerPassesDTOToService(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"2d6+1d4"}`)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, 1, svc.makeCalls)
	require.Equal(t, "user_test_01", svc.makeInput.UserID)
	require.Equal(t, "2d6+1d4", svc.makeInput.Formula)
	require.True(t, svc.makeInput.CharacterID.Valid)
	require.Nil(t, svc.makeInput.RoomID)
}

func TestDiceRollerPassesD100ModeToService(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d100","d100_mode":"bonus"}`)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, svc.makeInput.D100Mode)
	require.Equal(t, "bonus", string(*svc.makeInput.D100Mode))
}

func TestDiceRollerPassesRoomIDToService(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	checker := &fakeRoomAccessChecker{}
	router := newDiceRollerTestRouterWithRoom(svc, checker)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"2d6","room_id":"22222222-2222-2222-2222-222222222222"}`)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, 1, svc.makeCalls)
	require.NotNil(t, svc.makeInput.RoomID)
	require.Equal(t, "22222222-2222-2222-2222-222222222222", svc.makeInput.RoomID.String())
}

func TestDiceRollerMapsInvalidExpressionError(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{
		err: diceRollerServices.ErrInvalidExpression,
	}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"invalid"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "dice.invalid_expression")
}

func TestDiceRollerMapsCharacterNotFoundError(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{
		err: diceRollerServices.ErrCharacterNotFound,
	}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d20"}`)
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Contains(t, rec.Body.String(), "dice.character_not_found")
}

func TestDiceRollerListRollsPassesDTOToService(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodGet,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/lasts", "")
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 1, svc.listCalls)
	require.Equal(t, "user_test_01", svc.listInput.UserID)
	require.True(t, svc.listInput.CharacterID.Valid)
}

func TestDiceRollerListRollsMapsServiceError(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{
		err: errors.New("database error"),
	}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodGet,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/lasts", "")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDiceRollerGetLastsRejectsInvalidCharacterID(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodGet,
		"/api/dice-roll/not-a-uuid/lasts", "")
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Zero(t, svc.listCalls)
	require.Contains(t, rec.Body.String(), "dice.invalid_character_id")
}

func TestDiceRollerNoRoomCheckerWithRoomIDReturnsRoomNotAvailable(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d20","room_id":"22222222-2222-2222-2222-222222222222"}`)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Zero(t, svc.makeCalls)
	require.Contains(t, rec.Body.String(), "dice.room_not_available")
}

func TestDiceRollerRoomPreflightFailureReturnsRoomNotAvailable(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	checker := &fakeRoomAccessChecker{err: errors.New("not a member")}
	router := newDiceRollerTestRouterWithRoom(svc, checker)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d20","room_id":"22222222-2222-2222-2222-222222222222"}`)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Zero(t, svc.makeCalls)
	require.Contains(t, rec.Body.String(), "dice.room_not_available")
}
