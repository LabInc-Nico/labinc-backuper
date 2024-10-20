/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/labinc-nico/labinc-backuper/cmd"
	"github.com/labinc-nico/labinc-backuper/internal/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log     *zap.SugaredLogger
	version string
)

func init() {
	log = logger.GetLogger()
}

func main() {
	log.Debugf("use config file : %s", viper.ConfigFileUsed())
	cmd.Execute(GetVersion())
}

func GetVersion() string {
	if version == "" {
		version = "unknown"
	}
	return version
}
