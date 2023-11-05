package config

import (
	"fmt"
	"os"

	"github.com/pthomison/errcheck"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ViperDefaults() {
	viper.SetDefault("cache.directory", fmt.Sprintf("%v/.config/wikindexer", os.Getenv("HOME")))
	viper.SetDefault("cache.filename", "cache")
	viper.SetDefault("cores", 4)
}

func ViperInit() {
	ViperDefaults()

	viper.SetConfigName("config")                   // name of config file (without extension)
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.config/wikindexer") // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("No config file found, using default configuration")
		} else {
			errcheck.Check(err)
		}
	}
}
