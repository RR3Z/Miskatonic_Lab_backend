package characterErrors

import (
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapCharacterConstraintError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.ConstraintName {
	case "chk_health_states_current_lte_max":
		return myErrors.ErrCurrentHealthExceedsMax
	case "chk_magic_states_current_lte_max":
		return myErrors.ErrCurrentMagicExceedsMax
	case "chk_sanity_states_current_lte_max":
		return myErrors.ErrCurrentSanityExceedsMax
	case "chk_luck_states_current_lte_starting":
		return myErrors.ErrCurrentLuckExceedsStarting
	case "fk_finances_credit_rating_skill":
		return ErrInvalidFinances
	case "chk_backstory_items_section":
		return ErrInvalidBackstorySection
	case "derived_stats_speed_check",
		"derived_stats_physique_check",
		"derived_stats_damage_bonus_check",
		"derived_stats_dodge_value_check",
		"chk_derived_stats_damage_bonus_format":
		return ErrInvalidDerivedStats
	case "skills_base_value_check",
		"skills_value_check",
		"fk_skills_category",
		"fk_skills_specialty":
		return ErrInvalidSkill
	default:
		return err
	}
}

func MapDeleteSkillError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.ConstraintName == "fk_finances_credit_rating_skill" {
		return ErrSkillInUse
	}

	return err
}
