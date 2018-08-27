package services

import (
	"coding-exercise/dao"
	"coding-exercise/models"
	"errors"
)

const defaultMaxSuggestions = 10

type InventoryService interface {
	GetItemOrSuggestions(itemId string, maxNumberOfSuggestions int) (*models.InventoryItem, []models.InventoryItem, error)
}

type InventoryServiceImpl struct {
	dao               dao.InventoryDao
	suggestionService SuggestionService
}

//helper method to return your Inventory Service implementation
func NewInventoryService(dao dao.InventoryDao, suggestionService SuggestionService) *InventoryServiceImpl {
	return &InventoryServiceImpl{
		dao:               dao,
		suggestionService: suggestionService,
	}
}

// SuggestionService intentionally not implemented as unit test interactions with this interface should be mocked
type SuggestionService interface {
	GetSuggestions(itemId string, maxNumberOfSuggestions int) ([]models.InventoryItem, error)
}

func (service *InventoryServiceImpl) GetItemOrSuggestions(itemId string, maxNumberOfSuggestions int) (*models.InventoryItem, []models.InventoryItem, error) {

	item, error := service.dao.GetInventoryItem(itemId)

	if item == nil {
		return nil, nil, errors.New("Item does not exist")
	} else if error != nil {
		return nil, nil, error
	}

	if item.NumberInStock > 0 {
		return item, nil, nil
	}

	var suggestionCount = defaultMaxSuggestions
	if maxNumberOfSuggestions <= suggestionCount && maxNumberOfSuggestions > 0 {
		suggestionCount = maxNumberOfSuggestions
	}

	suggestedAlternatives, error := service.suggestionService.GetSuggestions(item.ItemId, suggestionCount)
	if error != nil {
		return nil, nil, error
	}

	return item, suggestedAlternatives, nil
}
