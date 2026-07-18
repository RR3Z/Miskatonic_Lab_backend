package logging

import (
	"context"
	"log/slog"
	"reflect"
	"strings"
	"unicode"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/events"
	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	diceEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/dice"
	roomEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/room"
)

type EventLogger struct {
	logger   *slog.Logger
	registry events.DescriptorRegistry
}

func NewEventLogger(logger *slog.Logger, registry events.DescriptorRegistry) *EventLogger {
	if logger == nil {
		logger = slog.Default()
	}

	return &EventLogger{logger: logger, registry: registry}
}

func NewDefaultEventLogger(logger *slog.Logger) *EventLogger {
	return NewEventLogger(logger, events.NewDescriptorRegistry(
		characterEvents.Descriptors(),
		diceEvents.Descriptors(),
		roomEvents.Descriptors(),
	))
}

func (l *EventLogger) Handle(ctx context.Context, event events.Event) {
	descriptor, ok := l.registry.Describe(event)
	if !ok {
		return
	}

	attrs := []any{
		"event", event.EventName(),
		"domain", descriptor.Domain,
		"resource", descriptor.Resource,
		"action", descriptor.Action,
		"outcome", string(descriptor.Outcome),
	}
	attrs = append(attrs, eventAttrs(event)...)

	switch descriptor.Outcome {
	case events.OutcomeFailed:
		l.logger.ErrorContext(ctx, "domain event failed", attrs...)
	case events.OutcomeSkipped:
		l.logger.WarnContext(ctx, "domain event skipped", attrs...)
	default:
		l.logger.InfoContext(ctx, "domain event succeeded", attrs...)
	}
}

func eventAttrs(event events.Event) []any {
	value := reflect.ValueOf(event)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	attrs := make([]any, 0, value.NumField()*2)
	eventType := value.Type()
	for i := range value.NumField() {
		field := eventType.Field(i)
		if field.PkgPath != "" || field.Name == "Details" {
			continue
		}

		fieldValue := value.Field(i)
		if fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		attrs = append(attrs, eventFieldName(field.Name), fieldValue.Interface())
	}

	return attrs
}

func eventFieldName(name string) string {
	if fieldName, ok := eventFieldNames[name]; ok {
		return fieldName
	}

	var builder strings.Builder
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(r))
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

var eventFieldNames = map[string]string{
	"ActorID":         "actor_id",
	"ActorUserID":     "actor_user_id",
	"BackstoryItemID": "backstory_item_id",
	"CharacterID":     "character_id",
	"DeletedCount":    "deleted_count",
	"DeletedRoomID":   "deleted_room_id",
	"Err":             "error",
	"EventID":         "event_id",
	"InactiveDeleted": "inactive_deleted",
	"InventoryID":     "inventory_id",
	"InvalidDeleted":  "invalid_deleted",
	"MemberID":        "member_id",
	"NewOwnerID":      "new_owner_id",
	"NoteID":          "note_id",
	"OwnerID":         "owner_id",
	"RollID":          "roll_id",
	"RoomID":          "room_id",
	"SkillID":         "skill_id",
	"TargetUserID":    "target_user_id",
	"UserID":          "user_id",
}
