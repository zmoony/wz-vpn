package acme

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/platform"
)

type Manager interface {
	IssueAndInstall(ctx context.Context, item domain.CertificateConfig) (fullchainPath string, keyPath string, err error)
	Renew(ctx context.Context, item domain.CertificateConfig) error
	ReadInstalledStatus(ctx context.Context, item domain.CertificateConfig) (domain.CertificateConfig, error)
}

type SystemManager struct {
	Runner      platform.CommandRunner
	ACMEShPath  string
	ReloadCmd   string
}

func (m SystemManager) IssueAndInstall(ctx context.Context, item domain.CertificateConfig) (string, string, error) {
	fullchainPath := filepath.Join(item.InstallDir, "fullchain.pem")
	keyPath := filepath.Join(item.InstallDir, "privkey.pem")

	if err := os.MkdirAll(item.InstallDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create cert dir: %w", err)
	}

	env := m.acmeEnv(item)
	if _, err := m.Runner.RunEnv(ctx, env, m.ACMEShPath, "--issue", "--dns", dnsProvider(item.Provider), "-d", "*."+item.RootDomain, "-d", item.RootDomain); err != nil {
		return "", "", err
	}
	installArgs := []string{
		"--install-cert",
		"-d", item.RootDomain,
		"--key-file", keyPath,
		"--fullchain-file", fullchainPath,
	}
	if strings.TrimSpace(m.ReloadCmd) != "" {
		installArgs = append(installArgs, "--reloadcmd", m.ReloadCmd)
	}
	if _, err := m.Runner.RunEnv(ctx, env, m.ACMEShPath, installArgs...); err != nil {
		return "", "", err
	}

	return fullchainPath, keyPath, nil
}

func (m SystemManager) Renew(ctx context.Context, item domain.CertificateConfig) error {
	env := m.acmeEnv(item)
	_, err := m.Runner.RunEnv(ctx, env, m.ACMEShPath, "--renew", "-d", item.RootDomain, "--force")
	return err
}

func (m SystemManager) ReadInstalledStatus(_ context.Context, item domain.CertificateConfig) (domain.CertificateConfig, error) {
	if item.FullchainPath == "" {
		return item, nil
	}
	content, err := os.ReadFile(item.FullchainPath)
	if err != nil {
		if os.IsNotExist(err) {
			return item, nil
		}
		return item, fmt.Errorf("read installed certificate: %w", err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return item, fmt.Errorf("decode pem certificate: invalid pem")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return item, fmt.Errorf("parse x509 certificate: %w", err)
	}
	notAfter := certificate.NotAfter
	item.NotAfter = &notAfter
	item.DaysRemaining = int(time.Until(notAfter).Hours() / 24)
	return item, nil
}

func (m SystemManager) acmeEnv(item domain.CertificateConfig) map[string]string {
	provider := strings.ToLower(strings.TrimSpace(item.Provider))
	if provider == "" {
		provider = "aliyun"
	}
	if provider == "aliyun" {
		return map[string]string{
			"Ali_Key":    item.AccessKeyID,
			"Ali_Secret": item.AccessKeySecretEnc,
		}
	}
	return nil
}

func dnsProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "aliyun":
		return "dns_ali"
	default:
		return "dns_ali"
	}
}
