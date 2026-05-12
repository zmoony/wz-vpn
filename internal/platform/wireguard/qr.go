package wireguard

import (
	"context"
	"encoding/base64"

	"github.com/zmoony/pi-gateway/internal/platform"
)

type QRCodeGenerator interface {
	Encode(ctx context.Context, content string) (string, error)
}

type Base64Generator struct {
	Runner platform.CommandRunner
}

func (g Base64Generator) Encode(ctx context.Context, content string) (string, error) {
	if g.Runner == nil {
		return base64.StdEncoding.EncodeToString([]byte(content)), nil
	}

	output, err := g.Runner.Run(ctx, "qrencode", "-t", "ANSIUTF8", content)
	if err != nil {
		return base64.StdEncoding.EncodeToString([]byte(content)), nil
	}

	return base64.StdEncoding.EncodeToString([]byte(output)), nil
}
