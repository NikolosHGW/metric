package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

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

func (ds DBStorage) SetMetric(m models.Metrics, ctx context.Context) {
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
	if err != nil {
		ds.log.Info("cannot exec set metric", zap.Error(err))
	}
}

func (ds DBStorage) GetMetric(name string, ctx context.Context) (models.Metrics, error) {
	row := ds.sql.QueryRowxContext(ctx, "SELECT id, type, delta, value FROM metrics WHERE id = &1", name)

	model := models.Metrics{}

	err := row.Scan(&model.ID, &model.MType, &model.Delta, &model.Value)

	if err != nil {
		ds.log.Info("cannot scan row when getting metric", zap.Error(err))
	}

	return model, err
}

func (ds DBStorage) SetGaugeMetric(name string, value models.Gauge, ctx context.Context) {

}
