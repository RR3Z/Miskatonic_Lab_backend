package tests

import (
	"context"
	"errors"
	"testing"

	characterEvents "github.com/RR3Z/Miskatonic_Lab_backend/pkg/events/character"
	inventoryDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character/inventory"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/stretchr/testify/require"
)

func TestEventPublishingCharacterServicePublishesInventorySuccessEvents(t *testing.T) {
	characterID := testUUID(testCharacterID)
	itemID := testUUID(testItemID)

	cases := []eventPublishingCase{
		{
			name: "list inventory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetInventoryItems(ctx, inventoryDTO.GetInventoryItemsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemsListSucceeded{UserID: testUserID, CharacterID: testCharacterID, Count: 1},
		},
		{
			name: "get inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetInventoryItem(ctx, inventoryDTO.GetInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemGetSucceeded{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Name: "Pocket Flashlight"},
		},
		{
			name: "create inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateInventoryItem(ctx, inventoryDTO.CreateInventoryItemInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemCreateSucceeded{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Name: "Pocket Flashlight"},
		},
		{
			name: "update inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateInventoryItem(ctx, inventoryDTO.UpdateInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemUpdateSucceeded{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Name: "Pocket Flashlight"},
		},
		{
			name: "delete inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteInventoryItem(ctx, inventoryDTO.DeleteInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterInventoryItemDeleteSucceeded{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, publisher, service := newEventPublishingTestSubject()
			require.NoError(t, tc.call(context.Background(), service))
			requirePublishedEvent(t, publisher, tc.expectedEvent)
		})
	}
}

func TestEventPublishingCharacterServicePublishesInventoryFailureEvents(t *testing.T) {
	characterID := testUUID(testCharacterID)
	itemID := testUUID(testItemID)
	expectedErr := errors.New("base service failed")

	cases := []eventPublishingCase{
		{
			name: "list inventory items",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetInventoryItems(ctx, inventoryDTO.GetInventoryItemsInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemsListFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "get inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.GetInventoryItem(ctx, inventoryDTO.GetInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemGetFailed{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Err: expectedErr},
		},
		{
			name: "create inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.CreateInventoryItem(ctx, inventoryDTO.CreateInventoryItemInput{UserID: testUserID, CharacterID: characterID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemCreateFailed{UserID: testUserID, CharacterID: testCharacterID, Err: expectedErr},
		},
		{
			name: "update inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				_, err := service.UpdateInventoryItem(ctx, inventoryDTO.UpdateInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
				return err
			},
			expectedEvent: characterEvents.CharacterInventoryItemUpdateFailed{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Err: expectedErr},
		},
		{
			name: "delete inventory item",
			call: func(ctx context.Context, service *characterServices.EventPublishingCharacterService) error {
				return service.DeleteInventoryItem(ctx, inventoryDTO.DeleteInventoryItemInput{UserID: testUserID, CharacterID: characterID, ItemID: itemID})
			},
			expectedEvent: characterEvents.CharacterInventoryItemDeleteFailed{UserID: testUserID, CharacterID: testCharacterID, InventoryID: testItemID, Err: expectedErr},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next, publisher, service := newEventPublishingTestSubject()
			next.Err = expectedErr
			require.ErrorIs(t, tc.call(context.Background(), service), expectedErr)
			requirePublishedEvent(t, publisher, tc.expectedEvent)
		})
	}
}
