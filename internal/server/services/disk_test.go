package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDiskStorage struct {
	mock.Mock
}

func (m *MockDiskStorage) WriteToDisk() {
	m.Called()
}

func (m *MockDiskStorage) WriteToStorage() {
	m.Called()
}

func (m *MockDiskStorage) CanWriteToDisk() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestDiskService_CollectMetrics(t *testing.T) {
	mockDiskStorage := new(MockDiskStorage)
	mockDiskStorage.On("CanWriteToDisk").Return(true)
	mockDiskStorage.On("WriteToDisk").Return()

	diskService := NewDiskService(mockDiskStorage, 1, false)

	go diskService.CollectMetrics(context.Background())
	time.Sleep(3 * time.Second)

	mockDiskStorage.AssertCalled(t, "WriteToDisk")
}

func TestDiskService_FillMetricStorage(t *testing.T) {
	mockDiskStorage := new(MockDiskStorage)
	mockDiskStorage.On("WriteToStorage").Return()

	diskService := NewDiskService(mockDiskStorage, 1, true)

	diskService.FillMetricStorage()

	mockDiskStorage.AssertCalled(t, "WriteToStorage")
}
