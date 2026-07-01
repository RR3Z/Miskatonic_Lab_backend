package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	diceRollerDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/diceRoller"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	diceRollerServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/diceRoller"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

type fakeDiceRollerHandlerService struct {
	roll      diceRollerDTO.DiceRollModel
	rolls     []diceRollerDTO.DiceRollModel
	err       error
	makeCalls int
	makeInput diceRollerDTO.MakeRollInput
	listCalls int
	listInput diceRollerDTO.GetLastDiceRollsInput
}

func (f *fakeDiceRollerHandlerService) MakeRoll(_ context.Context, input diceRollerDTO.MakeRollInput) (diceRollerDTO.DiceRollModel, error) {
	f.makeCalls++
	f.makeInput = input
	return f.roll, f.err
}

func (f *fakeDiceRollerHandlerService) GetLastDiceRolls(_ context.Context, input diceRollerDTO.GetLastDiceRollsInput) ([]diceRollerDTO.DiceRollModel, error) {
	f.listCalls++
	f.listInput = input
	return f.rolls, f.err
}

func newDiceRollerTestRouter(svc *fakeDiceRollerHandlerService) http.Handler {
	services := &service.Service{
		DiceRoller: svc,
	}
	h := handler.NewHandler(services)
	return h.InitRoutesWithAuth(diceRollerAuthMiddleware)
}

func diceRollerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := clerk.ContextWithSessionClaims(r.Context(), &clerk.SessionClaims{
			RegisteredClaims: clerk.RegisteredClaims{Subject: "user_test_01"},
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func performDiceRollerRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestDiceRollerRoutesAreMounted(t *testing.T) {
	svc := &fakeDiceRollerHandlerService{}
	router := newDiceRollerTestRouter(svc)

	rec := performDiceRollerRequest(router, http.MethodPost,
		"/api/dice-roll/11111111-1111-1111-1111-111111111111/",
		`{"expression":"1d20"}`)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, 1, svc.makeCalls)
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
	require.Contains(t, rec.Body.String(), "common.invalid_request")
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
