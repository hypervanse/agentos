package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/NilayYadav/agentos/pkg/runtime/container"
	"github.com/rs/zerolog/log"
)

func Deploy(ctx context.Context, containerId string, entrypoint string) error {

	cm := container.NewContainerManager(log.Logger)

	status, err := cm.GetContainerStatus(ctx, containerId)

	if err != nil {
		return err
	}

	if status == "running" {
		return errors.New("Container is already running")
	}

	return deployToContainer(ctx, cm, containerId, entrypoint)
}

func deployToContainer(ctx context.Context, cm *container.ContainerManager, containerId string, entrypoint string) error {

	tmpDir, err := os.MkdirTemp("", "deploy-*")

	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	repoURL := os.Getenv("REPO_URL")

	if repoURL == "" {
		return errors.New("REPO_URL is not set")
	}

	log.Info().Str("repo", repoURL).Str("dir", tmpDir).Msg("Cloning repo")

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	log.Info().Str("containerId", containerId).Str("entrypoint", entrypoint).Msg("Executing entrypoint")
	output, err := cm.ExecuteCommand(ctx, containerId, entrypoint)

	if err != nil {
		return fmt.Errorf("failed to execute entrypoint command: %w", err)
	}

	log.Info().Str("output", output).Msg("Deployment completed successfully")
	return nil
}
