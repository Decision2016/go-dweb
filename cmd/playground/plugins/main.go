package main

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/sirupsen/logrus"
	"log"
	"plugin"
)

func main() {
	LoadGlobalConfig("./config.yml")
	p, err := plugin.Open("plugin.so")
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := p.Lookup("Instance")
	if err != nil {
		log.Fatal(err)
	}

	pl := symbol.(IPlugin)

	pl.Load()
	fmt.Println(pl.PluginName())
}

func LoadGlobalConfig(filepath string) {
	config.WithOptions(config.ParseEnv)

	config.AddDriver(yaml.Driver)

	err := config.LoadFiles(filepath)
	if err != nil {
		logrus.WithField("error", err).Errorln("Load config file failed.")
	}
}
