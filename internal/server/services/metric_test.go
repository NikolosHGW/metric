package services

import (
	"context"
	"errors"
	"testing"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/stretchr/testify/assert"
)

func f(v float64) *float64 {
	return &v
}

func i(v int64) *int64 {
	return &v
}

type mockRepo struct{}

func (m *mockRepo) SetMetric(ctx context.Context, metric models.Metrics) error {
	if metric.ID == "error" {
		return errors.New("test error")
	}
	return nil
}

func (m *mockRepo) GetMetric(ctx context.Context, name string) (models.Metrics, error) {
	if name == "error" {
		return models.Metrics{}, errors.New("test error")
	}
	return models.Metrics{ID: name, MType: models.GaugeType, Value: f(42.0)}, nil
}

func (m *mockRepo) SetGaugeMetric(ctx context.Context, name string, value models.Gauge) error {
	return nil
}

func (m *mockRepo) SetCounterMetric(ctx context.Context, name string, value models.Counter) error {
	return nil
}

func (m *mockRepo) GetGaugeMetric(ctx context.Context, name string) (models.Gauge, error) {
	return 42.0, nil
}

func (m *mockRepo) GetCounterMetric(ctx context.Context, name string) (models.Counter, error) {
	return 42, nil
}

func (m *mockRepo) GetAllMetrics(ctx context.Context) []string {
	return []string{"testGauge", "testCounter"}
}

func (m *mockRepo) GetIsDBConnected() bool {
	return true
}

func (m *mockRepo) UpsertMetrics(ctx context.Context, mc models.MetricCollection) (models.MetricCollection, error) {
	return mc, nil
}

func TestSetMetric(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	err := service.SetMetric(context.Background(), models.GaugeType, "testGauge", "42.0")
	assert.NoError(t, err)

	err = service.SetMetric(context.Background(), models.CounterType, "testCounter", "42")
	assert.NoError(t, err)
}

func TestGetMetricValue(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	value, err := service.GetMetricValue(context.Background(), models.GaugeType, "testGauge")
	assert.NoError(t, err)
	assert.Equal(t, "42", value)

	value, err = service.GetMetricValue(context.Background(), models.CounterType, "testCounter")
	assert.NoError(t, err)
	assert.Equal(t, "42", value)
}

func TestSetJSONMetric(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	metric := models.Metrics{ID: "testGauge", MType: models.GaugeType, Value: f(42.0)}
	err := service.SetJSONMetric(context.Background(), metric)
	assert.NoError(t, err)

	metric = models.Metrics{ID: "error", MType: models.GaugeType, Value: f(42.0)}
	err = service.SetJSONMetric(context.Background(), metric)
	assert.Error(t, err)
}

func TestGetMetricByName(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	metric, err := service.GetMetricByName(context.Background(), "testGauge")
	assert.NoError(t, err)
	assert.Equal(t, models.Metrics{ID: "testGauge", MType: models.GaugeType, Value: f(42.0)}, metric)

	_, err = service.GetMetricByName(context.Background(), "error")
	assert.Error(t, err)
}

func TestGetAllMetrics(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	metrics := service.GetAllMetrics(context.Background())
	assert.Equal(t, []string{"testGauge", "testCounter"}, metrics)
}

func TestGetIsDBConnected(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	dbConnected := service.GetIsDBConnected()
	assert.True(t, dbConnected)
}

func TestUpsertMetrics(t *testing.T) {
	repo := &mockRepo{}
	service := NewMetricService(repo)

	mc := models.MetricCollection{
		Metrics: []models.Metrics{
			{ID: "testGauge", MType: models.GaugeType, Value: f(42.0)},
			{ID: "testCounter", MType: models.CounterType, Delta: i(42)},
		},
	}
	upsertedMc, err := service.UpsertMetrics(context.Background(), mc)
	assert.NoError(t, err)
	assert.Equal(t, mc, upsertedMc)
}
