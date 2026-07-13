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
	service, router := newCharacterHandlerTestSubject(nil)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("portrait", "portrait.png")
	require.NoError(t, err)
	_, err = part.Write(make([]byte, characterHandler.MaxPortraitUploadBytes+(1<<20)))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusRequestEntityTooLarge, "character.portrait_too_large")
	require.Equal(t, 1, service.replacePortraitCalls)
}

func TestReplacePortraitRejectsInvalidCharacterIDBeforeService(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)
	body, contentType := portraitMultipartBody(t, "portrait", "portrait.png", []byte("portrait"))
	request := httptest.NewRequest(http.MethodPatch, "/api/characters/not-a-uuid/", body)
	request.Header.Set("Content-Type", contentType)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.invalid_id")
	require.Zero(t, service.replacePortraitCalls)
}

func TestReplacePortraitRejectsMalformedMultipartBody(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)
	request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", bytes.NewBufferString("--broken\r\ninvalid"))
	request.Header.Set("Content-Type", "multipart/form-data; boundary=broken")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	requireCharacterError(t, recorder, http.StatusBadRequest, "character.invalid_input")
	require.Zero(t, service.replacePortraitCalls)
}

func TestReplacePortraitIgnoresNonPortraitPartsAndRequiresNamedFile(t *testing.T) {
	cases := []struct {
		name     string
		field    string
		filename string
	}{
		{name: "wrong field", field: "avatar", filename: "portrait.png"},
		{name: "empty filename", field: "portrait", filename: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service, router := newCharacterHandlerTestSubject(nil)
			body, contentType := portraitMultipartBody(t, tc.field, tc.filename, []byte("portrait"))
			request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
			request.Header.Set("Content-Type", contentType)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			requireCharacterError(t, recorder, http.StatusBadRequest, "character.portrait_required")
			require.Zero(t, service.replacePortraitCalls)
		})
	}
}

func TestReplacePortraitMapsStorageErrorsToExactHTTPContracts(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "too large", err: characterServiceErrors.ErrPortraitTooLarge, wantStatus: http.StatusRequestEntityTooLarge, wantCode: "character.portrait_too_large"},
		{name: "unsupported", err: characterServiceErrors.ErrPortraitUnsupported, wantStatus: http.StatusBadRequest, wantCode: "character.portrait_unsupported"},
		{name: "invalid", err: characterServiceErrors.ErrPortraitInvalid, wantStatus: http.StatusBadRequest, wantCode: "character.portrait_invalid"},
		{name: "storage unavailable", err: characterServiceErrors.ErrPortraitStorage, wantStatus: http.StatusServiceUnavailable, wantCode: "character.portrait_storage_unavailable"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service, router := newCharacterHandlerTestSubject(tc.err)
			body, contentType := portraitMultipartBody(t, "portrait", "portrait.png", []byte("portrait"))
			request := httptest.NewRequest(http.MethodPatch, "/api/characters/"+testCharacterID+"/", body)
			request.Header.Set("Content-Type", contentType)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			requireCharacterError(t, recorder, tc.wantStatus, tc.wantCode)
			require.Equal(t, 1, service.replacePortraitCalls)
		})
	}
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

func portraitMultipartBody(t *testing.T, field string, filename string, content []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	var partWriter interface {
		Write([]byte) (int, error)
	}
	var err error
	if filename == "" {
		partWriter, err = writer.CreateFormField(field)
	} else {
		partWriter, err = writer.CreateFormFile(field, filename)
	}
	require.NoError(t, err)
	_, err = partWriter.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.FormDataContentType()
}
