package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/NikolosHGW/metric/internal/models"
)

func NewDBStorage(sql *sqlx.DB, log customLogger) *DBStorage {
	return &DBStorage{
		sql: sql,
		log: log,
	}
}

type DBStorage struct {
	sql *sqlx.DB
	log customLogger
}

func (ds DBStorage) SetMetric(m models.Metrics, ctx context.Context) error {
	_, err := ds.sql.ExecContext(
		ctx,
		`INSERT INTO metrics (id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, 
			delta = EXCLUDED.delta, 
			value = EXCLUDED.value`,
		m.ID,
		m.MType,
		m.Delta,
		m.Value,
	)

	return err
}

func (ds DBStorage) GetMetric(name string, ctx context.Context) (models.Metrics, error) {
	row := ds.sql.QueryRowxContext(ctx, "SELECT id, type, delta, value FROM metrics WHERE id = $1", name)

	model := models.Metrics{}

	err := row.Scan(&model.ID, &model.MType, &model.Delta, &model.Value)

	if err != nil {
		ds.log.Info("cannot scan row when getting metric", zap.Error(err))
	}

	return model, err
}

func (ds DBStorage) SetGaugeMetric(name string, value models.Gauge, ctx context.Context) error {
	metric := models.Metrics{
		ID:    name,
		MType: models.GaugeType,
		Value: (*float64)(&value),
	}

	return ds.SetMetric(metric, ctx)
}

func (ds DBStorage) SetCounterMetric(name string, value models.Counter, ctx context.Context) error {
	metric := models.Metrics{
		ID:    name,
		MType: models.CounterType,
		Delta: (*int64)(&value),
	}

	return ds.SetMetric(metric, ctx)
}

func (ds DBStorage) GetGaugeMetric(name string, ctx context.Context) (models.Gauge, error) {
	metric, err := ds.GetMetric(name, ctx)

	if err != nil {
		var value models.Gauge
		return value, err
	}

	return models.Gauge(*metric.Value), err
}

func (ds DBStorage) GetCounterMetric(name string, ctx context.Context) (models.Counter, error) {
	metric, err := ds.GetMetric(name, ctx)

	if err != nil {
		var value models.Counter
		return value, err
	}

	return models.Counter(*metric.Delta), err
}

func (ds DBStorage) GetAllMetrics(ctx context.Context) []string {
	rows, err := ds.sql.QueryxContext(ctx, "SELECT id, type, delta, value FROM metrics ORDER BY id")

	var metricStrings []string
	if err != nil {
		ds.log.Info("cannot get all metric", zap.Error(err))
		return metricStrings
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var mType string
		var delta int64
		var value float64
		err = rows.Scan(&name, &mType, &delta, &value)

		if err != nil {
			ds.log.Info("cannot Scan", zap.Error(err))
			return metricStrings
		}

		result := fmt.Sprintf("%v", delta)
		if mType == models.GaugeType {
			result = fmt.Sprintf("%v", value)
		}

		metricStrings = append(metricStrings, fmt.Sprintf("%v: %v", name, result))
	}

	err = rows.Err()
	if err != nil {
		ds.log.Info("rows Err", zap.Error(err))
		return metricStrings
	}

	return metricStrings
}

func (ds DBStorage) GetIsDBConnected() bool {
	err := ds.sql.DB.Ping()

	return err == nil
}

func (ds *DBStorage) UpsertMetrics(metricCollection models.MetricCollection, ctx context.Context) (models.MetricCollection, error) {
	var upsertedMetrics []models.Metrics

	tx, err := ds.sql.BeginTxx(ctx, nil)
	if err != nil {
		return *models.NewMetricCollection(), err
	}

	for _, metric := range metricCollection.Metrics {
		var upsertedMetric models.Metrics
		err := tx.GetContext(ctx, &upsertedMetric,
			`INSERT INTO metrics (id, type, delta, value)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET
                type = EXCLUDED.type,
                delta = metrics.delta + EXCLUDED.delta,
                value = EXCLUDED.value
            RETURNING *`,
			metric.ID, metric.MType, metric.Delta, metric.Value,
		)
		if err != nil {
			tx.Rollback()
			return *models.NewMetricCollection(), err
		}
		upsertedMetrics = append(upsertedMetrics, upsertedMetric)
	}

	if err := tx.Commit(); err != nil {
		return *models.NewMetricCollection(), err
	}

	return models.MetricCollection{Metrics: upsertedMetrics}, nil
}
