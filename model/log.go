package model

import (
	"github.com/wonderivan/logger"
)

func InitLogConfig () {
	configFile := new(ConfigFile)
	configFile.Init("log.json","conf")
	logger.SetLogger(configFile.FilePath)
}