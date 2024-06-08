package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/deploy"
	"github.io/decision2016/go-dweb/utils"
	"time"
)

var appCmd = &cobra.Command{
	Use:              "app",
	Long:             "Decentralized application commands",
	TraverseChildren: true,
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
		diffs, err := utils.CreateDirIncrement(appGenerateStart, appGenerateEnd)
		if err != nil {
			logrus.Error(err)
			return
		}

		fmt.Println("The directory file diff between commits are:")
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

		commitMsg := fmt.Sprintf("Dweb commit: %s", timeString)
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

		logrus.Infof("New commit hash: %s", commit)
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
}
