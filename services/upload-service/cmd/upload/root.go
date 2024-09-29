package cmd

import (
	"context"
	"os/signal"
	"syscall"

	goversion "github.com/caarlos0/go-version"
	"github.com/sazonovItas/go-pastebin/pkg/config"
	"github.com/sazonovItas/go-pastebin/pkg/logger"
	"github.com/sazonovItas/go-pastebin/services/upload-service/internal/app"
	"github.com/spf13/cobra"
)

func Execute(version goversion.Info, args []string) {
	newRootCmd(version).Execute(args)
}

type rootCmd struct {
	cmd        *cobra.Command
	configPath string
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		panic(err)
	}
}

func newRootCmd(version goversion.Info) *rootCmd {
	root := &rootCmd{}

	cmd := &cobra.Command{
		Use:               "upload-service",
		Short:             "Run upload service.",
		Long:              "Run service for uploading texts.",
		Version:           version.String(),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Run: func(_ *cobra.Command, _ []string) {
			var cfg app.Config
			if err := config.Load(
				&cfg,
				root.configPath,
				true,
				config.WithEnvs("UPLOAD"),
				config.WithConfigType("yaml"),
				config.WithDefaults(map[string]any{
					"log.format": "json",
					"log.level":  "info",
				}),
			); err != nil {
				panic(err)
			}

			if err := logger.ConfigureLogger(
				logger.WithEncoding(""),
				logger.WithLevel(logger.ParseLevel("")),
			); err != nil {
				panic(err)
			}
			//nolint:errcheck
			defer logger.Sync()

			application := app.New(logger.CreateLogger(), cfg)
			defer application.CleanUp()

			ctx, stop := signal.NotifyContext(
				context.Background(),
				syscall.SIGTERM,
				syscall.SIGINT,
				syscall.SIGQUIT,
			)
			defer stop()

			application.MustRun(ctx)
		},
	}
	cmd.SetVersionTemplate("{{ .Version }}")

	cmd.PersistentFlags().
		StringVarP(&root.configPath, "config", "c", "", "Specify path to config file.")

	root.cmd = cmd
	return root
}
