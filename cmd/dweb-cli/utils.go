/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var requiredDir = []string{
	".cache",
	".cache/app",
	".cache/origin",
}

var cachedFiles = []string{}

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
