package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/go-cid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func GetUnixFsFile(path string) (files.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	f, err := files.NewReaderPathFile(path, file, stat)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GetUnixFsNode(path string) (files.Node, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, stat)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetFileCidV0 根据文件目录 path 计算文件的 cid
func GetFileCidV0(path string) (*cid.Cid, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		logrus.WithError(err).Error("read file errored")
		return nil, err
	}

	unixFSWrapped := unixfs.FilePBData(fileBytes, uint64(len(fileBytes)))

	c := merkledag.NodeWithData(unixFSWrapped).Cid()

	return &c, nil
}

// GetHeadCommit 获取 repo 的最新 commit
func GetHeadCommit(dir string) (string, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

// ListAllCommittedFiles 列出 repo 的所有已 commit 的文件
func ListAllCommittedFiles(dir string) ([]string, error) {
	var results []string
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	tree.Files().ForEach(func(f *object.File) error {
		results = append(results, f.Name)
		return nil
	})

	return results, nil
}

// CreateDirIndex 创建目录索引，path 通过类 hash 的 cid 标识
func CreateDirIndex(dir string) (*Index, error) {
	commitStr, err := GetHeadCommit(dir)
	if err != nil {
		return nil, err
	}

	fileList, err := ListAllCommittedFiles(dir)
	if err != nil {
		return nil, err
	}

	var full Index
	full.Commit = commitStr
	full.Paths = make(map[string]string)

	for _, filename := range fileList {
		absPath := filepath.Join(dir, filename)
		if filename == ".gitignore" {
			continue
		}

		c, err := GetFileCidV0(absPath)
		if err != nil {
			logrus.WithError(err).Debugf("get file %s cid failed", filename)
			return nil, err
		}

		full.Paths[filename] = c.String()
	}

	full.MerkleRoot()

	return &full, nil
}

func AbsPath(relative string) string {
	ex, _ := os.Executable()
	return filepath.Join(filepath.Dir(ex), relative)
}
