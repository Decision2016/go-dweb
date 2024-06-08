package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/deploy"
	"github.io/decision2016/go-dweb/utils"
)

var appCmd = &cobra.Command{
	Use:  "app",
	Long: "Decentralized application commands",
}

var appInitCmd = &cobra.Command{
	Use:  "init",
	Long: "Initialize decentralized application directory",
	Run: func(cmd *cobra.Command, args []string) {
		deploy.AppInitial()
	},
}

var (
	appGenerateStart  string
	appGenerateEnd    string
	appGenerateOutput string
)

var appGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate incremental data index",
	Long: "Generate incremental data index, optionally specifying the " +
		"start commit [begin] and end commit [end], " +
		"defaulting to the most recent two commits.",
	Run: func(cmd *cobra.Command, args []string) {
		diffs, err := utils.CreateDirIncrement(appGenerateStart, appGenerateEnd)
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Infoln("The directory file diff between commits are:")
		for _, diff := range diffs {
			fmt.Println(diff)
		}
	},
}

var appCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create new commit in workspace",
	Long:  "Create new commit in workspace",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	appInitCmd.Flags().StringVarP(&appGenerateStart, "start", "s", "", "start commit hash")
	appInitCmd.Flags().StringVarP(&appGenerateEnd, "end", "e", "", "end commit hash")
	appInitCmd.Flags().StringVarP(&appGenerateOutput, "output", "o", "", "output file path")
	appCmd.AddCommand(appInitCmd)
	appCmd.AddCommand(appGenerateCmd)
}
