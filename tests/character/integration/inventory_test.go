package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestInventoryItemsTableCRUDOwnershipAndCascade(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)
	quantity := int32(2)
	category := "Tools"
	description := "Fresh batteries."

	item, err := subject.queries.CreateInventoryItem(context.Background(), db.CreateInventoryItemParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Name:        "Flashlight",
		Quantity:    &quantity,
		Category:    &category,
		Description: &description,
	})
	require.NoError(t, err)
	require.Equal(t, "Flashlight", item.Name)
	require.Equal(t, int32(2), *item.Quantity)
	require.Equal(t, "Tools", *item.Category)
	require.Equal(t, "Fresh batteries.", *item.Description)

	items, err := subject.queries.GetInventoryItems(context.Background(), db.GetInventoryItemsParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, item.ID, items[0].ID)

	updated, err := subject.queries.UpdateInventoryItem(context.Background(), db.UpdateInventoryItemParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		ItemID:      item.ID,
		Name:        "Pocket flashlight",
		Quantity:    nil,
		Category:    nil,
		Description: nil,
	})
	require.NoError(t, err)
	require.Equal(t, "Pocket flashlight", updated.Name)
	require.Nil(t, updated.Quantity)
	require.Nil(t, updated.Category)
	require.Nil(t, updated.Description)

	_, err = subject.queries.GetInventoryItem(context.Background(), db.GetInventoryItemParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		ItemID:      item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteInventoryItem(context.Background(), db.DeleteInventoryItemParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		ItemID:      item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: owner.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetInventoryItem(context.Background(), db.GetInventoryItemParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		ItemID:      item.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
