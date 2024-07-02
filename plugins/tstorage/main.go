/**
  @author: decision
  @date: 2024/7/2
  @note:
**/

package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var Instance = TestStorage{}

type TestStorage struct {
}

func (t TestStorage) Initial(ctx context.Context) error {
	if _, err := os.Stat(".test"); !os.IsNotExist(err) {
		return err
	} else {
		err = os.RemoveAll(".test")
		if err != nil {
			return err
		}
	}

	err := os.Mkdir(".test", 0700)
	if err != nil {
		return err
	}

	return nil
}

func (t TestStorage) Ping(ctx context.Context) error {
	return nil
}

func (t TestStorage) Exists(ctx context.Context, src string) (bool, error) {
	return true, nil
}

func (t TestStorage) Upload(ctx context.Context, name string, src string) (string, error) {
	// 模拟上传文件到 dfs
	time.Sleep(time.Second / 100)
	cid, err := utils.GetFileCidV0(src)
	if err != nil {
		logrus.WithError(err).Debugln("calculate file cid failed")
		return "", err
	}

	// todo: 计算绝对路径
	dst := filepath.Join(utils.AbsPath(".test"), cid.String())
	cpCmd := exec.Command("cp", src, dst)
	err = cpCmd.Run()
	if err != nil {
		return "", err
	}

	return cid.String(), nil
}

func (t TestStorage) Download(ctx context.Context, identity string, dst string) error {
	src := filepath.Join(utils.AbsPath(".test"), identity)
	cpCmd := exec.Command("cp", src, dst)
	err := cpCmd.Run()

	// 模拟从 dfs 下载文件
	time.Sleep(time.Second / 100)

	if err != nil {
		logrus.WithError(err).Debugln("upload file failed")
		return err
	}

	return nil
}

func (t TestStorage) Delete(ctx context.Context, identity string) error {
	return nil
}
