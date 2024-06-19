package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/deploy"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"time"
)

var appCmd = &cobra.Command{
	Use:              "app",
	Short:            "Decentralized application commands",
	Long:             "Decentralized application commands",
	TraverseChildren: true,
}

var appInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize decentralized application directory",
	Long:  "Initialize decentralized application directory",
	Run: func(cmd *cobra.Command, args []string) {
		deploy.AppInitial()
	},
}

var (
	appGenerateStart  string
	appGenerateEnd    string
	appGenerateOutput string
)

// env
var (
	filePath string
)

var appGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate incremental data index",
	Long: "Generate incremental data index, optionally specifying the " +
		"start commit [begin] and end commit [end], " +
		"defaulting to the most recent two commits.",
	Run: func(cmd *cobra.Command, args []string) {
		diffs, err := utils.CreateDirIncrement(filePath, appGenerateStart,
			appGenerateEnd)
		if err != nil {
			logrus.WithError(err).Errorln(
				"error occurred when creating increment file")
			return
		}

		diffBytes, err := json.Marshal(diffs)
		if err != nil {
			logrus.WithError(err).Errorln("marshal diff to json failed")
			return
		}

		err = os.WriteFile(".increment", diffBytes, 0700)
		if err != nil {
			logrus.WithError(err).Errorln("error occurred writing json to file")
			return
		}

		logrus.Infoln("increment file created to .increment")
	},
}

var appCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create new commit in workspace",
	Long:  "Create new commit in workspace",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.PlainOpen(filePath)
		if err != nil {
			logrus.Error(err)
			return
		}

		tree, err := repo.Worktree()
		if err != nil {
			logrus.Error(err)
			return
		}

		_, err = tree.Add(".")
		if err != nil {
			logrus.Error(err)
			return
		}

		now := time.Now()
		timeString := now.Format("2006-01-02 15:04:05")

		commitMsg := fmt.Sprintf("dweb commit: %s", timeString)
		commit, err := tree.Commit(commitMsg, &git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  "dweb-app",
				Email: "apps@dweb.org",
			},
		})

		if err != nil {
			logrus.Error(fmt.Errorf("create new commit failed"))
			return
		}

		logrus.Infof("new commit hash: %s", commit)
	},
}

// 部署处理流程：
// (通过 generate 命令创建) 加载 repo 目录，获取与上一次 commit 的区别
// 1. 存在差异的文件加入到处理队列，加载进度条，创建一个上传索引文件，标识每个文件的上传状态
// 2.
var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy latest decentralized application to dweb",
	Long:  "Deploy latest decentralized application to dweb",
	Run: func(cmd *cobra.Command, args []string) {
		//incrementFiles, err := utils.CreateDirIncrement()
	},
}

func init() {
	filePath = utils.GetEnvDefault("FILE_PATH", ".")

	appGenerateCmd.Flags().StringVarP(&appGenerateStart, "start", "s", "", "start commit hash")
	appGenerateCmd.Flags().StringVarP(&appGenerateEnd, "end", "e", "", "end commit hash")
	appGenerateCmd.Flags().StringVarP(&appGenerateOutput, "output", "o", "", "output file path")
	appCmd.AddCommand(appInitCmd)
	appCmd.AddCommand(appGenerateCmd)
	appCmd.AddCommand(appCommitCmd)
	appCmd.AddCommand(appDeployCmd)
}
