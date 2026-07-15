package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchCharacterProfilePassesOnlyProvidedFields(t *testing.T) {
	service, router := newCharacterHandlerTestSubject(nil)

	recorder := performCharacterRequest(
		router,
		http.MethodPatch,
		"/api/characters/"+testCharacterID+"/",
		`{"occupation":"Antiquarian","age":null}`,
	)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, service.patchCalls)
	require.Zero(t, service.replacePortraitCalls)
	require.Equal(t, "user_1", service.patchInput.UserID)
	require.Equal(t, testCharacterUnitUUID(testCharacterID), service.patchInput.ID)
	require.True(t, service.patchInput.Occupation.Set)
	require.Equal(t, "Antiquarian", *service.patchInput.Occupation.Value)
	require.True(t, service.patchInput.Age.Set)
	require.Nil(t, service.patchInput.Age.Value)
	require.False(t, service.patchInput.Name.Set)
	require.False(t, service.patchInput.PlayerName.Set)
	require.False(t, service.patchInput.Sex.Set)
	require.False(t, service.patchInput.Residence.Set)
	require.False(t, service.patchInput.Birthplace.Set)
}

func TestPatchCharacterProfileRejectsInvalidPayloadBeforeService(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{name: "empty object", body: `{}`},
		{name: "unknown field", body: `{"unknown":"value"}`},
		{name: "portrait url", body: `{"portrait_url":"https://example.test/portrait.png"}`},
		{name: "null required name", body: `{"name":null}`},
		{name: "invalid field type", body: `{"age":"forty-two"}`},
		{name: "trailing json", body: `{"age":42}{"sex":"male"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service, router := newCharacterHandlerTestSubject(nil)

			recorder := performCharacterRequest(
				router,
				http.MethodPatch,
				"/api/characters/"+testCharacterID+"/",
				tc.body,
			)

			requireCharacterError(t, recorder, http.StatusBadRequest, "character.invalid_input")
			require.Zero(t, service.totalCalls())
		})
	}
}
