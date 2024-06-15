package deploy

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"path/filepath"
)

const ignoreContent = ".appConfig\n"

func appInitial(path string) error {
	// 在这里首先判断父目录不为空，以及指定目录不存在或为空文件夹
	exists, err := utils.CheckDirExistAndEmpty(path)
	if err != nil {
		return err
	}

	var formatted string
	var dirName string
	if !exists {
		parent := filepath.Dir(path)
		dirName = filepath.Base(path)

		formatted = filepath.Join(parent, dirName)
		err := os.Mkdir(formatted, 0755)
		if err != nil {
			return fmt.Errorf("create work directory failed")
		}
	}

	repo, err := git.PlainInit(formatted, false)
	if err != nil {
		return fmt.Errorf("create new git repo failed")
	}

	ignorePath := filepath.Join(formatted, ".gitignore")
	err = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)
	if err != nil {
		return fmt.Errorf("write ignore file failed")
	}

	tree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("can not get work tree")
	}

	_, err = tree.Add(".gitignore")
	if err != nil {
		return fmt.Errorf("add file .gitignore failed")
	}

	commitMsg := "initial"
	commit, err := tree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "dweb-app",
			Email: "apps@dweb.org",
		},
	})
	if err != nil {
		return fmt.Errorf("initial commit failed")
	}

	logrus.Infoln("app initial successful!")
	logrus.Infof("Initial commit hash: %s", commit)
	return nil
}
