package character

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
)

type PortraitStore interface {
	Save(ctx context.Context, file io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}

func (s *CharacterService) ReplacePortrait(ctx context.Context, input characterDTO.ReplacePortraitInput) (characterDTO.CharacterShortModel, error) {
	if input.File == nil {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrPortraitRequired
	}
	if s.portraitStore == nil {
		return characterDTO.CharacterShortModel{}, characterErrors.ErrPortraitStorage
	}

	if _, err := s.repos.Queries.GetCharacter(ctx, db.GetCharacterParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	}); err != nil {
		return characterDTO.CharacterShortModel{}, err
	}

	portraitKey, err := s.portraitStore.Save(ctx, input.File)
	if err != nil {
		return characterDTO.CharacterShortModel{}, mapPortraitStoreError(err)
	}
	cleanupNewPortrait := func() {
		s.removePortraitFile(
			context.WithoutCancel(ctx),
			portraitKey,
			"failed to remove portrait after database replacement failure",
		)
	}

	tx, err := s.repos.DB.Begin(ctx)
	if err != nil {
		cleanupNewPortrait()
		return characterDTO.CharacterShortModel{}, err
	}
	defer tx.Rollback(ctx)

	queries := s.repos.Queries.WithTx(tx)
	existing, err := queries.LockCharacterForPortraitReplacement(ctx, db.LockCharacterForPortraitReplacementParams{
		UserID: input.UserID,
		ID:     input.CharacterID,
	})
	if err != nil {
		cleanupNewPortrait()
		return characterDTO.CharacterShortModel{}, err
	}

	character, err := queries.SetCharacterPortraitKey(ctx, db.SetCharacterPortraitKeyParams{
		UserID:      input.UserID,
		ID:          input.CharacterID,
		PortraitKey: &portraitKey,
	})
	if err != nil {
		cleanupNewPortrait()
		return characterDTO.CharacterShortModel{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		cleanupNewPortrait()
		return characterDTO.CharacterShortModel{}, err
	}

	if existing.PortraitKey != nil && *existing.PortraitKey != portraitKey {
		s.removePortraitFile(
			context.WithoutCancel(ctx),
			*existing.PortraitKey,
			"failed to remove replaced character portrait",
			"character_id", input.CharacterID.String(),
		)
	}

	result := characterDTO.ToCharacterShortModel(character)
	result.PortraitUrl = s.portraitURL(character.PortraitKey)
	return result, nil
}

func (s *CharacterService) portraitURL(key *string) *string {
	if key == nil || s.portraitStore == nil {
		return nil
	}
	publicURL := s.portraitStore.PublicURL(*key)
	if publicURL == "" {
		return nil
	}
	return &publicURL
}

func (s *CharacterService) removePortraitFile(ctx context.Context, key string, message string, attributes ...any) {
	if s.portraitStore == nil {
		return
	}
	if err := s.portraitStore.Delete(ctx, key); err != nil {
		attributes = append(attributes, "portrait_key", key, "error", err)
		slog.Warn(message, attributes...)
	}
}

func mapPortraitStoreError(err error) error {
	switch {
	case errors.Is(err, portraitStorage.ErrPortraitRequired):
		return characterErrors.ErrPortraitRequired
	case errors.Is(err, portraitStorage.ErrPortraitTooLarge):
		return characterErrors.ErrPortraitTooLarge
	case errors.Is(err, portraitStorage.ErrUnsupportedImage):
		return characterErrors.ErrPortraitUnsupported
	case errors.Is(err, portraitStorage.ErrInvalidImage):
		return characterErrors.ErrPortraitInvalid
	default:
		return fmt.Errorf("%w: %v", characterErrors.ErrPortraitStorage, err)
	}
}
