/**
  @author: decision
  @date: 2024/6/7
  @note:
**/

package utils

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/sirupsen/logrus"
)

func LoadGlobalConfig(filepath string) error {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles(filepath)
	if err != nil {
		logrus.WithField("error", err).Errorln("Load config file failed.")
		return err
	}

	logrus.Infof("global config file loaded")
	return nil
}
