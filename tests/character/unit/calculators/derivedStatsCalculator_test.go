package tests

import (
	"testing"

	calculators "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/calculators"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCalculateDerivedStatsPreservesTargetIdentity(t *testing.T) {
	characterID := testUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	stats := calculators.CalculateDerivedStats("user_id", characterID, testCharacteristics(50, 50, 50))

	require.Equal(t, "user_id", stats.UserID)
	require.Equal(t, characterID, stats.CharacterID)
}

func TestCalculateDerivedStatsCalculatesBaseSpeed(t *testing.T) {
	tests := []struct {
		name          string
		strength      int16
		dexterity     int16
		size          int16
		expectedSpeed int16
	}{
		{
			name:          "strength and dexterity below size",
			strength:      40,
			dexterity:     45,
			size:          50,
			expectedSpeed: 7,
		},
		{
			name:          "strength equals size",
			strength:      50,
			dexterity:     45,
			size:          50,
			expectedSpeed: 8,
		},
		{
			name:          "dexterity equals size",
			strength:      45,
			dexterity:     50,
			size:          50,
			expectedSpeed: 8,
		},
		{
			name:          "all three equal",
			strength:      50,
			dexterity:     50,
			size:          50,
			expectedSpeed: 8,
		},
		{
			name:          "strength and dexterity above size",
			strength:      60,
			dexterity:     55,
			size:          50,
			expectedSpeed: 9,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stats := calculators.CalculateDerivedStats("user_id", pgtype.UUID{}, testCharacteristics(tc.strength, tc.size, tc.dexterity))

			requireInt16PointerValue(t, stats.Speed, tc.expectedSpeed)
		})
	}
}

func TestCalculateDerivedStatsCalculatesDodgeAsHalfDexterityRoundedDown(t *testing.T) {
	tests := []struct {
		name          string
		dexterity     int16
		expectedDodge int16
	}{
		{name: "even dexterity", dexterity: 60, expectedDodge: 30},
		{name: "odd dexterity", dexterity: 55, expectedDodge: 27},
		{name: "zero dexterity", dexterity: 0, expectedDodge: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stats := calculators.CalculateDerivedStats("user_id", pgtype.UUID{}, testCharacteristics(50, 50, tc.dexterity))

			requireInt16PointerValue(t, stats.DodgeValue, tc.expectedDodge)
		})
	}
}

func TestCalculateDerivedStatsCalculatesPhysiqueAndDamageBonusBoundaries(t *testing.T) {
	tests := []struct {
		name                string
		strengthPlusSize    int16
		expectedPhysique    int16
		expectedDamageBonus string
	}{
		{name: "minimum flat penalty", strengthPlusSize: 2, expectedPhysique: -2, expectedDamageBonus: "-2"},
		{name: "64 upper penalty boundary", strengthPlusSize: 64, expectedPhysique: -2, expectedDamageBonus: "-2"},
		{name: "65 lower penalty boundary", strengthPlusSize: 65, expectedPhysique: -1, expectedDamageBonus: "-1"},
		{name: "84 upper penalty boundary", strengthPlusSize: 84, expectedPhysique: -1, expectedDamageBonus: "-1"},
		{name: "85 lower zero boundary", strengthPlusSize: 85, expectedPhysique: 0, expectedDamageBonus: "0"},
		{name: "124 upper zero boundary", strengthPlusSize: 124, expectedPhysique: 0, expectedDamageBonus: "0"},
		{name: "125 lower d4 boundary", strengthPlusSize: 125, expectedPhysique: 1, expectedDamageBonus: "+1d4"},
		{name: "164 upper d4 boundary", strengthPlusSize: 164, expectedPhysique: 1, expectedDamageBonus: "+1d4"},
		{name: "165 lower d6 boundary", strengthPlusSize: 165, expectedPhysique: 2, expectedDamageBonus: "+1d6"},
		{name: "204 upper d6 boundary", strengthPlusSize: 204, expectedPhysique: 2, expectedDamageBonus: "+1d6"},
		{name: "205 lower 2d6 boundary", strengthPlusSize: 205, expectedPhysique: 3, expectedDamageBonus: "+2d6"},
		{name: "284 upper 2d6 boundary", strengthPlusSize: 284, expectedPhysique: 3, expectedDamageBonus: "+2d6"},
		{name: "285 lower 3d6 boundary", strengthPlusSize: 285, expectedPhysique: 4, expectedDamageBonus: "+3d6"},
		{name: "364 upper 3d6 boundary", strengthPlusSize: 364, expectedPhysique: 4, expectedDamageBonus: "+3d6"},
		{name: "365 lower 4d6 boundary", strengthPlusSize: 365, expectedPhysique: 5, expectedDamageBonus: "+4d6"},
		{name: "444 upper 4d6 boundary", strengthPlusSize: 444, expectedPhysique: 5, expectedDamageBonus: "+4d6"},
		{name: "445 keeps increasing", strengthPlusSize: 445, expectedPhysique: 6, expectedDamageBonus: "+5d6"},
		{name: "525 keeps increasing again", strengthPlusSize: 525, expectedPhysique: 7, expectedDamageBonus: "+6d6"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			strength, size := splitTotal(tc.strengthPlusSize)
			stats := calculators.CalculateDerivedStats("user_id", pgtype.UUID{}, testCharacteristics(strength, size, 50))

			requireInt16PointerValue(t, stats.Physique, tc.expectedPhysique)
			requireStringPointerValue(t, stats.DamageBonus, tc.expectedDamageBonus)
		})
	}
}
