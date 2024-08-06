package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use:          "raccoon",
		Short:        "Scalable event ingestion tool",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Args: cobra.NoArgs,
		Long: heredoc.Doc(`
			Raccoon is a high-throughput, low-latency service to collect 
			events in real-time from your web, mobile apps, and services 
			using multiple network protocols.`),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	root.AddCommand(serverCommand())
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	return root
}
