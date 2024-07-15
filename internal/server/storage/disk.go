package storage

import (
	"context"
	"encoding/json"
	"os"

	"github.com/NikolosHGW/metric/internal/models"
	"go.uber.org/zap"
)

type Storage interface {
	GetMetricsModels(context.Context) []models.Metrics
	SetMetric(context.Context, models.Metrics) error
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteMetric(metric *models.Metrics) error {
	return p.encoder.Encode(&metric)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetric() (*models.Metrics, error) {
	metric := &models.Metrics{}
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, err
	}

	return metric, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

type customLogger interface {
	Info(string, ...zap.Field)
}

type DiskStorage struct {
	strg     Storage
	log      customLogger
	fileName string
}

func NewDiskStorage(strg Storage, log customLogger, fileName string) *DiskStorage {
	return &DiskStorage{
		strg:     strg,
		log:      log,
		fileName: fileName,
	}
}

func (ds DiskStorage) WriteToDisk() {
	Producer, err := NewProducer(ds.fileName)
	if err != nil {
		ds.log.Info("cannot open file", zap.Error(err))
	}
	defer func() {
		err := Producer.Close()
		if err != nil {
			ds.log.Info("cannot close Producer", zap.Error(err))
		}
	}()

	for _, metric := range ds.strg.GetMetricsModels(context.Background()) {
		if err := Producer.WriteMetric(&metric); err != nil {
			ds.log.Info("cannot encode", zap.Error(err))
		}
	}
}

func (ds DiskStorage) WriteToStorage() {
	Consumer, err := NewConsumer(ds.fileName)
	if err != nil {
		ds.log.Info("cannot open file", zap.Error(err))
	}
	defer func() {
		err := Consumer.Close()
		if err != nil {
			ds.log.Info("cannot close Consumer", zap.Error(err))
		}
	}()

	for {
		metric, err := Consumer.ReadMetric()
		if err != nil {
			ds.log.Info("cannot decode", zap.Error(err))
			break
		}
		err = ds.strg.SetMetric(context.Background(), *metric)
		if err != nil {
			ds.log.Info("cannot SetMetric", zap.Error(err))
			break
		}
	}
}

func (ds DiskStorage) CanWriteToDisk() bool {
	return ds.fileName != ""
}
