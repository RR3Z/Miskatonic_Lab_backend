package characterErrors

import "errors"

var (
	ErrNameRequired          = errors.New("character name is required")
	ErrNameTooLong           = errors.New("character name exceeds max length")
	ErrAgeNegative           = errors.New("character age must be >= 0")
	ErrSexInvalid            = errors.New("character sex must be male or female")
	ErrPatchInvalid          = errors.New("character patch input is invalid")
	ErrCharacterLimitReached = errors.New("character limit reached")
	ErrPortraitRequired      = errors.New("portrait is required")
	ErrPortraitTooLarge      = errors.New("portrait is too large")
	ErrPortraitUnsupported   = errors.New("portrait format is unsupported")
	ErrPortraitInvalid       = errors.New("portrait image is invalid")
	ErrPortraitStorage       = errors.New("portrait storage is unavailable")

	ErrCharacteristicsNegative = errors.New("characteristics values must be >= 0")

	ErrDerivedStatsNegative = errors.New("derived stats speed/dodge must be >= 0")
	ErrInvalidDamageBonus   = errors.New("invalid derived stats damage bonus")
	ErrInvalidDerivedStats  = errors.New("invalid derived stats")

	ErrStateNegative          = errors.New("state values must be >= 0")
	ErrStateCurrentExceedsMax = errors.New("current value cannot exceed max value")

	ErrInvalidBackstorySection = errors.New("invalid backstory item section")
	ErrSectionTooLong          = errors.New("backstory item section exceeds max length")
	ErrBackstoryTitleRequired  = errors.New("backstory item title is required")
	ErrBackstoryTitleTooLong   = errors.New("backstory item title exceeds max length")
	ErrBackstoryTextRequired   = errors.New("backstory item text is required")

	ErrSkillNameRequired  = errors.New("skill name is required")
	ErrSkillNameTooLong   = errors.New("skill name exceeds max length")
	ErrSkillValueNegative = errors.New("skill values must be >= 0")
	ErrInvalidSkill       = errors.New("invalid skill")
	ErrProtectedSkill     = errors.New("protected skill cannot be renamed, rebased, or deleted")

	ErrFinancesMoneyTooLong = errors.New("finances money field exceeds max length")

	ErrNoteTitleRequired = errors.New("note title is required")
	ErrNoteTitleTooLong  = errors.New("note title exceeds max length")
	ErrNoteBodyRequired  = errors.New("note body is required")

	ErrInventoryItemNameRequired    = errors.New("inventory item name is required")
	ErrInventoryItemNameTooLong     = errors.New("inventory item name exceeds max length")
	ErrInventoryItemQuantityInvalid = errors.New("inventory item quantity must be >= 1")
	ErrInventoryItemCategoryTooLong = errors.New("inventory item category exceeds max length")
)
