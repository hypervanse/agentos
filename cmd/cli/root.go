package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agentos",
	Short: "AgentOS is a environment for running agents",
}

func Execute() error {
	return rootCmd.Execute()
}

func InitializeCli() error {

	err := Execute()

	if err != nil {
		return fmt.Errorf("error executing root command: %v", err)
	}

	initializeCmds()

	return nil
}

func initializeCmds() {
	rootCmd.AddCommand(newCreateCmd())
}
