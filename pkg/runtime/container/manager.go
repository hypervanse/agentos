package container

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/containerd/go-runc"
	"github.com/google/shlex"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/rs/zerolog"
)

type Container struct {
	Id   string
	Name string
}

type ContainerManager struct {
	runc   *runc.Runc
	logger zerolog.Logger
}

func NewContainerManager(logger zerolog.Logger) *ContainerManager {
	return &ContainerManager{
		runc: &runc.Runc{
			Command: "runc",
		},
		logger: logger,
	}
}

func (m *ContainerManager) CreateContainer(ctx context.Context, name string, bundlePath string) (*Container, error) {

	containerId := GenerateContainerUID()

	m.logger.Info().Str("containerId", containerId).Str("bundlePath", bundlePath).Msg("Creating container")

	opts := &runc.CreateOpts{}

	_, err := m.runc.Run(ctx, containerId, bundlePath, opts)

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to run container")
		return nil, fmt.Errorf("failed to run container: %w", err)
	}

	return &Container{
		Id:   containerId,
		Name: name,
	}, nil
}

func (m *ContainerManager) ExecuteCommand(ctx context.Context, containerID string, command string) error {
	cmd := fmt.Sprintf("bash -c '%s'", command)
	parsedCmd, err := shlex.Split(cmd)

	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}

	execIO, err := runc.NewSTDIO()
	if err != nil {
		return fmt.Errorf("failed to create IO: %w", err)
	}

	processSpec := specs.Process{
		Args: parsedCmd,
		Env:  []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Cwd:  "/",
	}

	execOpts := &runc.ExecOpts{
		IO: execIO,
	}

	err = m.runc.Exec(ctx, containerID, processSpec, execOpts)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}

func (m *ContainerManager) GetBrowserConnectURL(ctx context.Context) (string, error) {

	cmd := exec.Command("curl", "-s", "ifconfig.me")
	ip, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to get IP: %w", err)
	}

	return fmt.Sprintf("ws://%s:3000", ip), nil
}
