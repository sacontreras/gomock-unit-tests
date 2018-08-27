package services

import (
	"coding-exercise/mocks"
	"coding-exercise/models"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

//Unit tests go here
//You may choose any go mocking framework you think works best for you
//Please aim to cover all the scenarios/paths in the GetItemOrSuggestions method

type TestCaseId uint8
const aplhabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_=+;,./"
type Util struct {}
func (util* Util) InitRand() {
	rand.Seed(time.Now().UnixNano())
}
func (util* Util) RandStringBytes(n int) string {	//borrowed from StackOverflow
	util.InitRand()
	b := make([]byte, n)
	for i := range b {
		b[i] = aplhabet[rand.Intn(len(aplhabet))]
	}
	return string(b)
}
func (util* Util) RandPositiveInt(n int) int {
	util.InitRand()
	return int(math.Max(1, float64(rand.Intn(n))))
}

func TestInventoryServiceImplGetItemOrSuggestions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//construct mock objects
	mockInventoryDao := mocks.NewMockInventoryDao(mockCtrl)
	mockSuggestionService := mocks.NewMockSuggestionService(mockCtrl)

	//construct InventoryServiceImpl with mock objects
	inventoryServiceImpl := NewInventoryService(mockInventoryDao, mockSuggestionService)

	//****************************************** ARRANGE - begin ******************************************/
	const (
		NonexistentInventoryItem 						TestCaseId = iota
		DevilInTheDAO
		ItemExistsAndInStock
		ItemExistsButOutOfStock
		ItemExistsButOutOfStock_InvalidMaxSuggestions
		ItemExistsButOutOfStock_SuggestionServiceOffline
	)
	var (
		testCaseResultAsserters = map[TestCaseId]func(*testing.T, string, int, *models.InventoryItem, []models.InventoryItem, error) {
			NonexistentInventoryItem:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					//requestedItemId and requestedMaxSuggestions irrelevant here
					assert.Nil(t, inventoryItemResult)
					assert.Nil(t, suggestAlternativesResult)
					assert.NotNil(t, errorResult)
					assert.Equal(t, "Item does not exist", errorResult.Error())
				},

			ItemExistsAndInStock:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					//requestedMaxSuggestions irrelevant here
					assert.NotNil(t, inventoryItemResult)
					assert.Equal(t, requestedItemId, inventoryItemResult.ItemId)
					assert.True(t, inventoryItemResult.NumberInStock > 0)
					assert.Nil(t, suggestAlternativesResult)
					assert.Nil(t, errorResult)
				},

			DevilInTheDAO:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					//requestedItemId and requestedMaxSuggestions irrelevant here
					assert.Nil(t, inventoryItemResult)
					assert.Nil(t, suggestAlternativesResult)
					assert.NotNil(t, errorResult)
					assert.Equal(t, "the DEVIL is the only explanation for this! mauahahahah...", errorResult.Error())
				},

			ItemExistsButOutOfStock_InvalidMaxSuggestions:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					assert.NotNil(t, inventoryItemResult)
					assert.Equal(t, requestedItemId, inventoryItemResult.ItemId)
					assert.True(t, inventoryItemResult.NumberInStock == 0)
					assert.NotNil(t, suggestAlternativesResult)
					assert.NotEqual(t, requestedMaxSuggestions, len(suggestAlternativesResult))
					assert.Equal(t, defaultMaxSuggestions, len(suggestAlternativesResult))
					assert.Nil(t, errorResult)
				},

			ItemExistsButOutOfStock:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					assert.NotNil(t, inventoryItemResult)
					assert.Equal(t, requestedItemId, inventoryItemResult.ItemId)
					assert.True(t, inventoryItemResult.NumberInStock == 0)
					assert.NotNil(t, suggestAlternativesResult)
					assert.Equal(t, requestedMaxSuggestions, len(suggestAlternativesResult))
					assert.Nil(t, errorResult)
				},

			ItemExistsButOutOfStock_SuggestionServiceOffline:
				func(t *testing.T, requestedItemId string, requestedMaxSuggestions int, inventoryItemResult *models.InventoryItem, suggestAlternativesResult []models.InventoryItem, errorResult error) {
					assert.Nil(t, inventoryItemResult)
					assert.Nil(t, suggestAlternativesResult)
					assert.NotNil(t, errorResult)
					assert.Equal(t, "Suggestion service offline", errorResult.Error())
				},
		}
	)
	const (
		BDDTrigger__WellformedId__doesnotexist           						= "doesnotexist"
		BDDTrigger__WellformedId__finditembutreturnerror 						= "finditembutreturnerror"
		BDDTrigger__WellformedId__existsandinstock       						= "existsandinstock"
		BDDTrigger__WellformedId__existsbutoutofstocksuggestionserviceoffline	= "existsbutoutofstocksuggestionserviceoffline"
	)

	mockInventoryDao.EXPECT().GetInventoryItem(BDDTrigger__WellformedId__doesnotexist).
		DoAndReturn(
			func(itemId string) (*models.InventoryItem, error) {
				fmt.Printf("mockInventoryDao.GetInventoryItem: \tnon-existent-item handler; return nil, nil\n")
				return nil, nil
			},
		).
		AnyTimes()
	mockInventoryDao.EXPECT().GetInventoryItem(BDDTrigger__WellformedId__finditembutreturnerror).
		DoAndReturn(
			func(itemId string) (*models.InventoryItem, error) {
				item := models.InventoryItem{ItemId: "NOTtheEntityIdYouExpected?", NumberInStock: -666}
				err := errors.New("the DEVIL is the only explanation for this! mauahahahah...")
				fmt.Printf("mockInventoryDao.GetInventoryItem: \tdevil-in-the-dao; return item{\"%s\", %d}, error(\"%s\")\n", item.ItemId, item.NumberInStock, err.Error())
				return &item, err
			},
		).
		AnyTimes()
	mockInventoryDao.EXPECT().GetInventoryItem(BDDTrigger__WellformedId__existsandinstock).
		DoAndReturn(
			func(itemId string) (*models.InventoryItem, error) {
				item := models.InventoryItem{ItemId: BDDTrigger__WellformedId__existsandinstock, NumberInStock: 1}
				fmt.Printf("mockInventoryDao.GetInventoryItem: \texists-and-in-stock-item handler; return item{\"%s\", %d}, nil\n", item.ItemId, item.NumberInStock)
				return &item, nil
			},
		).
		AnyTimes()
	mockInventoryDao.EXPECT().GetInventoryItem(gomock.Any()).
		DoAndReturn(
			func(itemId string) (*models.InventoryItem, error) {
				if len(strings.TrimSpace(itemId)) == 0 {	//this covers the malformed id cases - falls into non-existent item behavior since nil inventory item is returned
					errMsg := "itemId string cannot be empty"
					fmt.Printf("mockInventoryDao.GetInventoryItem: \tempty-id handler; return nil, error(\"%s\")\n", errMsg)
					return nil, errors.New(errMsg)		//note that this has NO effect though since maps every case for nil inv item to the same error "Item does not exist" (which is okay)
				} else {
					if strings.Contains(strings.TrimSpace(itemId), " ") {
						errMsg := fmt.Sprintf("itemId '%s' is malformed!", itemId)
						fmt.Printf("mockInventoryDao.GetInventoryItem: \tmalformed-id handler; return nil, error(\"%s\")\n", errMsg)
						return nil, errors.New(errMsg)		//note that this has NO effect though since maps every case for nil inv item to the same error "Item does not exist" (which is okay)
					} else { //every other case should emulate existing item that is out of stock
						item := models.InventoryItem{ItemId: itemId, NumberInStock: 0}
						fmt.Printf("mockInventoryDao.GetInventoryItem: \texists-and-out-of-stock-item handler; return item{\"%s\", %d}, nil\n", item.ItemId, item.NumberInStock)
						return &item, nil
					}
				}
			},
		).
		AnyTimes()

	mockSuggestionService.EXPECT().GetSuggestions(BDDTrigger__WellformedId__existsbutoutofstocksuggestionserviceoffline, gomock.Any()).
		DoAndReturn(
			func(itemId string, maxNumberOfSuggestions int) ([]models.InventoryItem, error) {
				fmt.Printf("mockSuggestionService.GetSuggestions: \tid: \"%s\", maxNumSuggestions: %d\n", itemId, maxNumberOfSuggestions)
				err := errors.New("Suggestion service offline")
				fmt.Printf("mockSuggestionService.GetSuggestions: \t\tid: error(\"%s\")\n", err.Error())
				return nil, err
			},
		).
		AnyTimes()
	mockSuggestionService.EXPECT().GetSuggestions(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(itemId string, maxNumberOfSuggestions int) ([]models.InventoryItem, error) {
				fmt.Printf("mockSuggestionService.GetSuggestions: \tid: \"%s\", maxNumSuggestions: %d\n", itemId, maxNumberOfSuggestions)
				alternatives := make([]models.InventoryItem, maxNumberOfSuggestions)
				util := Util{}
				for i := 0; i < maxNumberOfSuggestions; i++ {
					alternatives[i] = models.InventoryItem{util.RandStringBytes(8), util.RandPositiveInt(100)}
					fmt.Printf("mockSuggestionService.GetSuggestions: \t\tsuggestedAlternative[%d] - id: \"%s\", stockCount: %d\n", i, alternatives[i].ItemId, alternatives[i].NumberInStock)
				}
				return alternatives, nil
			},
		).
		AnyTimes()
	//****************************************** ARRANGE - end ******************************************/

	t.Run(
		"MalformedId",
		func(t *testing.T) {
			//each of this subcases should result in the same outcome - non-existent inventory item
			var malformedId string
			t.Run(
				"IdIsEmptyString",
				func(t *testing.T) {
					//Act
					malformedId = ""
					maxSuggestions := 1
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(malformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[NonexistentInventoryItem](
						t,
						malformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)

			t.Run(
				"IdIsOnlyWhitespace",
				func(t *testing.T) {
					//Act
					malformedId = "    "
					maxSuggestions := 1
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(malformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[NonexistentInventoryItem](
						t,
						malformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)

			t.Run(
				"IdHasWhiteSpaceInterspersed",
				func(t *testing.T) {
					//Act
					malformedId = "  s lkkd  $*( llk"
					maxSuggestions := 1
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(malformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[NonexistentInventoryItem](
						t,
						malformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)
		},
	)

	t.Run(
		"WellformedId",
		func(t *testing.T) {
			t.Run(
				"NonexistentInventoryItem",
				func(t *testing.T) {
					//Act
					wellformedId := BDDTrigger__WellformedId__doesnotexist
					maxSuggestions := 130536	//why not since this value should be irrelevant
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[NonexistentInventoryItem](
						t,
						wellformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)

			t.Run(
				"DaoFindsItemForThisIdButSomehowReturnsError",
				func(t *testing.T) {
					//Act
					wellformedId := BDDTrigger__WellformedId__finditembutreturnerror
					maxSuggestions := -34878	//why not since this value should be irrelevant
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[DevilInTheDAO](
						t,
						wellformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)

			t.Run(
				"ItemForIdExistsAndIsInStock",
				func(t *testing.T) {
					//Act
					wellformedId := BDDTrigger__WellformedId__existsandinstock
					maxSuggestions := 0	//why not since this value should be irrelevant
					inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, maxSuggestions)

					//Assert
					testCaseResultAsserters[ItemExistsAndInStock](
						t,
						wellformedId,
						maxSuggestions,
						inventoryItemResult,
						suggestedAlternativesResult,
						errResult,
					)
				},
			)

			t.Run(
				"ItemForIdExistsButIsOutOfStock",
				func(t *testing.T) {
					//generate a random id in order to ensure we trigger the out-of-stock handler (see above)
					util := Util{}
					wellformedId := util.RandStringBytes(8)

					t.Run(
						"InvalidNumMaxSuggestions",
						func(t *testing.T) {
							var invalidMaxSuggestions int
							t.Run(
								"Negative",
								func(t *testing.T) {
									//Act
									invalidMaxSuggestions = -1	//this value is relevant!
									inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, invalidMaxSuggestions)

									//Assert
									testCaseResultAsserters[ItemExistsButOutOfStock_InvalidMaxSuggestions](
										t,
										wellformedId,
										invalidMaxSuggestions,
										inventoryItemResult,
										suggestedAlternativesResult,
										errResult,
									)
								},
							)

							t.Run(
								"Zero",
								func(t *testing.T) {
									//Act
									invalidMaxSuggestions = 0	//this value is relevant!
									inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, invalidMaxSuggestions)

									//Assert
									testCaseResultAsserters[ItemExistsButOutOfStock_InvalidMaxSuggestions](
										t,
										wellformedId,
										invalidMaxSuggestions,
										inventoryItemResult,
										suggestedAlternativesResult,
										errResult,
									)
								},
							)

							t.Run(
								"GreaterThanDefault",
								func(t *testing.T) {
									//Act
									invalidMaxSuggestions = defaultMaxSuggestions + 1	//this value is relevant!
									inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, invalidMaxSuggestions)

									//Assert
									testCaseResultAsserters[ItemExistsButOutOfStock_InvalidMaxSuggestions](
										t,
										wellformedId,
										invalidMaxSuggestions,
										inventoryItemResult,
										suggestedAlternativesResult,
										errResult,
									)
								},
							)
						},
					)
					t.Run(
						"ValidNumMaxSuggestions",
						func(t *testing.T) {
							var validMaxSuggestions int
							t.Run(
								"LessThanDefault",
								func(t *testing.T) {
									//Act
									validMaxSuggestions = util.RandPositiveInt(defaultMaxSuggestions) //[1, defaultMaxSuggestions)
									inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, validMaxSuggestions)

									//Assert
									testCaseResultAsserters[ItemExistsButOutOfStock](
										t,
										wellformedId,
										validMaxSuggestions,
										inventoryItemResult,
										suggestedAlternativesResult,
										errResult,
									)
								},
							)

							t.Run(
								"SuggestionServiceOffline",
								func(t *testing.T) {
									//Act
									wellformedId = BDDTrigger__WellformedId__existsbutoutofstocksuggestionserviceoffline
									validMaxSuggestions = util.RandPositiveInt(defaultMaxSuggestions) //[1, defaultMaxSuggestions)
									inventoryItemResult, suggestedAlternativesResult, errResult := inventoryServiceImpl.GetItemOrSuggestions(wellformedId, validMaxSuggestions)

									//Assert
									testCaseResultAsserters[ItemExistsButOutOfStock_SuggestionServiceOffline](
										t,
										wellformedId,
										validMaxSuggestions,
										inventoryItemResult,
										suggestedAlternativesResult,
										errResult,
									)
								},
							)
						},
					)
				},
			)
		},
	)
}
