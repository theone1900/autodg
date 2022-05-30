
package service

import (
	"strings"

	"github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// 初始化日志 logger
func NewZapLogger(cfg *CfgFile) (err error) {
	Logger, _, err = log.InitLogger(&log.Config{
		Level: strings.ToLower(cfg.LogConfig.LogLevel),
		File: log.FileLogConfig{
			Filename:   cfg.LogConfig.LogFile,
			MaxSize:    cfg.LogConfig.MaxSize,
			MaxDays:    cfg.LogConfig.MaxDays,
			MaxBackups: cfg.LogConfig.MaxBackups,
		},
	}, zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		return err
	}
	return
}
