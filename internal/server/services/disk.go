package services

import "time"

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

func (dService DiskService) CollectMetrics() {
	if dService.diskStrg.CanWriteToDisk() {
		for {
			dService.diskStrg.WriteToDisk()

			time.Sleep(time.Duration(dService.storeInterval) * time.Second)
		}
	}
}

func (dService DiskService) FillMetricStorage() {
	if dService.restore {
		dService.diskStrg.WriteToStorage()
	}
}
