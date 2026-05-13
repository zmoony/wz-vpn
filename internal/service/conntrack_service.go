package service

import (
	"context"
	"strings"

	"github.com/zmoony/pi-gateway/internal/platform/conntrack"
)

type ConntrackService struct {
	Manager conntrack.Manager
}

func (s ConntrackService) List(ctx context.Context, sourceIP string) ([]conntrack.Entry, error) {
	items, err := s.Manager.List(ctx)
	if err != nil {
		return nil, err
	}
	sourceIP = strings.TrimSpace(sourceIP)
	if sourceIP == "" {
		return items, nil
	}
	filtered := make([]conntrack.Entry, 0, len(items))
	for _, item := range items {
		if item.SourceIP == sourceIP {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
