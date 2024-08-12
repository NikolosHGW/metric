package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikolosHGW/metric/internal/server/config"
	"github.com/NikolosHGW/metric/internal/server/db"
	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/server/middlewares"
	"github.com/NikolosHGW/metric/internal/server/routes"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

const defaultTagValue = "N/A"

var (
	buildVersion = defaultTagValue
	buildDate    = defaultTagValue
	buildCommit  = defaultTagValue
)

func main() {
	if err := run(); err != nil {
		log.Fatal(fmt.Errorf("server/main run_ListenAndServe: %w", err))
	}
}

func run() error {
	config := config.NewConfig()

	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}

	database, err := db.InitDB(config.GetDBConnection())
	if err != nil {
		logger.Log.Info("init db", zap.Error(err))
	}
	if database != nil {
		defer func() {
			err := database.Close()
			if err != nil {
				logger.Log.Info("err close database", zap.Error(err))
			}
		}()
	}

	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)
	if database != nil {
		databaseStrg := storage.NewDBStorage(database, logger.Log)
		metricService = services.NewMetricService(databaseStrg)
	}
	handler := handlers.NewHandler(metricService, logger.Log)
	diskStrg := storage.NewDiskStorage(strg, logger.Log, config.GetFileStoragePath())
	diskService := services.NewDiskService(diskStrg, config.GetStoreInterval(), config.GetRestore())
	diskService.FillMetricStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go diskService.CollectMetrics(ctx)

	hashMiddleware := middlewares.NewHashMiddleware(config.GetKey())
	decryptMiddleware := middlewares.NewDecryptMiddleware(config.GetCryptoKeyPath(), logger.Log)

	r := routes.InitRouter(handler, hashMiddleware, decryptMiddleware)

	fmt.Println(
		"Build version: ", buildVersion, "\n",
		"Build date: ", buildDate, "\n",
		"Build commit: ", buildCommit,
	)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	server := &http.Server{
		Addr:    config.Address,
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case sig := <-signalChan:
		logger.Log.Info("Received signal, shutting down", zap.String("signal", sig.String()))
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	}

	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()

	cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Log.Info("Server exited gracefully")

	return nil
}
