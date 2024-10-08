/**
  @author: decision
  @date: 2024/6/5
  @note: dweb 的命令行，控制本地的 dweb 进行服务更新/部署，工作目录初始化？9
**/

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/utils"
)

// env
var (
	filePath string
)

var rootCmd = &cobra.Command{
	Use:   "dweb-cli",
	Short: "DWeb is an extensible decentralized web service framework",
	Long: "DWeb is an extensible decentralized web service framework " +
		"that can be used for decentralized deployment of web applications " +
		"including React, Vue, etc",
	TraverseChildren: true,
}

func main() {
	rootCmd.Execute()
}

func init() {
	filePath = utils.GetEnvDefault("FILE_PATH", ".")
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})

	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(testCmd)
}
