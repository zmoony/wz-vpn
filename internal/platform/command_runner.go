package platform

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
	RunInput(ctx context.Context, input string, name string, args ...string) (string, error)
	RunEnv(ctx context.Context, env map[string]string, name string, args ...string) (string, error)
	RunInputEnv(ctx context.Context, env map[string]string, input string, name string, args ...string) (string, error)
}

type ExecCommandRunner struct{}

func (r ExecCommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	return r.RunInputEnv(ctx, nil, "", name, args...)
}

func (r ExecCommandRunner) RunInput(ctx context.Context, input string, name string, args ...string) (string, error) {
	return r.RunInputEnv(ctx, nil, input, name, args...)
}

func (r ExecCommandRunner) RunEnv(ctx context.Context, env map[string]string, name string, args ...string) (string, error) {
	return r.RunInputEnv(ctx, env, "", name, args...)
}

func (r ExecCommandRunner) RunInputEnv(ctx context.Context, env map[string]string, input string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if input != "" {
		cmd.Stdin = bytes.NewBufferString(input)
	}
	if len(env) > 0 {
		cmd.Env = os.Environ()
		for key, value := range env {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %v: %w: %s", name, args, err, stderr.String())
	}

	return stdout.String(), nil
}
