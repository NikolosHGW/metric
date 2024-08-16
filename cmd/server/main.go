package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikolosHGW/metric/internal/proto"
	"github.com/NikolosHGW/metric/internal/server/config"
	"github.com/NikolosHGW/metric/internal/server/db"
	"github.com/NikolosHGW/metric/internal/server/grpcserver"
	"github.com/NikolosHGW/metric/internal/server/interceptor"
	"github.com/NikolosHGW/metric/internal/server/logger"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"

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
	// handler := handlers.NewHandler(metricService, logger.Log)
	diskStrg := storage.NewDiskStorage(strg, logger.Log, config.GetFileStoragePath())
	diskService := services.NewDiskService(diskStrg, config.GetStoreInterval(), config.GetRestore())
	diskService.FillMetricStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go diskService.CollectMetrics(ctx)

	// hashMiddleware := middlewares.NewHashMiddleware(config.GetKey())
	// decryptMiddleware := middlewares.NewDecryptMiddleware(config.GetCryptoKeyPath(), logger.Log)
	// checkIP := middlewares.NewCheckIP(config.GetTrustedSubnet(), logger.Log)

	fmt.Println(
		"Build version: ", buildVersion, "\n",
		"Build date: ", buildDate, "\n",
		"Build commit: ", buildCommit,
	)

	logger.Log.Info("Running server", zap.String("address", config.Address))

	// server := &http.Server{
	// 	Addr:    config.Address,
	// 	Handler: r,
	// }

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	errChan := make(chan error, 1)
	grpcServerChan := make(chan *grpc.Server)

	go func() {
		grpcServer, err := startGRPCServer(config, *metricService, logger.Log)
		if err != nil {
			errChan <- err
		}
		grpcServerChan <- grpcServer
	}()

	select {
	case sig := <-signalChan:
		logger.Log.Info("Received signal, shutting down", zap.String("signal", sig.String()))
	case err := <-errChan:
		if err != nil && err != grpc.ErrServerStopped {
			return fmt.Errorf("server error: %w", err)
		}
	}

	cancel()

	grpcServer := <-grpcServerChan
	grpcServer.GracefulStop()

	logger.Log.Info("Server exited gracefully")

	return nil
}

type configer interface {
	GetAddress() string
	GetKey() string
	GetCryptoKeyPath() string
	GetTrustedSubnet() string
}

type customLogger interface {
	Info(string, ...zap.Field)
}

func startGRPCServer(config configer, metricService services.MetricService, log customLogger) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryLoggingInterceptor),
		grpc.UnaryInterceptor(interceptor.UnaryGzipInterceptor),
		grpc.UnaryInterceptor(interceptor.NewHashMiddleware(config.GetKey()).UnaryHashInterceptor),
		grpc.UnaryInterceptor(interceptor.NewDecryptMiddleware(config.GetCryptoKeyPath(), logger.Log).UnaryDecryptInterceptor),
		grpc.UnaryInterceptor(interceptor.NewCheckIP(config.GetTrustedSubnet(), logger.Log).UnaryCheckIPInterceptor),
	)
	proto.RegisterMetricServiceServer(grpcServer, grpcserver.NewMetricServiceServer(metricService, logger.Log))

	log.Info("Starting gRPC server at", zap.String("address", config.GetAddress()))

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Info("gRPC server stopped with error", zap.Error(err))
		}
	}()

	return grpcServer, nil
}
