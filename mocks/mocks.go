package mocks

import (
	models "coding-exercise/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockInventoryDao is a mock of InventoryDao interface
type MockInventoryDao struct {
	ctrl     *gomock.Controller
	recorder *MockInventoryDaoMockRecorder
}

// MockInventoryDaoMockRecorder is the mock recorder for MockInventoryDao
type MockInventoryDaoMockRecorder struct {
	mock *MockInventoryDao
}

// NewMockInventoryDao creates a new mock instance
func NewMockInventoryDao(ctrl *gomock.Controller) *MockInventoryDao {
	mock := &MockInventoryDao{ctrl: ctrl}
	mock.recorder = &MockInventoryDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInventoryDao) EXPECT() *MockInventoryDaoMockRecorder {
	return m.recorder
}

// GetInventoryItem mocks base method
func (m *MockInventoryDao) GetInventoryItem(arg0 string) (*models.InventoryItem, error) {
	ret := m.ctrl.Call(m, "GetInventoryItem", arg0)
	ret0, _ := ret[0].(*models.InventoryItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInventoryItem indicates an expected call of GetInventoryItem
func (mr *MockInventoryDaoMockRecorder) GetInventoryItem(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInventoryItem", reflect.TypeOf((*MockInventoryDao)(nil).GetInventoryItem), arg0)
}


// MockSuggestionService is a mock of SuggestionService interface
type MockSuggestionService struct {
	ctrl     *gomock.Controller
	recorder *MockSuggestionServiceMockRecorder
}

// MockSuggestionServiceMockRecorder is the mock recorder for MockSuggestionService
type MockSuggestionServiceMockRecorder struct {
	mock *MockSuggestionService
}

// NewMockSuggestionService creates a new mock instance
func NewMockSuggestionService(ctrl *gomock.Controller) *MockSuggestionService {
	mock := &MockSuggestionService{ctrl: ctrl}
	mock.recorder = &MockSuggestionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSuggestionService) EXPECT() *MockSuggestionServiceMockRecorder {
	return m.recorder
}

// GetSuggestions mocks base method
func (m *MockSuggestionService) GetSuggestions(arg0 string, arg1 int) ([]models.InventoryItem, error) {
	ret := m.ctrl.Call(m, "GetSuggestions", arg0, arg1)
	ret0, _ := ret[0].([]models.InventoryItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSuggestions indicates an expected call of GetSuggestions
func (mr *MockSuggestionServiceMockRecorder) GetSuggestions(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuggestions", reflect.TypeOf((*MockSuggestionService)(nil).GetSuggestions), arg0, arg1)
}