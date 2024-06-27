/**
  @author: decision
  @date: 2024/6/27
  @note:
**/

package main

import (
	"flag"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "c", "./conf.yml", "Config file path")
}

func workDirInitial() {

}
