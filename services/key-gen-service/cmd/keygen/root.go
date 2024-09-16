package cmd

import (
	"context"
	"os/signal"
	"syscall"

	goversion "github.com/caarlos0/go-version"
	"github.com/sazonovItas/go-pastebin/pkg/config"
	"github.com/sazonovItas/go-pastebin/pkg/logger"
	"github.com/sazonovItas/go-pastebin/services/key-gen-service/internal/app"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd        *cobra.Command
	configPath string
	exit       func(int)
}

func Execute(version goversion.Info, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
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
		Use:               "key-gen-service",
		Short:             "Run ukey-gen service.",
		Long:              "Run ukey generation service for unique text's ids.",
		Version:           version.String(),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Run: func(cmd *cobra.Command, _ []string) {
			var cfg app.Config
			if err := config.Load(
				&cfg,
				root.configPath,
				true,
				config.WithEnvs("KEYGEN"),
				config.WithConfigType("yaml"),
				config.WithDefaults(map[string]any{
					"log.format":      "json",
					"log.level":       "info",
					"core.key_buffer": 10,
					"core.key_length": 15,
				}),
			); err != nil {
				panic(err)
			}

			if err := logger.ConfigureLogger(
				logger.WithEncoding(cfg.Log.Format),
				logger.WithLevel(logger.ParseLevel(cfg.Log.Level)),
			); err != nil {
				panic(err)
			}
			//nolint:errcheck
			defer logger.Sync()

			application := app.NewApp(logger.CreateLogger(), cfg)
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
