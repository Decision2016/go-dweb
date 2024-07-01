package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/config/v2"
	"github.com/schollz/progressbar/v3"
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
var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy latest decentralized application to dweb",
	Long:  "Deploy latest decentralized application to dweb",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: 这部分的实现较为复杂，需要合理地划分为多个函数
		// 1. load chain plugin by config file
		ctx := context.Background()
		err := utils.LoadGlobalConfig("config.yml")
		if err != nil {
			logrus.WithError(err).Errorln("load config file failed")
			return
		}

		var chainIdent, storageIdent utils.Ident
		chainPath := config.String("chain", "")
		err = chainIdent.FromString(chainPath)
		if chainPath == "" || err != nil {
			logrus.Errorln("load chain config failed")
			return
		}

		chain, err := utils.ParseOnChain(chainPath)
		if err != nil {
			logrus.WithError(err).Errorln("load chain plugin failed")
			return
		}

		// 2. get on-chain storage identity
		onChainStorageIdent, err := (*chain).Identity()
		if err != nil {
			logrus.WithError(err).Errorln("read on-chain storage identity failed")
			return
		}

		err = storageIdent.FromString(onChainStorageIdent)
		if err != nil {
			logrus.WithError(err).Errorln("convert identity to object failed")
			return
		}

		// 3. calculate local app repo index and merkle && check merkle equal
		//incrementFiles, err := utils.CreateDirIncrement()
		index, err := utils.CreateDirIndex(filePath)
		if err != nil {
			logrus.WithError(err).Infof("create full directory index failed")
			return
		}
		index.MerkleRoot()

		if index.Root[:8] == storageIdent.Merkle {
			logrus.Infof("merkle root equal, deploy process exit")
			return
		}

		// 4. download original files and check diff
		indexIdent, storage, err := utils.ParseFileStorage(ctx,
			onChainStorageIdent)
		if err != nil {
			logrus.WithError(err).Errorln("parse storage ident failed")
			return
		}
		// todo: create .cache file if not exist and delete if exists
		err = (*storage).Download(ctx, indexIdent, ".origin")
		if err != nil {
			logrus.WithError(err).Errorln("download original index failed")
		}

		origin, err := utils.LoadIndex(".origin")
		if err != nil {
			logrus.WithError(err).Errorln("load origin file failed")
			return
		}

		diffBar := progressbar.Default(int64(len(index.Paths)))
		uploadList := make([]string, 0)
		for k, _ := range index.Paths {
			p, ok := origin.Paths[k]
			if ok {
				err := (*storage).Download(ctx, p, ".compare")
				if err != nil {
					logrus.WithError(err).Errorln("download origin file failed")
					return
				}

				originCid, err := utils.GetFileCidV0(".compare")
				if err != nil {
					logrus.WithError(err).Errorln(
						"calculate cid for original file failed")
					return
				}

				currentCid, err := utils.GetFileCidV0(k)
				if err != nil {
					logrus.WithError(err).Errorln(
						"calculate cid for current file failed")
					return
				}

				if originCid.String() != currentCid.String() {
					uploadList = append(uploadList, k)
				}
			}

			err = diffBar.Add(1)
			if err != nil {
				logrus.WithError(err).Errorln("progress bar inc failed")
			}
		}
		//indexBytes, err := yaml.Marshal(index)
		//if err != nil {
		//	logrus.WithError(err).Errorln("marshal index to bytes failed")
		//	return
		//}
		//
		//err = os.WriteFile(".index", indexBytes, 0700)
		//if err != nil {
		//	logrus.WithError(err).Errorln("write bytes to file failed")
		//}

		//storagePath := config.String("deploy.plugin.storage", "")
		//if storagePath == "" {
		//	logrus.Error("config item 'deploy.plugin.storage' is empty")
		//	return
		//}
		//
		//symbol, err := utils.LoadSymbol(storagePath)
		//if err != nil {
		//	logrus.WithError(err).Errorln("load plugin failed")
		//}
		//
		//storage, ok := symbol.(interfaces.IFileStorage)
		//if !ok {
		//	logrus.Error("convert symbol to storage interface failed")
		//	return
		//}
		//
		//uploader := managers.NewUploader()
		//uploader.Setup(index.Commit, index.Paths)

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
