/**
  @author: decision
  @date: 2024/6/17
  @note:
**/

package main

import (
	"context"
	"fmt"
	"github.com/everFinance/goar"
	"github.com/gookit/config/v2"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type ArweaveStorage struct {
	wallet *goar.Wallet
	client *goar.Client
}

func (s *ArweaveStorage) Initial(ctx context.Context) error {
	walletPath := config.String("deploy.storage.path")
	url := config.String("deploy.storage.url")

	wallet, err := goar.NewWalletFromPath(walletPath, url)
	if err != nil {
		logrus.WithError(err).Errorln("load wallet failed")
		return err
	}

	s.wallet = wallet
	s.client = goar.NewClient(url)

	return nil
}

func (s *ArweaveStorage) Ping(ctx context.Context) error {
	return nil
}
func (s *ArweaveStorage) Exists(ctx context.Context, identity string) (bool,
	error) {
	return s.client.ExistTxData(identity)
}

func (s *ArweaveStorage) Upload(ctx context.Context, name string, source string) (string, error) {
	if !filepath.IsAbs(source) {
		return "", fmt.Errorf("source path is not absolute path")
	}

	fileBytes, err := os.ReadFile(source)
	if err != nil {
		logrus.WithError(err).Debugf("read data from %s failed", source)
		return "", err
	}

	tx, err := s.wallet.SendData(fileBytes, nil)
	if err != nil {
		logrus.WithError(err).Errorln("save data on arweave failed")
		return "", err
	}

	logrus.Infof("save data to arweave with tx: %s", tx.ID)

	return tx.ID, nil
}

func (s *ArweaveStorage) Download(ctx context.Context, identity string,
	dst string) error {
	if !filepath.IsAbs(dst) {
		return fmt.Errorf("destination path is not absolute path")
	}

	fd, err := s.client.DownloadChunkData(identity)
	if err != nil {
		logrus.WithError(err).Errorln("open download data stream failed")
		return err
	}

	err = os.WriteFile(dst, fd, 0700)
	if err != nil {
		logrus.WithError(err).Debugf("write data to file %s failed", dst)
		return err
	}

	return nil
}

func (s *ArweaveStorage) Delete(ctx context.Context, identity string) error {
	return nil
}

func (s *ArweaveStorage) check() bool {
	var cfgs = []string{}

	for _, item := range cfgs {
		if !config.Exists(item) {
			logrus.Errorf("config item %s not exist", item)
			return false
		}
	}
	return true
}
