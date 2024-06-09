package main

type IPlugin interface {
	Load()
	PluginName() string
}
