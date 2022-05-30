/*
Copyright © 2020 Marvin
*/
package service

import (
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
)

// 版本信息
var (
	Version   = "V2022-05-23"
	BuildTS   = "2022-05-23"
	GitHash   = "None"
	GitBranch = "Lixora"
	GitOwner  = "HuangLinjie-17767151782"
)

func GetAppVersion(version bool) {
	if version {
		fmt.Printf("%v", getRawVersion())
		os.Exit(1)
	}
}

// 版本信息输出重定向到日志
func RecordAppVersion(app string, logger *zap.Logger, cfg *CfgFile) {
	logger.Info("Welcome to "+app,
		zap.String("Release Version", Version),
		zap.String("Git Commit Hash", GitHash),
		zap.String("Git Branch", GitBranch),
		zap.String("UTC Build Time", BuildTS),
		zap.String("Go Version", runtime.Version()),
		zap.String("Release Version", GitOwner),
	)
	logger.Info(app+" config", zap.Stringer("config", cfg))
}

func getRawVersion() string {
	info := ""
	info += fmt.Sprintf("Release Version: %s\n", Version)
	info += fmt.Sprintf("Release Owner: %s\n", GitOwner)
	info += fmt.Sprintf("Git Commit Hash: %s\n", GitHash)
	info += fmt.Sprintf("Git Branch: %s\n", GitBranch)
	info += fmt.Sprintf("UTC Build Time: %s\n", BuildTS)
	info += fmt.Sprintf("Go Version: %s\n", runtime.Version())
	return info
}
