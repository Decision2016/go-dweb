package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckDirExistAndEmpty 判断指定目录是否存在，并且是否为空
func CheckDirExistAndEmpty(path string) (bool, error) {
	parent := filepath.Dir(path)

	// 判断父目录是否存在
	_, err := os.Stat(parent)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("parent diectory not exists")
	}

	// 判断目录是否存在
	_, err = os.Stat(path)
	// 如果目录不存在直接返回 nil，创建的工作交由上层函数
	if os.IsNotExist(err) {
		return false, nil
	}

	// 否则判断目录是否为空
	dir, err := os.Open(path)
	if err != nil {
		return true, fmt.Errorf("can not open directory")
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return true, fmt.Errorf("can not read directory")
	}

	if len(entries) != 0 {
		return true, fmt.Errorf("directory is not empty")
	}

	return true, nil
}

func CreateWorkDir() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	current := filepath.Dir(ex)

	paths := []string{
		"./.service/app",
		"./.service/index",
		"./plugins",
	}

	for _, path := range paths {
		target := filepath.Join(current, path)

		if _, err = os.Stat(target); os.IsNotExist(err) {
			err = os.MkdirAll(target, 0700)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return err
}
