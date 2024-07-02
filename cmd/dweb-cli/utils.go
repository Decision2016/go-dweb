/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var requiredDir = []string{
	".cache",
	".cache/app",
	".cache/origin",
}

var cachedFiles = []string{
	".cache/.index",    // 临时 index 文件
	".cache/.progress", // 上传进度保存文件
	".cache/.origin",   // 链上存放拉取到的 index 文件
}

// workDirInitial 初始化工作目录
func workDirInitial() error {
	ex, err := os.Executable()
	if err != nil {
		logrus.WithError(err).Errorln("get executable path failed")
		return err
	}

	runDir := filepath.Dir(ex)
	for _, dir := range requiredDir {
		abs := filepath.Join(runDir, dir)
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			err = os.Mkdir(abs, 0700)
			if err != nil {
				logrus.WithError(err).Debugln("mkdir %s failed", dir)
				return err
			}
		} else if err != nil {
			logrus.WithError(err).Debugln("get %s stat failed", dir)
			return err
		}
	}

	return nil
}

func cleanCachedFiles() error {
	ex, err := os.Executable()
	if err != nil {
		logrus.WithError(err).Errorln("get executable path failed")
		return err
	}

	runDir := filepath.Dir(ex)
	for _, file := range cachedFiles {
		abs := filepath.Join(runDir, file)
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			continue
		} else if err != nil {
			logrus.WithError(err).Debugln("get %s stat failed", file)
			return err
		}

		err = os.Remove(abs)
		if err != nil {
			logrus.WithError(err).Debugln("delete file %s failed", file)
			return err
		}
	}

	return nil
}

func selectCleanCache() error {
	ex, err := os.Executable()
	if err != nil {
		logrus.WithError(err).Errorln("get executable path failed")
		return err
	}

	exists := false
	var opt string
	runDir := filepath.Dir(ex)
	for _, file := range cachedFiles {
		abs := filepath.Join(runDir, file)
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			continue
		} else if err != nil {
			logrus.WithError(err).Debugln("get %s stat failed", file)
			return err
		}

		exists = true
		break
	}

	if exists {
		fmt.Printf("found uploaded cache data, delete? (Y/N):")
		_, err := fmt.Scan(&opt)
		if err != nil {
			logrus.Debugln("read option string failed")
			return err
		}
	}

	if opt == "Y" {
		err = cleanCachedFiles()
		if err != nil {
			return err
		}
	}

	return nil
}
