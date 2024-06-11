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
