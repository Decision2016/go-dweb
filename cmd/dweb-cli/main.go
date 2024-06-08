/**
  @author: decision
  @date: 2024/6/5
  @note: dweb 的命令行，控制本地的 dweb 进行服务更新/部署，工作目录初始化？9
**/

package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dweb",
	Short: "DWeb is an extensible decentralized web service framework",
	Long: "DWeb is an extensible decentralized web service framework " +
		"that can be used for decentralized deployment of web applications " +
		"such as React, Vue, etc",
}

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(appCmd)
}
