package nginx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform"
)

type Manager interface {
	Apply(ctx context.Context, host domain.ProxyHost) error
	Remove(ctx context.Context, serverName string) error
	Render(host domain.ProxyHost) string
}

type SystemManager struct {
	Runner   platform.CommandRunner
	SitesDir string
}

var invalidFileChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func (m SystemManager) Apply(ctx context.Context, host domain.ProxyHost) error {
	if err := os.MkdirAll(m.SitesDir, 0o755); err != nil {
		return fmt.Errorf("create nginx sites dir: %w", err)
	}

	target := filepath.Join(m.SitesDir, fileName(host.ServerName))
	content := m.Render(host)

	previous, hadPrevious, err := readIfExists(target)
	if err != nil {
		return err
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write nginx config: %w", err)
	}

	if _, err := m.Runner.Run(ctx, "nginx", "-t"); err != nil {
		_ = rollbackFile(target, hadPrevious, previous)
		return fmt.Errorf("nginx config test failed: %w", err)
	}
	if _, err := m.Runner.Run(ctx, "nginx", "-s", "reload"); err != nil {
		_ = rollbackFile(target, hadPrevious, previous)
		return fmt.Errorf("nginx reload failed: %w", err)
	}

	return nil
}

func (m SystemManager) Remove(ctx context.Context, serverName string) error {
	target := filepath.Join(m.SitesDir, fileName(serverName))
	previous, hadPrevious, err := readIfExists(target)
	if err != nil {
		return err
	}
	if hadPrevious {
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("remove nginx config: %w", err)
		}
	}
	if _, err := m.Runner.Run(ctx, "nginx", "-t"); err != nil {
		_ = rollbackFile(target, hadPrevious, previous)
		return fmt.Errorf("nginx config test failed: %w", err)
	}
	if _, err := m.Runner.Run(ctx, "nginx", "-s", "reload"); err != nil {
		_ = rollbackFile(target, hadPrevious, previous)
		return fmt.Errorf("nginx reload failed: %w", err)
	}
	return nil
}

func (m SystemManager) Render(host domain.ProxyHost) string {
	lines := []string{
		"server {",
		"    listen 80;",
		"    listen [::]:80;",
		"    listen 443 ssl http2;",
		"    listen [::]:443 ssl http2;",
		"    server_name " + host.ServerName + ";",
		"",
		"    ssl_certificate " + host.CertificateCertPath + ";",
		"    ssl_certificate_key " + host.CertificateKeyPath + ";",
		"",
		"    location / {",
		"        proxy_pass " + host.UpstreamURL + ";",
		"        proxy_http_version 1.1;",
		"        proxy_set_header Host $host;",
		"        proxy_set_header X-Real-IP $remote_addr;",
		"        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;",
		"        proxy_set_header X-Forwarded-Proto $scheme;",
		"    }",
		"}",
	}

	return strings.Join(lines, "\n") + "\n"
}

func fileName(serverName string) string {
	safe := invalidFileChars.ReplaceAllString(serverName, "-")
	return safe + ".conf"
}

func readIfExists(path string) ([]byte, bool, error) {
	content, err := os.ReadFile(path)
	if err == nil {
		return content, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("read file %s: %w", path, err)
}

func rollbackFile(path string, hadPrevious bool, previous []byte) error {
	if hadPrevious {
		return os.WriteFile(path, previous, 0o644)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
