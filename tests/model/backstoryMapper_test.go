package tests

import (
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/model"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/stretchr/testify/require"
)

func TestToBackstoryItemModelCopiesAllFields(t *testing.T) {
	item := testBackstoryItem()

	result := model.ToBackstoryItemModel(item)

	require.Equal(t, item.ID, result.ID)
	require.Equal(t, item.Section, result.Section)
	require.Equal(t, item.Title, result.Title)
	require.Equal(t, item.Text, result.Text)
	require.Equal(t, item.CreatedAt.Time, result.CreatedAt.Time)
	require.Equal(t, item.UpdatedAt.Time, result.UpdatedAt.Time)
}

func TestToBackstoryModelCopiesBackstoryAndMapsItems(t *testing.T) {
	backstory := testBackstory()
	item := testBackstoryItem()

	result := model.ToBackstoryModel(backstory, []db.BackstoryItem{item})

	require.Equal(t, backstory.ID, result.ID)
	require.Equal(t, backstory.CharacterID, result.CharacterID)
	require.Equal(t, backstory.PersonalDescription, result.PersonalDescription)
	require.Equal(t, backstory.CreatedAt.Time, result.CreatedAt.Time)
	require.Equal(t, backstory.UpdatedAt.Time, result.UpdatedAt.Time)
	require.Len(t, result.Items, 1)
	require.Equal(t, item.ID, result.Items[0].ID)
	require.Equal(t, item.Section, result.Items[0].Section)
}

func TestToBackstoryModelPreservesNilDescriptionAndEmptyItems(t *testing.T) {
	backstory := testBackstory()
	backstory.PersonalDescription = nil

	result := model.ToBackstoryModel(backstory, nil)

	require.Nil(t, result.PersonalDescription)
	require.Empty(t, result.Items)
}
