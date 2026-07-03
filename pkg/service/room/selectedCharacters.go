package room

import (
	"context"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	roomHelpers "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room/helpers"
)

func (s *RoomService) ListSelectedCharacters(ctx context.Context, input model.ListSelectedCharactersInput) ([]model.SelectedCharacterModel, error) {
	if err := s.EnsureMember(ctx, input.RoomID, input.UserID); err != nil {
		return nil, err
	}

	members, err := s.repos.Queries.ListVisibleSelectedRoomMembers(ctx, db.ListVisibleSelectedRoomMembersParams{
		RoomID: input.RoomID,
		UserID: input.UserID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.SelectedCharacterModel, 0, len(members))
	for _, member := range members {
		if !member.CharacterID.Valid {
			continue
		}

		character, err := roomHelpers.LoadCharacterSheet(ctx, s.repos.Queries, member.UserID, member.CharacterID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.ToSelectedCharacterModel(member, character))
	}

	return result, nil
}
