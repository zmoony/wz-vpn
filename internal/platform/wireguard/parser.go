package wireguard

import (
	"strconv"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
)

func ParseDump(output string) []domain.WireGuardRuntimePeer {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return nil
	}

	peers := make([]domain.WireGuardRuntimePeer, 0, len(lines)-1)
	for _, line := range lines[1:] {
		fields := strings.Split(line, "\t")
		if len(fields) < 8 {
			continue
		}

		var handshake *time.Time
		if fields[5] != "0" {
			if ts, err := strconv.ParseInt(fields[5], 10, 64); err == nil && ts > 0 {
				t := time.Unix(ts, 0)
				handshake = &t
			}
		}

		rx, _ := strconv.ParseInt(fields[6], 10, 64)
		tx, _ := strconv.ParseInt(fields[7], 10, 64)

		peers = append(peers, domain.WireGuardRuntimePeer{
			PublicKey:       fields[0],
			LastHandshakeAt: handshake,
			RxBytes:         rx,
			TxBytes:         tx,
			Online:          handshake != nil,
		})
	}

	return peers
}
