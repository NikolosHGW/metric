package db

import (
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var DB *sqlx.DB

func InitDB() {
	db, err := sqlx.Connect("postgres", "user=nikolos password=abc123 dbname=metric sslmode=disable")
	if err != nil {
		logger.Log.Info("connect to postgres", zap.Error(err))

		return
	}

	DB = db
}
