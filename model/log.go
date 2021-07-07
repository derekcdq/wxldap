package model

import (
	"github.com/wonderivan/logger"
)

func InitLogConfig () {
	configFile := new(ConfigFile)
	configFile.Init("log.json","conf")
	err := logger.SetLogger(configFile.FilePath)
	if err != nil {
		logger.Info(err)
	}
}