package acme

import (
	"context"
	"strings"
	"testing"

	"github.com/zmoony/pi-gateway/internal/domain"
)

type runnerCall struct {
	env   map[string]string
	input string
	name  string
	args  []string
}

type runnerStub struct {
	calls []runnerCall
}

func (s *runnerStub) Run(ctx context.Context, name string, args ...string) (string, error) {
	return s.RunEnv(ctx, nil, name, args...)
}

func (s *runnerStub) RunInput(ctx context.Context, input string, name string, args ...string) (string, error) {
	return s.RunInputEnv(ctx, nil, input, name, args...)
}

func (s *runnerStub) RunEnv(_ context.Context, env map[string]string, name string, args ...string) (string, error) {
	s.calls = append(s.calls, runnerCall{env: env, name: name, args: append([]string(nil), args...)})
	return "", nil
}

func (s *runnerStub) RunInputEnv(_ context.Context, env map[string]string, input string, name string, args ...string) (string, error) {
	s.calls = append(s.calls, runnerCall{
		env:   env,
		input: input,
		name:  name,
		args:  append([]string(nil), args...),
	})
	return "", nil
}

func TestSystemManagerIssueAndInstallUsesNativeAcmeCommand(t *testing.T) {
	runner := &runnerStub{}
	manager := SystemManager{
		Runner:     runner,
		ACMEShPath: "/root/.acme.sh/acme.sh",
		ReloadCmd:  "nginx -s reload",
	}

	_, _, err := manager.IssueAndInstall(context.Background(), domain.CertificateConfig{
		RootDomain:         "example.com",
		Provider:           "aliyun",
		AccessKeyID:        "ak",
		AccessKeySecretEnc: "secret",
		InstallDir:         t.TempDir(),
	})
	if err != nil {
		t.Fatalf("IssueAndInstall() error = %v", err)
	}
	if len(runner.calls) != 2 {
		t.Fatalf("expected 2 runner calls, got %d", len(runner.calls))
	}

	for _, call := range runner.calls {
		if call.name == "cmd" {
			t.Fatalf("expected native command execution, got cmd /c call: %#v", call)
		}
	}

	issueCall := runner.calls[0]
	if issueCall.name != "/root/.acme.sh/acme.sh" {
		t.Fatalf("expected acme.sh path as command, got %s", issueCall.name)
	}
	if issueCall.env["Ali_Key"] != "ak" || issueCall.env["Ali_Secret"] != "secret" {
		t.Fatalf("expected Aliyun credentials in env, got %#v", issueCall.env)
	}
	if strings.Join(issueCall.args, " ") != "--issue --dns dns_ali -d *.example.com -d example.com" {
		t.Fatalf("unexpected issue args: %#v", issueCall.args)
	}
}

func TestSystemManagerRenewUsesNativeAcmeCommand(t *testing.T) {
	runner := &runnerStub{}
	manager := SystemManager{
		Runner:     runner,
		ACMEShPath: "/root/.acme.sh/acme.sh",
	}

	err := manager.Renew(context.Background(), domain.CertificateConfig{
		RootDomain:         "example.com",
		Provider:           "aliyun",
		AccessKeyID:        "ak",
		AccessKeySecretEnc: "secret",
	})
	if err != nil {
		t.Fatalf("Renew() error = %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected 1 runner call, got %d", len(runner.calls))
	}

	call := runner.calls[0]
	if call.name == "cmd" {
		t.Fatalf("expected native command execution, got cmd /c call: %#v", call)
	}
	if call.name != "/root/.acme.sh/acme.sh" {
		t.Fatalf("expected acme.sh path as command, got %s", call.name)
	}
	if call.env["Ali_Key"] != "ak" || call.env["Ali_Secret"] != "secret" {
		t.Fatalf("expected Aliyun credentials in env, got %#v", call.env)
	}
	if strings.Join(call.args, " ") != "--renew -d example.com --force" {
		t.Fatalf("unexpected renew args: %#v", call.args)
	}
}
