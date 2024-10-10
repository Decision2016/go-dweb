package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/config/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/deploy"
	"github.io/decision2016/go-dweb/interfaces"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"time"
)

var (
	appGenerateStart  string
	appGenerateEnd    string
	appGenerateOutput string

	appDeployConfig  string
	appDeployWorkdir string
)

// vars for deploy
var (
	chainIdent   *utils.Ident
	storageIdent *utils.Ident
	index        *utils.Index
	origin       = &utils.Index{}
	storage      *interfaces.IFileStorage
	chain        *interfaces.IChain
	appDir       string
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
var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy latest decentralized application to dweb",
	Long:  "Deploy latest decentralized application to dweb",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		err := utils.LoadGlobalConfig("config.yml")
		if err != nil {
			logrus.WithError(err).Errorln("load config file failed")
			return
		}

		identStr := config.String("chain.identity", "")
		ident := utils.Ident{}
		err = ident.FromString(identStr)
		if err != nil {
			logrus.WithError(err).Errorln("load identity from config failed")
			return
		}
		chainIdent = &ident

		err = workDirInit()
		if err != nil {
			logrus.WithError(err).Errorln("work directory initialization failed")
			return
		}

		status, err := checkChainIdentity(ctx)
		if err != nil {
			logrus.WithError(err).Errorln("error occurred when checking on-chain identity")
			return
		}

		switch status {
		case onChainUpdate:
			err = checkStorageDiff(ctx)
			if err != nil {
				logrus.WithError(err).Errorln("check storage diff failed")
				return
			}
			err = processUpload(ctx)
			if err != nil {
				logrus.WithError(err).Errorln(
					"error occurred when uploading files")
				return
			}
		case onChainUpload:
			err = checkStorageDiff(ctx)
			if err != nil {
				logrus.WithError(err).Errorln("check storage diff failed")
				return
			}
			err = processUpload(ctx)
			if err != nil {
				logrus.WithError(err).Errorln(
					"error occurred when upload files")
				return
			}
		case onChainNone:
			logrus.Infof("merkle root not change, deploy canceled")
			return
		}
	},
}

func init() {
	appGenerateCmd.Flags().StringVarP(&appGenerateStart, "start", "s", "", "start commit hash")
	appGenerateCmd.Flags().StringVarP(&appGenerateEnd, "end", "e", "", "end commit hash")
	appGenerateCmd.Flags().StringVarP(&appGenerateOutput, "output", "o", "", "output file path")
	appCmd.AddCommand(appInitCmd)

	appCmd.AddCommand(appGenerateCmd)
	appCmd.AddCommand(appCommitCmd)

	appDeployCmd.Flags().StringVarP(&appDeployWorkdir, "workdir", "o",
		".deploy", "work directory path")
	appDeployCmd.Flags().StringVarP(&appDeployConfig, "config", "c",
		"config.yml", "config file path")
	appCmd.AddCommand(appDeployCmd)
}
