package container

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/beam-cloud/go-runc"
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

func (m ContainerManager) CreateContainer(ctx context.Context, name string, bundlePath string) (Container, error) {
	containerId := GenerateContainerUID()
	m.logger.Info().Str("containerId", containerId).Str("bundlePath", bundlePath).Msg("Creating container")

	opts := &runc.CreateOpts{}

	go func() {
		_, err := m.runc.Run(ctx, containerId, bundlePath, opts)
		if err != nil {
			m.logger.Error().Err(err).Msg("Failed to run container")
		}
	}()

	time.Sleep(2 * time.Second)

	m.logger.Info().Str("containerId", containerId).Msg("Container created")

	return Container{
		Id:   containerId,
		Name: name,
	}, nil
}

func (m *ContainerManager) ExecuteCommand(ctx context.Context, containerID string, command string) (string, error) {
	cmd := fmt.Sprintf("bash -c '%s'", command)
	parsedCmd, err := shlex.Split(cmd)

	if err != nil {
		return "", fmt.Errorf("failed to parse command: %w", err)
	}

	processSpec := specs.Process{
		Args: parsedCmd,
		Env:  []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Cwd:  "/",
	}

	var outputBuffer bytes.Buffer
	outputWriter := os.Stdout

	execOpts := &runc.ExecOpts{
		OutputWriter: outputWriter,
	}

	fmt.Printf("Executing '%s' in container %s\n", command, containerID)

	err = m.runc.Exec(ctx, containerID, processSpec, execOpts)

	if err != nil {
		fmt.Printf("Exec error: %v\n", err)
	}

	return outputBuffer.String(), nil
}

func (m *ContainerManager) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	m.logger.Info().Str("containerId", containerID).Msg("Checking container status")

	status, err := m.runc.State(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to get container status: %w", err)
	}

	m.logger.Info().Str("containerId", containerID).Str("status", status.Status).Msg("Container status retrieved")
	return status.Status, nil
}

func (m *ContainerManager) GetBrowserConnectURL(ctx context.Context, containerID string) (string, error) {
	ip, err := m.ExecuteCommand(ctx, containerID, "curl -s ifconfig.me")
	if err != nil {
		return "", fmt.Errorf("failed to get IP: %w", err)
	}

	return fmt.Sprintf("ws://%s:3000", ip), nil
}
