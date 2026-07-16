package character

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *CharacterService) syncDynamicSkillBases(ctx context.Context, userID string, characterID pgtype.UUID, dexterity, education *int16) error {
	return s.repos.Queries.SyncDynamicSkillBases(ctx, db.SyncDynamicSkillBasesParams{
		UserID:      userID,
		CharacterID: characterID,
		Dexterity:   dexterity,
		Education:   education,
	})
}
