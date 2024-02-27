package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/NikolosHGW/metric/internal/models"
)

// type repository interface {
// 	SetMetric(models.Metrics, context.Context)
// 	GetMetric(string, context.Context) (models.Metrics, error)
// 	SetGaugeMetric(string, models.Gauge, context.Context)
// 	SetCounterMetric(string, models.Counter, context.Context)
// 	GetGaugeMetric(string, context.Context) (models.Gauge, error)
// 	GetCounterMetric(string, context.Context) (models.Counter, error)
// 	GetAllMetrics(context.Context) []string
// }

func NewDBStorage(sql *sqlx.DB) *DBStorage {
	return &DBStorage{
		sql: sql,
	}
}

type DBStorage struct {
	sql *sqlx.DB
}

func (ds DBStorage) SetMetric(m models.Metrics, ctx context.Context) {
	ds.sql.ExecContext(ctx, "INSERT INTO metrics (country, telcode) VALUES (?, ?)")
}
