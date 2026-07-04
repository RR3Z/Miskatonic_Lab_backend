package tests

import (
	"context"

	characteristicsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/characteristics"
	derivedStatsDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/derivedstats"
	financesDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/finances"
	healthDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/health"
	luckDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/luck"
	magicDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/magic"
	sanityDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/sanity"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
)

func (f *fakeCharacterHandlerService) GetHealth(_ context.Context, input healthDTO.GetHealthInput) (db.HealthState, error) {
	f.getHealthCalls++
	f.getHealthInput = input
	return db.HealthState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertHealth(_ context.Context, input healthDTO.UpsertHealthInput) (db.HealthState, error) {
	f.upsertHealthCalls++
	f.upsertHealthInput = input
	return db.HealthState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteHealth(_ context.Context, input healthDTO.DeleteHealthInput) error {
	f.deleteHealthCalls++
	f.deleteHealthInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetCharacteristics(_ context.Context, input characteristicsDTO.GetCharacteristicsInput) (db.Characteristic, error) {
	f.getCharacteristicsCalls++
	f.getCharacteristicsInput = input
	return db.Characteristic{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertCharacteristics(_ context.Context, input characteristicsDTO.UpsertCharacteristicsInput) (db.Characteristic, error) {
	f.upsertCharacteristicsCalls++
	f.upsertCharacteristicsInput = input
	return db.Characteristic{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteCharacteristics(_ context.Context, input characteristicsDTO.DeleteCharacteristicsInput) error {
	f.deleteCharacteristicsCalls++
	f.deleteCharacteristicsInput = input
	return f.err
}

func (f *fakeCharacterHandlerService) GetSanity(_ context.Context, _ sanityDTO.GetSanityInput) (db.SanityState, error) {
	return db.SanityState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertSanity(_ context.Context, _ sanityDTO.UpsertSanityInput) (db.SanityState, error) {
	return db.SanityState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteSanity(_ context.Context, _ sanityDTO.DeleteSanityInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetMagic(_ context.Context, _ magicDTO.GetMagicInput) (db.MagicState, error) {
	return db.MagicState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertMagic(_ context.Context, _ magicDTO.UpsertMagicInput) (db.MagicState, error) {
	return db.MagicState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteMagic(_ context.Context, _ magicDTO.DeleteMagicInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetLuck(_ context.Context, _ luckDTO.GetLuckInput) (db.LuckState, error) {
	return db.LuckState{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertLuck(_ context.Context, _ luckDTO.UpsertLuckInput) (db.LuckState, error) {
	return db.LuckState{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteLuck(_ context.Context, _ luckDTO.DeleteLuckInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetFinances(_ context.Context, _ financesDTO.GetFinancesInput) (db.Finance, error) {
	return db.Finance{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertFinances(_ context.Context, _ financesDTO.UpsertFinancesInput) (db.Finance, error) {
	return db.Finance{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteFinances(_ context.Context, _ financesDTO.DeleteFinancesInput) error {
	return f.err
}

func (f *fakeCharacterHandlerService) GetDerivedStats(_ context.Context, _ derivedStatsDTO.GetDerivedStatsInput) (db.DerivedStat, error) {
	return db.DerivedStat{}, f.err
}

func (f *fakeCharacterHandlerService) UpsertDerivedStats(_ context.Context, _ derivedStatsDTO.UpsertDerivedStatsInput) (db.DerivedStat, error) {
	return db.DerivedStat{}, f.err
}

func (f *fakeCharacterHandlerService) DeleteDerivedStats(_ context.Context, _ derivedStatsDTO.DeleteDerivedStatsInput) error {
	return f.err
}
