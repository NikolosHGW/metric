package grpcserver

import (
	"context"

	"github.com/NikolosHGW/metric/internal/models"
	"github.com/NikolosHGW/metric/internal/proto"
	"go.uber.org/zap"
)

type metricService interface {
	GetMetricByName(context.Context, string) (models.Metrics, error)
	UpsertMetrics(context.Context, models.MetricCollection) (models.MetricCollection, error)
}

type customLogger interface {
	Info(string, ...zap.Field)
}

type MetricServiceServer struct {
	proto.UnimplementedMetricServiceServer
	metricService metricService
	logger        customLogger
}

func NewMetricServiceServer(ms metricService, logger customLogger) *MetricServiceServer {
	return &MetricServiceServer{
		metricService: ms,
		logger:        logger,
	}
}

func (s *MetricServiceServer) GetMetric(ctx context.Context, req *proto.MetricRequest) (*proto.MetricResponse, error) {
	metric, err := s.metricService.GetMetricByName(ctx, req.Id)
	if err != nil {
		s.logger.Info("metric not found", zap.Error(err))
		return nil, err
	}

	return &proto.MetricResponse{
		Id:    metric.ID,
		Type:  metric.MType,
		Value: *metric.Value,
		Delta: *metric.Delta,
	}, nil
}

func (s *MetricServiceServer) UpsertMetrics(ctx context.Context, req *proto.UpsertMetricRequest) (*proto.UpsertMetricResponse, error) {
	metricCollection := models.MetricCollection{}
	for _, m := range req.Metrics {
		metric := models.Metrics{
			ID:    m.Id,
			MType: m.Type,
			Value: &m.Value,
			Delta: &m.Delta,
		}
		metricCollection.Metrics = append(metricCollection.Metrics, metric)
	}

	metrics, err := s.metricService.UpsertMetrics(ctx, metricCollection)
	if err != nil {
		s.logger.Info("cannot upsert metrics", zap.Error(err))
		return nil, err
	}

	var responseMetrics []*proto.Metric
	for _, m := range metrics.Metrics {
		responseMetrics = append(responseMetrics, &proto.Metric{
			Id:    m.ID,
			Type:  m.MType,
			Value: *m.Value,
			Delta: *m.Delta,
		})
	}

	return &proto.UpsertMetricResponse{
		Metrics: responseMetrics,
	}, nil
}
