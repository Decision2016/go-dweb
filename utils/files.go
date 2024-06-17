package utils

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
	"os"
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

func GetFileCidV0(path string) (*cid.Cid, error) {
	f, err := os.Open(path)
	if err != nil {
		logrus.Debugf("read file %s failed", path)
		return nil, err
	}
	defer f.Close()
	var fileBytes []byte
	_, err = f.Read(fileBytes)
	if err != nil {
		logrus.WithError(err).Error("read file errored")
		return nil, err
	}

	prefix := cid.Prefix{
		Version:  0,
		Codec:    uint64(multicodec.Raw),
		MhType:   multihash.SHA2_256,
		MhLength: -1,
	}

	c, err := prefix.Sum(fileBytes)
	if err != nil {
		logrus.WithError(err).Errorf("calculate cid v0 for file %s failed",
			path)

		return nil, err
	}

	return &c, nil
}

func ListAllCommittedFiles(dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	tree.Files().ForEach(func(f *object.File) error {
		fmt.Println(f.Name)
		return nil
	})

	return nil
}
