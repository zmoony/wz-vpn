package service

import (
	"context"
	"testing"

	"github.com/zmoony/pi-gateway/internal/platform/conntrack"
)

type conntrackManagerStub struct {
	items []conntrack.Entry
}

func (s conntrackManagerStub) List(context.Context) ([]conntrack.Entry, error) {
	return s.items, nil
}

func TestConntrackListFiltersBySourceIP(t *testing.T) {
	service := ConntrackService{
		Manager: conntrackManagerStub{
			items: []conntrack.Entry{
				{Protocol: "tcp", SourceIP: "10.66.66.2", DestinationIP: "192.168.1.10"},
				{Protocol: "udp", SourceIP: "10.66.66.3", DestinationIP: "192.168.1.11"},
			},
		},
	}

	items, err := service.List(context.Background(), "10.66.66.2")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 || items[0].SourceIP != "10.66.66.2" {
		t.Fatalf("unexpected filtered items: %#v", items)
	}
}
