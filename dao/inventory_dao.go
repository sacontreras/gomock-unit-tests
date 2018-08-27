package dao

import "coding-exercise/models"

type InventoryDao interface {
	GetInventoryItem(itemId string) (*models.InventoryItem, error)
}

//intentionally not implemented as unit test interactions with this interface should be mocked
