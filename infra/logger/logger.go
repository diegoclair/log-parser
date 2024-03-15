package logger

import (
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/log-parser/infra/config"
)

func New(cfg *config.Config) logger.Logger {
	params := logger.LogParams{
		AppName:    cfg.AppName,
		DebugLevel: cfg.LogDebug,
	}

	return logger.New(params)
}
