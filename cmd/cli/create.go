package cli

import (
	"fmt"
	"github.com/NilayYadav/agentos/pkg/runtime/container"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new container",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := container.createContainer(cmd.Flag("name").Value.String())

			if err != nil {
				return fmt.Errorf("Error creating container: %v", err)
			}

			fmt.Printf("Created container: %v\n", c)

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the container")

	return cmd
}
