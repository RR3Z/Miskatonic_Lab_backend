package tests

import (
	"context"
	"testing"

	characterModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateCharacterRejectsBlankNameBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{})

	_, err := service.CreateCharacter(context.Background(), characterModel.CreateCharacterInput{
		UserID: "user_1",
		Name:   "   ",
	})

	require.ErrorIs(t, err, characterErrors.ErrNameRequired)
}

func TestUpdateCharacterRejectsBlankNameBeforeRepository(t *testing.T) {
	service := characterServices.NewCharacterService(&repository.Repository{})

	_, err := service.UpdateCharacter(context.Background(), characterModel.UpdateCharacterInput{
		UserID: "user_1",
		ID:     testCoreUUID("11111111-1111-1111-1111-111111111111"),
		Name:   "",
	})

	require.ErrorIs(t, err, characterErrors.ErrNameRequired)
}

func testCoreUUID(value string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		panic(err)
	}
	return id
}
