package utils

import (
	"github.com/ipfs/boxo/files"
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
