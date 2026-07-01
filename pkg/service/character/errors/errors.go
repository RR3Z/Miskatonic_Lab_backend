package characterErrors

import "errors"

var (
	ErrNameRequired            = errors.New("character name is required")
	ErrInvalidBackstorySection = errors.New("invalid backstory item section")
	ErrInvalidDerivedStats     = errors.New("invalid derived stats")
	ErrInvalidFinances         = errors.New("invalid finances")
	ErrInvalidSkill            = errors.New("invalid skill")
	ErrSkillInUse              = errors.New("skill is referenced by another character resource")
)
