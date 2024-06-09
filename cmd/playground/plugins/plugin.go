package main

import (
	"github.com/gookit/config/v2"
)

type Plugin struct {
	name string
}

func (p *Plugin) Load() {
	// 测试在外界加载的 config，插件加载后是否可读
	p.name = config.String("plugin.name")
}

func (p *Plugin) PluginName() string {
	return p.name
}

var Instance Plugin
