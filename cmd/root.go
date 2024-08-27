package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cmdx"
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
		Annotations: map[string]string{
			"group": "core",
			"help:learn": heredoc.Doc(`
				Use 'raccoon <command> --help' for more information about a command.
				Read the manual at https://raystack.github.io/raccoon/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/raystack/raccoon/issues
			`),
		},
	}

	cmdx.SetHelp(root)
	root.AddCommand(cmdx.SetCompletionCmd("raccoon"))
	root.AddCommand(cmdx.SetRefCmd(root))

	root.AddCommand(serverCommand())
	return root
}
