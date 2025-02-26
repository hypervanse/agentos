package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/NilayYadav/agentos/pkg/image"
	"github.com/NilayYadav/agentos/pkg/runtime/container"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new container",
		RunE:  runCreateContainer,
	}

	cmd.Flags().StringP("name", "n", "", "Name of the container")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runCreateContainer(cmd *cobra.Command, args []string) error {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get 'name' flag")
		return fmt.Errorf("failed to get 'name' flag: %w", err)
	}

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	manager := container.NewContainerManager(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	image, err := image.PullBrowserlessImage(ctx, logger)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to pull image")
	}

	bundlePath := fmt.Sprintf("%s/bundle", image)

	c, err := manager.CreateContainer(ctx, name, bundlePath)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to create container")
		return fmt.Errorf("failed to create container: %w", err)
	}

	logger.Info().Str("containerId", c.Id).Str("name", c.Name).Msg("Container created")

	connectURL, err := manager.GetBrowserConnectURL(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get browser connect URL")
	}

	logger.Info().Str("connectURL", connectURL).Msg("Connect with Playwright using")

	return nil
}
