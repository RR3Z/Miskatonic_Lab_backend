package tests

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestE2ECharacterPortraitUploadReplaceAndDeleteLifecycle(t *testing.T) {
	subject := newE2ESubject(t)
	characterID := subject.createCharacter(t, "E2E Portrait Investigator")
	t.Cleanup(func() { subject.deleteCharacter(t, characterID) })

	firstContent := e2ePortraitPNG(t, 1)
	var first e2eCharacterShortResponse
	subject.doMultipartFile(t, http.MethodPatch, "/api/characters/"+characterID+"/", "portrait", "first.png", firstContent, http.StatusOK, &first)
	require.Equal(t, characterID, first.ID)
	require.NotNil(t, first.PortraitURL)

	body, headers := subject.doPublicPortraitRequest(t, http.MethodGet, *first.PortraitURL, http.StatusOK)
	require.Equal(t, firstContent, body)
	require.Equal(t, "image/png", headers.Get("Content-Type"))
	require.Equal(t, "public, max-age=31536000, immutable", headers.Get("Cache-Control"))
	require.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"))
	headBody, headHeaders := subject.doPublicPortraitRequest(t, http.MethodHead, *first.PortraitURL, http.StatusOK)
	require.Empty(t, headBody)
	require.Equal(t, "image/png", headHeaders.Get("Content-Type"))

	var full e2eCharacterResponse
	subject.doJSON(t, http.MethodGet, "/api/characters/"+characterID+"/", nil, http.StatusOK, &full)
	require.Equal(t, first.PortraitURL, full.PortraitURL)
	var summaries []e2eCharacterSummaryResponse
	subject.doJSON(t, http.MethodGet, "/api/characters/", nil, http.StatusOK, &summaries)
	require.Equal(t, first.PortraitURL, portraitURLForCharacter(summaries, characterID))

	secondContent := e2ePortraitPNG(t, 2)
	var second e2eCharacterShortResponse
	subject.doMultipartFile(t, http.MethodPatch, "/api/characters/"+characterID+"/", "portrait", "second.png", secondContent, http.StatusOK, &second)
	require.NotNil(t, second.PortraitURL)
	require.NotEqual(t, *first.PortraitURL, *second.PortraitURL)
	subject.doPublicPortraitRequest(t, http.MethodGet, *first.PortraitURL, http.StatusNotFound)
	secondBody, _ := subject.doPublicPortraitRequest(t, http.MethodGet, *second.PortraitURL, http.StatusOK)
	require.Equal(t, secondContent, secondBody)

	subject.doJSON(t, http.MethodDelete, "/api/characters/"+characterID+"/", nil, http.StatusNoContent, nil)
	subject.doPublicPortraitRequest(t, http.MethodGet, *second.PortraitURL, http.StatusNotFound)
}

func TestE2ECharacterPortraitRejectsInvalidMissingOversizedAndForeignUploads(t *testing.T) {
	subject := newE2ESubject(t)
	characterID := subject.createCharacter(t, "E2E Invalid Portrait Investigator")
	t.Cleanup(func() { subject.deleteCharacter(t, characterID) })

	cases := []struct {
		name       string
		field      string
		filename   string
		content    []byte
		wantStatus int
		wantCode   string
	}{
		{name: "unsupported", field: "portrait", filename: "portrait.txt", content: []byte("not an image"), wantStatus: http.StatusBadRequest, wantCode: "character.portrait_unsupported"},
		{name: "corrupt", field: "portrait", filename: "portrait.png", content: []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, wantStatus: http.StatusBadRequest, wantCode: "character.portrait_invalid"},
		{name: "missing", field: "avatar", filename: "", content: []byte("not a portrait field"), wantStatus: http.StatusBadRequest, wantCode: "character.portrait_required"},
		{name: "oversized", field: "portrait", filename: "portrait.png", content: make([]byte, (5<<20)+1), wantStatus: http.StatusRequestEntityTooLarge, wantCode: "character.portrait_too_large"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var response e2eErrorResponse
			subject.doMultipartFile(t, http.MethodPatch, "/api/characters/"+characterID+"/", tc.field, tc.filename, tc.content, tc.wantStatus, &response)
			require.Equal(t, tc.wantCode, response.Code)
		})
	}

	otherUserID := "e2e_portrait_other_" + e2eHash(time.Now().String())
	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{ID: otherUserID, Username: otherUserID, Email: otherUserID + "@example.com"})
	require.NoError(t, err)
	t.Cleanup(func() { _ = subject.queries.DeleteUserByClerkID(context.Background(), otherUserID) })
	otherCharacter, err := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{UserID: otherUserID, Name: "Foreign Portrait Investigator"})
	require.NoError(t, err)
	var foreignResponse e2eErrorResponse
	subject.doMultipartFile(t, http.MethodPatch, "/api/characters/"+otherCharacter.ID.String()+"/", "portrait", "portrait.png", e2ePortraitPNG(t, 3), http.StatusNotFound, &foreignResponse)
	require.Equal(t, "character.not_found", foreignResponse.Code)
}

func TestE2ECharacterLimitReturnsConflict(t *testing.T) {
	subject := newE2ESubject(t)
	existing, err := subject.queries.GetAllUserCharacters(context.Background(), subject.userID)
	require.NoError(t, err)
	seedCapacity := int(characterService.MaxCharactersPerUser) - len(existing)
	if seedCapacity < 0 {
		seedCapacity = 0
	}
	seeded := make([]pgtype.UUID, 0, seedCapacity)
	for i := len(existing); i < int(characterService.MaxCharactersPerUser); i++ {
		character, createErr := subject.queries.CreateCharacter(context.Background(), db.CreateCharacterParams{UserID: subject.userID, Name: "E2E Limit Seed"})
		require.NoError(t, createErr)
		seeded = append(seeded, character.ID)
	}
	t.Cleanup(func() {
		for _, id := range seeded {
			_, _ = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{ID: id, UserID: subject.userID})
		}
	})

	var response e2eErrorResponse
	subject.doJSON(t, http.MethodPost, "/api/characters/", map[string]any{"name": "Character 31"}, http.StatusConflict, &response)
	require.Equal(t, "character.limit_reached", response.Code)
}

func e2ePortraitPNG(t *testing.T, marker byte) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: marker, A: 255})
	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img))
	return buffer.Bytes()
}

func portraitURLForCharacter(characters []e2eCharacterSummaryResponse, characterID string) *string {
	for _, character := range characters {
		if character.ID == characterID {
			return character.PortraitURL
		}
	}
	return nil
}
