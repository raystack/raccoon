package cmd

import "github.com/spf13/cobra"

func serverCommand() *cobra.Command {
	// todo
	return &cobra.Command{
		Use:   "server",
		Short: "Start raccoon server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
