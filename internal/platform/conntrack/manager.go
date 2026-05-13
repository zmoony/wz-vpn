package conntrack

import (
	"context"
	"strings"

	"github.com/zmoony/pi-gateway/internal/platform"
)

type Entry struct {
	Protocol    string `json:"protocol"`
	SourceIP    string `json:"sourceIp"`
	SourcePort  string `json:"sourcePort"`
	DestinationIP   string `json:"destinationIp"`
	DestinationPort string `json:"destinationPort"`
	State       string `json:"state"`
}

type Manager interface {
	List(ctx context.Context) ([]Entry, error)
}

type SystemManager struct {
	Runner platform.CommandRunner
}

func (m SystemManager) List(ctx context.Context) ([]Entry, error) {
	output, err := m.Runner.Run(ctx, "conntrack", "-L")
	if err != nil {
		return nil, err
	}
	return Parse(output), nil
}

func Parse(output string) []Entry {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	items := make([]Entry, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if entry, ok := parseLine(line); ok {
			items = append(items, entry)
		}
	}
	return items
}

func parseLine(line string) (Entry, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return Entry{}, false
	}

	entry := Entry{
		Protocol: fields[0],
	}
	for _, field := range fields {
		switch {
		case field == "ESTABLISHED" || field == "TIME_WAIT" || field == "SYN_SENT" || field == "SYN_RECV" || field == "FIN_WAIT" || field == "CLOSE_WAIT" || field == "UNREPLIED" || field == "ASSURED":
			if entry.State == "" {
				entry.State = field
			}
		case strings.HasPrefix(field, "src=") && entry.SourceIP == "":
			entry.SourceIP = strings.TrimPrefix(field, "src=")
		case strings.HasPrefix(field, "dst=") && entry.DestinationIP == "":
			entry.DestinationIP = strings.TrimPrefix(field, "dst=")
		case strings.HasPrefix(field, "sport=") && entry.SourcePort == "":
			entry.SourcePort = strings.TrimPrefix(field, "sport=")
		case strings.HasPrefix(field, "dport=") && entry.DestinationPort == "":
			entry.DestinationPort = strings.TrimPrefix(field, "dport=")
		}
	}
	if entry.State == "" {
		entry.State = "UNKNOWN"
	}
	if entry.SourceIP == "" || entry.DestinationIP == "" {
		return Entry{}, false
	}
	return entry, true
}
