package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/NilayYadav/agentos/pkg/deploy"
	"github.com/NilayYadav/agentos/pkg/image"

	"github.com/NilayYadav/agentos/pkg/runtime/container"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy app inside a container",
		RunE:  runDeployCode,
	}

	cmd.Flags().StringP("repo", "n", "", "Repository name")
	cmd.Flags().StringP("entrypoint", "c", "", "Entrypoint command")

	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("entrypoint")

	return cmd
}

func runDeployCode(cmd *cobra.Command, args []string) error {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		return fmt.Errorf("failed to get 'repo' flag: %w", err)
	}
	if repo == "" {
		return fmt.Errorf("repo cannot be empty")
	}

	entrypoint, err := cmd.Flags().GetString("entrypoint")
	if err != nil {
		return fmt.Errorf("failed to get 'entrypoint' flag: %w", err)
	}

	manager := container.NewContainerManager(logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	image, err := image.PullBrowserlessImage(ctx, logger)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to pull image")
	}

	bundlePath := fmt.Sprintf("%s/bundle", image)

	container, err := manager.CreateContainer(ctx, repo, bundlePath)

	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err = deploy.Deploy(ctx, container.Id, entrypoint); err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	return nil
}
