package db

import (
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var DB *sqlx.DB

func InitDB(dataSourceName string) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		logger.Log.Info("connect to postgres", zap.Error(err))

		return
	}

	DB = db
}
