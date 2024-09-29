package app

import (
	"github.com/sazonovItas/go-pastebin/pkg/app"
	"go.uber.org/zap"
)

type Config struct{}

func New(log *zap.Logger, cfg Config) *app.App {
	return &app.App{
		Cfg:      cfg,
		CleanUps: []func(){},
	}
}
