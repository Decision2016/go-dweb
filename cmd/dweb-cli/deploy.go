/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package main

import (
	"context"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/interfaces"
	"github.io/decision2016/go-dweb/managers"
	"github.io/decision2016/go-dweb/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"time"
)

func workDirInit() error {
	// 假定在这之前 identity 已经检查过了
	uid, _ := chainIdent.Uid()
	appDir = filepath.Join(appDeployWorkdir, uid)
	opt := "N"

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, 0700)
		if err != nil {
			logrus.WithError(err).Debugln("create all directory failed")
			return err
		}
		logrus.Debugf("create directory with path %s", appDeployWorkdir)
	} else if err != nil {
		return err
	}

	d, err := os.Open(appDir)
	if err != nil {
		logrus.WithError(err).Debugf("open dir %s failed", appDir)
		return err
	}
	defer d.Close()

	entities, err := d.ReadDir(0)
	if err != nil {
		logrus.WithError(err).Debugf("cann not read application cache dir")
		return err
	}

	if len(entities) != 0 {
		fmt.Printf("found cached file in %s, continue to upload? (Y/N): ")
		_, err = fmt.Scan(&opt)
		if err != nil {
			logrus.WithError(err).Debugf("read user option failed")
			return err
		}
	}

	if opt == "N" {
		err = os.RemoveAll(appDir)
		if err != nil {
			logrus.WithError(err).Debugf("remove files from %s failed", appDir)
			return err
		}

		err = os.MkdirAll(appDir, 0700)
		if err != nil {
			logrus.WithError(err).Debugf("create application dir %s failed", appDir)
			return err
		}

		sourceDir := filepath.Join(appDir, "source")
		err = os.MkdirAll(sourceDir, 0700)
		if err != nil {
			logrus.WithError(err).Debugf("create source dir failed")
			return err
		}

		originDir := filepath.Join(appDir, "origin")
		err = os.MkdirAll(originDir, 0700)
		if err != nil {
			logrus.WithError(err).Debugf("create source dir failed")
			return err
		}
	}

	return nil
}

func checkChainIdentity(ctx context.Context) (onChainStatus, error) {
	var err error
	chain, err = utils.ParseOnChain(chainIdent)
	if err != nil {
		logrus.WithError(err).Errorln("load chain plugin failed")
		return onChainDefault, err
	}

	chainStorageIdent, err := (*chain).Identity()
	if err != nil {
		logrus.WithError(err).Errorln("get storage ident failed")
		return onChainDefault, err
	}

	index, err = utils.CreateDirIndex(filePath)
	if err != nil {
		logrus.WithError(err).Errorln("create full directory index failed")
		return onChainDefault, err
	}
	index.MerkleRoot()

	storageIdent = &utils.Ident{}
	err = storageIdent.FromString(chainStorageIdent)
	if err != nil {
		logrus.Infof("on-chain identity invalid, first time upload")
		return onChainUpload, nil
	}

	if index.Root[:8] == storageIdent.Merkle {
		logrus.Infof("merkle root equal, deploy process exit")
		return onChainNone, fmt.Errorf("merkle equal")
	}

	var originIndex string
	originIndexPath := filepath.Join(appDir, ".origin")
	originIndex, storage, err = utils.ParseFileStorage(ctx, chainStorageIdent)
	if err != nil {
		logrus.WithError(err).Errorln("parse storage identity failed")
		return onChainDefault, err
	}

	err = (*storage).Download(ctx, originIndex, originIndexPath)
	if err != nil {
		logrus.WithError(err).Errorln("download origin index failed")
		return onChainDefault, err
	} else {
		origin, err = utils.LoadIndex(originIndexPath)
		if err != nil {
			return onChainDefault, err
		}
	}

	return onChainUpdate, nil
}

func checkStorageDiff(ctx context.Context) error {
	originPath := filepath.Join(appDir, "origin")
	diffBar := progressbar.Default(int64(len(index.Paths)))
	for k, v := range index.Paths {
		p, ok := origin.Paths[k]
		if ok {
			downloadPath := filepath.Join(originPath, p)
			// todo: 如果有文件占用则删除，这个逻辑交由插件实现？

			err := (*storage).Download(ctx, p, downloadPath)
			if err != nil {
				logrus.WithError(err).Errorln("download origin file failed")
				return err
			}

			originCid, err := utils.GetFileCidV0(downloadPath)
			if err != nil {
				logrus.WithError(err).Errorln("calculate file failed")
				return err
			}

			if originCid.String() != v {
				// 通过设置 path 为空标识上传文件
				index.Paths[k] = ""
			} else {
				index.Paths[k] = p
			}

		} else {
			index.Paths[k] = ""
		}

		err := diffBar.Add(1)
		if err != nil {
			logrus.WithError(err).Errorln("progress bar increase failed")
			return err
		}
	}

	return nil
}

// processUpload 上传文件到 DFS 并更新链上索引
func processUpload(ctx context.Context) error {
	logrus.Infof("current timestamp t0 = %d (ms)", time.Now().UnixMilli())
	indexPath := filepath.Join(appDir, ".index")
	progressPath := filepath.Join(appDir, ".archive")
	storagePath := config.String("plugins.storage", "")

	if storagePath == "" {
		logrus.Error("config item 'plugins.storage' is empty")
		return fmt.Errorf("config item not exists")
	}

	symbol, err := utils.LoadSymbol(storagePath)
	if err != nil {
		logrus.WithError(err).Errorln("load plugin failed")
		return err
	}

	s, ok := symbol.(interfaces.IFileStorage)
	if !ok {
		logrus.Error("convert symbol to storage interface failed")
		return fmt.Errorf("convert symbol to interface failed")
	}

	err = s.Initial(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("storage plugin initial failed")
		return err
	}

	uploader := managers.NewUploader(progressPath)
	err = uploader.Setup(index, &s)
	if err != nil {
		logrus.WithError(err).Errorln("setup uploader index failed")
		return err
	}

	err = uploader.Process(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("process upload task failed")
		return err
	}

	indexBytes, err := yaml.Marshal(index)
	if err != nil {
		logrus.WithError(err).Errorln("marshal index failed")
		return err
	}

	err = os.WriteFile(indexPath, indexBytes, 0700)
	if err != nil {
		logrus.WithError(err).Errorln("write index file failed")
		return err
	}

	indexNewAddr, err := s.Upload(ctx, ".index", indexPath)
	if err != nil {
		logrus.WithError(err).Errorln("upload index file to fs failed")
		return err
	}
	newIdent := utils.Ident{
		Type:    "storage",
		SubType: storagePath,
		Merkle:  index.Root[:8],
		Address: indexNewAddr,
	}

	identStr, err := newIdent.String()
	if err != nil {
		logrus.WithError(err).Errorln("ident obj to string failed")
		return err
	}
	logrus.Infof("current timestamp t1 = %d (ms)", time.Now().UnixMilli())
	logrus.Infof("DWApp deployed on %s with MID: %s", storagePath, identStr)

	err = (*chain).SetIdentity(identStr)
	if err != nil {
		logrus.WithError(err).Errorln("update on-chain identity failed")
		return err
	}

	return nil
}
