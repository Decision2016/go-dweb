/**
  @author: decision
  @date: 2024/6/6
  @note:
**/

package utils

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CreateDirIncrement(filepath string, from string, to string) ([]string,
	error) {
	repo, err := git.PlainOpen(filepath)
	if err != nil {
		return nil, err
	}
	var result []string
	var fromCommit, toCommit *object.Commit

	if from == "" && to == "" {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}

		commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
		if err != nil {
			return nil, err
		}

		var commits []*object.Commit
		err = commitIter.ForEach(func(commit *object.Commit) error {
			if len(commits) < 2 {
				commits = append(commits, commit)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		if len(commits) < 2 {
			return nil, fmt.Errorf("commit times too less < 2")
		}
		fromCommit = commits[1]
		toCommit = commits[0]
	} else {
		fromCommit, err = repo.CommitObject(plumbing.NewHash(from))
		if err != nil {
			return nil, err
		}

		toCommit, err = repo.CommitObject(plumbing.NewHash(to))
		if err != nil {
			return nil, err
		}

	}
	diff, err := toCommit.Patch(fromCommit)
	if err != nil {
		return nil, err
	}
	for _, file := range diff.Stats() {
		result = append(result, file.Name)
	}

	return result, nil
}
