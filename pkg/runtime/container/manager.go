package container

import (
	"context"
	"fmt"

	"github.com/containerd/go-runc"
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

func (m *ContainerManager) CreateContainer(ctx context.Context, name string) (*Container, error) {

	containerId := GenerateContainerUID()

	bundleDir := "."

	m.logger.Info().Str("containerId", containerId).Str("bundleDir", bundleDir).Msg("Creating container")

	opts := &runc.CreateOpts{}

	if err := m.runc.Create(ctx, containerId, bundleDir, opts); err != nil {
		m.logger.Error().Err(err).Msg("Failed to create container")
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	m.logger.Info().Str("containerId", containerId).Msg("Starting container")

	if err := m.runc.Start(ctx, containerId); err != nil {
		m.logger.Error().Err(err).Msg("Failed to start container")
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return &Container{
		Id:   containerId,
		Name: name,
	}, nil
}
