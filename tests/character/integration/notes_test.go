package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestNotesTableCreateListGetUpdateAndDeleteNote(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	firstNote, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Title:       "First session",
		Body:        "The investigator found a locked cabinet.",
	})
	require.NoError(t, err)

	require.True(t, firstNote.ID.Valid)
	require.Equal(t, character.ID, firstNote.CharacterID)
	require.Equal(t, "First session", firstNote.Title)
	require.Equal(t, "The investigator found a locked cabinet.", firstNote.Body)
	require.True(t, firstNote.CreatedAt.Valid)
	require.True(t, firstNote.UpdatedAt.Valid)

	time.Sleep(5 * time.Millisecond)

	secondNote, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Title:       "Second session",
		Body:        "The key was hidden in a library ledger.",
	})
	require.NoError(t, err)

	notes, err := subject.queries.GetNotes(context.Background(), db.GetNotesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Len(t, notes, 2)
	require.Equal(t, secondNote.ID, notes[0].ID)
	require.Equal(t, firstNote.ID, notes[1].ID)

	fetchedNote, err := subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      firstNote.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstNote.ID, fetchedNote.ID)

	updatedNote, err := subject.queries.UpdateNote(context.Background(), db.UpdateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      firstNote.ID,
		Title:       "Updated session",
		Body:        "The cabinet contained a silver key.",
	})
	require.NoError(t, err)
	require.Equal(t, firstNote.ID, updatedNote.ID)
	require.Equal(t, "Updated session", updatedNote.Title)
	require.Equal(t, "The cabinet contained a silver key.", updatedNote.Body)
	require.True(t, updatedNote.UpdatedAt.Time.After(firstNote.UpdatedAt.Time) || updatedNote.UpdatedAt.Time.Equal(firstNote.UpdatedAt.Time))

	deletedNote, err := subject.queries.DeleteNote(context.Background(), db.DeleteNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      firstNote.ID,
	})
	require.NoError(t, err)
	require.Equal(t, firstNote.ID, deletedNote.ID)

	_, err = subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      firstNote.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestNotesTableListReturnsEmptyForCharacterWithoutNotes(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	notes, err := subject.queries.GetNotes(context.Background(), db.GetNotesParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, notes)
}

func TestNotesTableRequiresCharacterOwnerForCreateListGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, owner.ID)

	_, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		Title:       "Unauthorized note",
		Body:        "This should not be inserted.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	note, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		Title:       "Owner note",
		Body:        "Only the owner can manage this note.",
	})
	require.NoError(t, err)

	notes, err := subject.queries.GetNotes(context.Background(), db.GetNotesParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Empty(t, notes)

	_, err = subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateNote(context.Background(), db.UpdateNoteParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
		Title:       "Unauthorized update",
		Body:        "This should not update.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteNote(context.Background(), db.DeleteNoteParams{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	deletedNote, err := subject.queries.DeleteNote(context.Background(), db.DeleteNoteParams{
		UserID:      owner.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
	})
	require.NoError(t, err)
	require.Equal(t, note.ID, deletedNote.ID)
}

func TestNotesTableListsOnlyNotesForRequestedCharacter(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	firstCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	secondCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	firstNote, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
		Title:       "First character note",
		Body:        "This belongs to the first character.",
	})
	require.NoError(t, err)

	secondNote, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
		Title:       "Second character note",
		Body:        "This belongs to the second character.",
	})
	require.NoError(t, err)

	firstCharacterNotes, err := subject.queries.GetNotes(context.Background(), db.GetNotesParams{
		UserID:      testUser.ID,
		CharacterID: firstCharacter.ID,
	})
	require.NoError(t, err)
	require.Len(t, firstCharacterNotes, 1)
	require.Equal(t, firstNote.ID, firstCharacterNotes[0].ID)

	secondCharacterNotes, err := subject.queries.GetNotes(context.Background(), db.GetNotesParams{
		UserID:      testUser.ID,
		CharacterID: secondCharacter.ID,
	})
	require.NoError(t, err)
	require.Len(t, secondCharacterNotes, 1)
	require.Equal(t, secondNote.ID, secondCharacterNotes[0].ID)
}

func TestNotesTableRequiresMatchingCharacterForGetUpdateAndDelete(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	owningCharacter := createCharacterTestCharacter(t, subject, testUser.ID)
	otherCharacter := createCharacterTestCharacter(t, subject, testUser.ID)

	note, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: owningCharacter.ID,
		Title:       "Character scoped note",
		Body:        "This note must not be reachable through another character.",
	})
	require.NoError(t, err)

	_, err = subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		NoteID:      note.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateNote(context.Background(), db.UpdateNoteParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		NoteID:      note.ID,
		Title:       "Wrong character update",
		Body:        "This should not update.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteNote(context.Background(), db.DeleteNoteParams{
		UserID:      testUser.ID,
		CharacterID: otherCharacter.ID,
		NoteID:      note.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	fetchedNote, err := subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: owningCharacter.ID,
		NoteID:      note.ID,
	})
	require.NoError(t, err)
	require.Equal(t, note.ID, fetchedNote.ID)
	require.Equal(t, "Character scoped note", fetchedNote.Title)
}

func TestNotesTableReturnsNoRowsForMissingCharacterOrNote(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)
	missingCharacterID := characterTestUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	missingNoteID := characterTestUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	_, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: missingCharacterID,
		Title:       "Missing character",
		Body:        "This should not be inserted.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      missingNoteID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.UpdateNote(context.Background(), db.UpdateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      missingNoteID,
		Title:       "Missing note",
		Body:        "This should not update.",
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = subject.queries.DeleteNote(context.Background(), db.DeleteNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      missingNoteID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestNotesTableAllowsEmptyTitleAndBody(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	note, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Title:       "",
		Body:        "",
	})
	require.NoError(t, err)
	require.Equal(t, "", note.Title)
	require.Equal(t, "", note.Body)

	updatedNote, err := subject.queries.UpdateNote(context.Background(), db.UpdateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
		Title:       "",
		Body:        "",
	})
	require.NoError(t, err)
	require.Equal(t, "", updatedNote.Title)
	require.Equal(t, "", updatedNote.Body)
}

func TestNotesTableRejectsMissingRequiredFields(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	_, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Title:       strings.Repeat("a", 121),
		Body:        "Body",
	})
	requirePostgresErrorCode(t, err, "22001")
}

func TestNotesTableDeletingCharacterCascadesNotes(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	testUser := createCharacterTestUser(t, subject)
	character := createCharacterTestCharacter(t, subject, testUser.ID)

	note, err := subject.queries.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		Title:       "Cascade note",
		Body:        "This note should be removed with its character.",
	})
	require.NoError(t, err)

	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: testUser.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)

	_, err = subject.queries.GetNote(context.Background(), db.GetNoteParams{
		UserID:      testUser.ID,
		CharacterID: character.ID,
		NoteID:      note.ID,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
