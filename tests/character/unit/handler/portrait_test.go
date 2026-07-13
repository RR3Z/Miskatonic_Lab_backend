package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	characterHandler "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	characterServiceErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestReplacePortraitPassesMultipartFileAndOwnershipContext(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)
	png := append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, make([]byte, 32)...)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("portrait", "portrait.png")
	require.NoError(t, err)
	_, err = part.Write(png)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, service.replacePortraitCalls)
	require.Equal(t, "user_1", service.replacePortraitInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.replacePortraitInput.CharacterID)
	require.Equal(t, png, service.replacePortraitContent)
}

func TestReplacePortraitRequiresAuthentication(t *testing.T) {
	characterService := &fakeCharacterHandlerService{}
	router := handler.NewHandler(handler.Dependencies{Services: &service.Service{Character: characterService}}).InitRoutes(func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	})

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Zero(t, characterService.replacePortraitCalls)
}

func TestReplacePortraitHidesForeignCharacterAsNotFound(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(pgx.ErrNoRows)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("portrait", "portrait.png")
	require.NoError(t, err)
	_, err = part.Write([]byte("portrait"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusNotFound, "character.not_found")
	require.Equal(t, 1, service.replacePortraitCalls)
}

func TestReplacePortraitRequiresMultipartFile(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.portrait_required")
	require.Zero(t, service.replacePortraitCalls)
}

func TestReplacePortraitRejectsFileLargerThanFiveMiB(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(characterServiceErrors.ErrPortraitTooLarge)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("portrait", "portrait.png")
	require.NoError(t, err)
	_, err = part.Write(make([]byte, characterHandler.MaxPortraitUploadBytes+1))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusRequestEntityTooLarge, "character.portrait_too_large")
	require.Equal(t, 1, service.replacePortraitCalls)
}

func TestPortraitFileServerIsMountedForGetAndHeadOnly(t *testing.T) {
	characterService := &fakeCharacterHandlerService{}
	fileServerCalls := 0
	fileServer := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fileServerCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	router := handler.NewHandler(handler.Dependencies{
		Services:           &service.Service{Character: characterService},
		PortraitFileServer: fileServer,
	}).InitRoutes(func(next http.Handler) http.Handler { return next })

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(method, "/uploads/portraits/test.png", nil))
		require.Equal(t, http.StatusNoContent, recorder.Code)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/uploads/portraits/test.png", nil))
	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	require.Equal(t, 2, fileServerCalls)
}
