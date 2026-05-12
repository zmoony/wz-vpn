package service

import (
	"context"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform/systemstats"
)

type SystemStatsService struct {
	Collector systemstats.Collector
}

func (s SystemStatsService) Get(ctx context.Context) (domain.SystemStats, error) {
	return s.Collector.Collect(ctx)
}
