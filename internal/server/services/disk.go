package services

import (
	"context"
	"time"
)

type DiskStorage interface {
	WriteToDisk()
	WriteToStorage()
	CanWriteToDisk() bool
}

type DiskService struct {
	diskStrg      DiskStorage
	storeInterval int
	restore       bool
}

func NewDiskService(diskStorage DiskStorage, storeInterval int, restore bool) *DiskService {
	return &DiskService{
		diskStrg:      diskStorage,
		storeInterval: storeInterval,
		restore:       restore,
	}
}

func (dService *DiskService) CollectMetrics(ctx context.Context) {
	if dService.diskStrg.CanWriteToDisk() {
		ticker := time.NewTicker(time.Duration(dService.storeInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				dService.diskStrg.WriteToDisk()
			case <-ctx.Done():
				dService.diskStrg.WriteToDisk()
				return
			}
		}
	}
}

func (dService DiskService) FillMetricStorage() {
	if dService.restore {
		dService.diskStrg.WriteToStorage()
	}
}
