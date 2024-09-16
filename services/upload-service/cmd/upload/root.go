package cmd

import (
	goversion "github.com/caarlos0/go-version"
	"github.com/spf13/cobra"
)

func Execute(version goversion.Info, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd  *cobra.Command
	exit func(int)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		panic(err)
	}
}

func newRootCmd(version goversion.Info, exit func(int)) *rootCmd {
	root := &rootCmd{exit: exit}

	cmd := &cobra.Command{
		Use:               "upload-service",
		Short:             "Run upload service.",
		Long:              "Run service for uploading texts.",
		Version:           version.String(),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}
	cmd.SetVersionTemplate("{{ .Version }}")

	root.cmd = cmd
	return root
}
