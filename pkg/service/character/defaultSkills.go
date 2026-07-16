package character

import (
	"context"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	dodgeBaseRule          = "dodge"
	nativeLanguageBaseRule = "native_language"
)

var defaultSkillsCategoryID = pgtype.UUID{
	Bytes: uuid.MustParse("1f81c838-4c15-4bdc-aabf-fc9699595cc8"),
	Valid: true,
}

type defaultCharacterSkill struct {
	name        string
	baseValue   int16
	isProtected bool
	baseRule    *string
}

var defaultCharacterSkills = []defaultCharacterSkill{
	{name: "Антропология", baseValue: 1, isProtected: true},
	{name: "Археология", baseValue: 1, isProtected: true},
	{name: "Ближний бой (драка)", baseValue: 25, isProtected: true},
	{name: "Бухгалтерское дело", baseValue: 5, isProtected: true},
	{name: "Верховая езда", baseValue: 5, isProtected: true},
	{name: "Взлом", baseValue: 1, isProtected: true},
	{name: "Внимание", baseValue: 25, isProtected: true},
	{name: "Вождение", baseValue: 20, isProtected: true},
	{name: "Выживание", baseValue: 10},
	{name: "Естествознание", baseValue: 10, isProtected: true},
	{name: "Запугивание", baseValue: 15, isProtected: true},
	{name: "Искусство/ремесло", baseValue: 5},
	{name: "История", baseValue: 5, isProtected: true},
	{name: "Красноречие", baseValue: 5, isProtected: true},
	{name: "Лазание", baseValue: 20, isProtected: true},
	{name: "Ловкость рук", baseValue: 10, isProtected: true},
	{name: "Маскировка", baseValue: 5, isProtected: true},
	{name: "Медицина", baseValue: 1, isProtected: true},
	{name: "Метание", baseValue: 20, isProtected: true},
	{name: "Механика", baseValue: 10, isProtected: true},
	{name: "Мифы Ктулху", baseValue: 0, isProtected: true},
	{name: "Наука", baseValue: 1},
	{name: "Обаяние", baseValue: 15, isProtected: true},
	{name: "Обоняние", baseValue: 15, isProtected: true},
	{name: "Оккультизм", baseValue: 5, isProtected: true},
	{name: "Ориентирование", baseValue: 10, isProtected: true},
	{name: "Оценка", baseValue: 5, isProtected: true},
	{name: "Первая помощь", baseValue: 30, isProtected: true},
	{name: "Пилотирование", baseValue: 1, isProtected: true},
	{name: "Плавание", baseValue: 20, isProtected: true},
	{name: "Прыжки", baseValue: 20, isProtected: true},
	{name: "Психоанализ", baseValue: 1, isProtected: true},
	{name: "Психология", baseValue: 10, isProtected: true},
	{name: "Работа в библиотеке", baseValue: 20, isProtected: true},
	{name: "Скрытность", baseValue: 20, isProtected: true},
	{name: "Слух", baseValue: 20, isProtected: true},
	{name: "Стрельба (винт./дроб.)", baseValue: 25, isProtected: true},
	{name: "Стрельба (пистолет)", baseValue: 20, isProtected: true},
	{name: "Убеждение", baseValue: 10, isProtected: true},
	{name: "Уклонение", isProtected: true, baseRule: stringPointer(dodgeBaseRule)},
	{name: "Управление тяжелыми машинами", baseValue: 1, isProtected: true},
	{name: "Чтение следов", baseValue: 10, isProtected: true},
	{name: "Электрика", baseValue: 10, isProtected: true},
	{name: "Юриспруденция", baseValue: 5, isProtected: true},
	{name: "Язык, иностранный", baseValue: 1},
	{name: "Язык, родной", isProtected: true, baseRule: stringPointer(nativeLanguageBaseRule)},
}

func createDefaultCharacterSkills(ctx context.Context, queries *db.Queries, userID string, characterID pgtype.UUID) error {
	for _, skill := range defaultCharacterSkills {
		_, err := queries.CreateCharacterSkill(ctx, db.CreateCharacterSkillParams{
			UserID:      userID,
			CharacterID: characterID,
			Name:        skill.name,
			CategoryID:  defaultSkillsCategoryID,
			BaseValue:   skill.baseValue,
			Value:       0,
			Checked:     false,
			IsProtected: skill.isProtected,
			BaseRule:    skill.baseRule,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func stringPointer(value string) *string {
	return &value
}
