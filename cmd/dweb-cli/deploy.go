/**
  @author: decision
  @date: 2024/7/1
  @note: DWApp 部署的主要逻辑代码
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

	// 读取到上传的缓存，判断是否断点续传或重新上传
	if len(entities) != 0 {
		fmt.Printf("found cached file in %s, continue to upload? (Y/N): ")
		_, err = fmt.Scan(&opt)
		if err != nil {
			logrus.WithError(err).Debugf("read user option failed")
			return err
		}
	}

	// 如果不进行断点续传，则清除缓存文件
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

// checkChainIdentity 判断链上状态与本地 DWApp
func checkChainIdentity(ctx context.Context) (onChainStatus, error) {
	var err error

	// 根据 MID 加载 IChain 类型的插件
	chain, err = utils.ParseOnChain(chainIdent)
	if err != nil {
		logrus.WithError(err).Errorln("load chain plugin failed")
		return onChainDefault, err
	}

	// 获取链上存放的 DID 信息
	chainStorageIdent, err := (*chain).Identity()
	if err != nil {
		logrus.WithError(err).Errorln("get storage ident failed")
		return onChainDefault, err
	}

	// 根据本地的 DWApp 创建索引
	index, err = utils.CreateDirIndex(filePath)
	if err != nil {
		logrus.WithError(err).Errorln("create full directory index failed")
		return onChainDefault, err
	}
	// 计算 Merkle Tree 哈希值
	index.MerkleRoot()

	storageIdent = &utils.Ident{}
	err = storageIdent.FromString(chainStorageIdent)
	if err != nil {
		logrus.Infof("on-chain identity invalid, first time upload")
		return onChainUpload, nil
	}

	// 与链上的 DID 中的 merkle tree 哈希值进行对比
	if index.Root[:8] == storageIdent.Merkle {
		logrus.Infof("merkle root equal, deploy process exit")
		// 返回链上信息一致，无需更新
		return onChainNone, fmt.Errorf("merkle equal")
	}

	var originIndex string
	originIndexPath := filepath.Join(appDir, ".origin")
	originIndex, storage, err = utils.ParseFileStorage(ctx, chainStorageIdent)
	if err != nil {
		logrus.WithError(err).Errorln("parse storage identity failed")
		return onChainDefault, err
	}

	// 根据链上 DID 下载元文件
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

	// 本次的 DWApp 需要更新，返回 onChainUpdate
	return onChainUpdate, nil
}

// checkStorageDiff 文件差异判断
func checkStorageDiff(ctx context.Context) error {
	originPath := filepath.Join(appDir, "origin")
	diffBar := progressbar.Default(int64(len(index.Paths)))

	// 遍历链上 DID 对应的 DWApp 文件
	for k, v := range index.Paths {
		p, ok := origin.Paths[k]
		if ok {
			// 将文件下载到本地
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

			// 通过 CID 判断后，如果不一致则需要更新，设置为空字符串
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
	indexPath := filepath.Join(appDir, ".index")
	progressPath := filepath.Join(appDir, ".archive")
	storagePath := config.String("plugins.storage", "")

	if storagePath == "" {
		logrus.Error("config item 'plugins.storage' is empty")
		return fmt.Errorf("config item not exists")
	}

	// 加载文件上传插件
	symbol, err := utils.LoadSymbol(storagePath)
	if err != nil {
		logrus.WithError(err).Errorln("load plugin failed")
	}

	s, ok := symbol.(interfaces.IFileStorage)
	if !ok {
		logrus.Error("convert symbol to storage interface failed")
		return fmt.Errorf("convert symbol to interface failed")
	}

	// 插件初始化
	err = s.Initial(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("storage plugin initial failed")
		return err
	}

	// 创建上传器实例，并上传 DWApp
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

	// 将元文件信息写入到本地
	err = os.WriteFile(indexPath, indexBytes, 0700)
	if err != nil {
		logrus.WithError(err).Errorln("write index file failed")
		return err
	}

	// 上传元文件，并取得 CID
	indexNewAddr, err := s.Upload(ctx, ".index", indexPath)
	if err != nil {
		logrus.WithError(err).Errorln("upload index file to fs failed")
		return err
	}
	// 根据 CID 创建 DID
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
	logrus.Infof("DWApp deployed on %s with MID: %s", storagePath, identStr)

	// 更新链上的 DID 信息
	err = (*chain).SetIdentity(identStr)
	if err != nil {
		logrus.WithError(err).Errorln("update on-chain identity failed")
		return err
	}

	return nil
}
