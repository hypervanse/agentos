package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/NilayYadav/agentos/pkg/runtime/container"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute command inside a container",
		RunE:  runExecuteContainer,
	}

	cmd.Flags().StringP("cmd", "n", "", "Shell command to execute")
	cmd.Flags().StringP("containerId", "c", "", "Container ID")

	_ = cmd.MarkFlagRequired("containerId")
	_ = cmd.MarkFlagRequired("cmd")

	return cmd
}

func runExecuteContainer(cmd *cobra.Command, args []string) error {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmdStr, err := cmd.Flags().GetString("cmd")
	if err != nil {
		return fmt.Errorf("failed to get 'cmd' flag: %w", err)
	}
	if cmdStr == "" {
		return fmt.Errorf("cmd cannot be empty")
	}

	containerID, err := cmd.Flags().GetString("containerId")
	if err != nil {
		return fmt.Errorf("failed to get 'containerId' flag: %w", err)
	}

	manager := container.NewContainerManager(logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status, err := manager.GetContainerStatus(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to get container status: %w", err)
	}

	logger.Info().Str("containerId", containerID).Str("status", status).Msg("Container status retrieved")

	if status != "running" {
		return fmt.Errorf("container is not running (status: %s)", status)
	}

	output, err := manager.ExecuteCommand(ctx, containerID, cmdStr)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	logger.Info().Str("output", output).Msg("Command executed successfully")
	return nil
}
