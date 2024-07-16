/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/utils"
)

var rootCmd = &cobra.Command{
	Use:   "dweb",
	Short: "DWeb is an extensible decentralized web service framework",
	Long: "DWeb is an extensible decentralized web service framework " +
		"that can be used for decentralized deployment of web applications " +
		"such as React, Vue, etc",
	TraverseChildren: true,
}

func main() {
	rootCmd.Execute()
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})

	rootCmd.AddCommand(serviceCmd)
}
