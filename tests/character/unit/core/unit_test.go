package tests

import (
	"bytes"
	"context"
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateCharacterRejectsBlankNameBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)

	_, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: "user_1",
		Name:   "   ",
	})

	require.ErrorIs(t, err, characterErrors.ErrNameRequired)
}

func TestUpdateCharacterRejectsBlankNameBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)

	_, err := service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID: "user_1",
		ID:     testCoreUUID("11111111-1111-1111-1111-111111111111"),
		Name:   "",
	})

	require.ErrorIs(t, err, characterErrors.ErrNameRequired)
}

func TestCreateCharacterRejectsNameTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longName := string(make([]byte, 256))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}

	_, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: "user_1",
		Name:   longName,
	})
	require.ErrorIs(t, err, characterErrors.ErrNameTooLong)
}

func TestUpdateCharacterRejectsNameTooLong(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	longName := string(make([]byte, 256))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}

	_, err := service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID: "user_1",
		ID:     testCoreUUID("11111111-1111-1111-1111-111111111111"),
		Name:   longName,
	})
	require.ErrorIs(t, err, characterErrors.ErrNameTooLong)
}

func TestCreateCharacterRejectsNegativeAgeBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	negativeAge := int16(-1)

	_, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: "user_1",
		Name:   "Investigator",
		Age:    &negativeAge,
	})

	require.ErrorIs(t, err, characterErrors.ErrAgeNegative)
}

func TestUpdateCharacterRejectsNegativeAgeBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	negativeAge := int16(-1)

	_, err := service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID: "user_1",
		ID:     testCoreUUID("11111111-1111-1111-1111-111111111111"),
		Name:   "Investigator",
		Age:    &negativeAge,
	})

	require.ErrorIs(t, err, characterErrors.ErrAgeNegative)
}

func TestCreateCharacterRejectsInvalidSexBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	sex := "other"

	_, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: "user_1",
		Name:   "Investigator",
		Sex:    &sex,
	})

	require.ErrorIs(t, err, characterErrors.ErrSexInvalid)
}

func TestUpdateCharacterRejectsInvalidSexBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)
	sex := "other"

	_, err := service.UpdateCharacter(context.Background(), characterDTO.UpdateCharacterInput{
		UserID: "user_1",
		ID:     testCoreUUID("11111111-1111-1111-1111-111111111111"),
		Name:   "Investigator",
		Sex:    &sex,
	})

	require.ErrorIs(t, err, characterErrors.ErrSexInvalid)
}

func TestReplacePortraitRejectsMissingFileBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)

	_, err := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID:      "user_1",
		CharacterID: testCoreUUID("11111111-1111-1111-1111-111111111111"),
	})

	require.ErrorIs(t, err, characterErrors.ErrPortraitRequired)
}

func TestReplacePortraitRejectsUnavailableStorageBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{}, nil, nil)

	_, err := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID:      "user_1",
		CharacterID: testCoreUUID("11111111-1111-1111-1111-111111111111"),
		File:        bytes.NewReader([]byte("portrait")),
	})

	require.ErrorIs(t, err, characterErrors.ErrPortraitStorage)
}

func testCoreUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
