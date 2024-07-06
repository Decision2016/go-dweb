/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"log"
)

func main() {
	//repo, err := git.PlainInit("/Users/decision/Repos/go-dweb/cmd/playground"+
	//	"/git", false)
	//if err != nil {
	//	log.Fatal(err)
	//}
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})
	err := fmt.Errorf("test error info")
	logrus.WithError(err).Errorf("print error")

	repo, err := git.PlainOpen("/Users/decision/Repos/go-dweb/cmd/playground" +
		"/git")
	if err != nil {
		log.Fatal(err)
	}

	//w, err := repo.Worktree()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//commit, err := w.Commit("commit second", &git.CommitOptions{
	//	All:       true,
	//	Author:    nil,
	//	Committer: nil,
	//	Parents:   nil,
	//	SignKey:   nil,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//obj, err := repo.CommitObject(commit)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(obj)

	// 获取 HEAD 引用
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	// 获取 HEAD 对应的 commit 对象
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}

	// 获取最新的两个 commit
	var commits []*object.Commit
	err = commitIter.ForEach(func(commit *object.Commit) error {
		if len(commits) < 2 {
			commits = append(commits, commit)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// 比较两个 commit
	fromCommit := commits[1]
	toCommit := commits[0]
	patch, err := toCommit.Patch(fromCommit)
	if err != nil {
		log.Fatal(err)
	}

	// 输出修改的文件名
	for _, file := range patch.Stats() {
		fmt.Println(file.Name)
	}
}
