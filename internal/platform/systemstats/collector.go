package systemstats

import (
	"context"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type Collector interface {
	Collect(ctx context.Context) (domain.SystemStats, error)
}

type StaticCollector struct{}

func (StaticCollector) Collect(context.Context) (domain.SystemStats, error) {
	return domain.SystemStats{
		CPUPercent:    12.5,
		Load1:         0.14,
		Load5:         0.20,
		Load15:        0.24,
		MemoryPercent: 35.2,
		MemoryUsedMB:  712,
		MemoryTotalMB: 2048,
		TemperatureC:  46.3,
		DiskPercent:   51.4,
		DiskUsedGB:    18.7,
		DiskTotalGB:   36.4,
		Network: []domain.NetworkCounters{
			{Name: "eth0", RxBytes: 18200123, TxBytes: 9045123},
			{Name: "wg0", RxBytes: 90231, TxBytes: 108233},
		},
	}, nil
}
