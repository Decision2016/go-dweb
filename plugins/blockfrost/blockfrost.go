/**
  @author: decision
  @date: 2024/6/21
  @note:
**/

package blockfrost

import (
	"context"
	"fmt"
	"github.com/blockfrost/blockfrost-go"
	"github.com/gookit/config/v2"
	"github.com/ipfs/boxo/path"
	"github.com/sirupsen/logrus"
	"os"
)

type BlockFrostIPFS struct {
	id string

	api *blockfrost.IPFSClient
}

var Instance BlockFrostIPFS

func (i *BlockFrostIPFS) check() bool {
	var cfgs = []string{
		"deploy.storage.id",
	}

	for _, item := range cfgs {
		if !config.Exists(item) {
			logrus.Errorf("config item %s not exist", item)
			return false
		}
	}

	return true
}

func (i *BlockFrostIPFS) Initial(ctx context.Context) error {
	if !i.check() {
		return fmt.Errorf("required configuration not exist")
	}

	i.id = config.String("deploy.storage.id")

	client := blockfrost.NewIPFSClient(blockfrost.IPFSClientOptions{
		ProjectID: i.id,
	})

	i.api = &client
	return nil
}

func (i *BlockFrostIPFS) Ping(ctx context.Context) error {
	return nil
}

func (i *BlockFrostIPFS) Exists(ctx context.Context, identity string) (bool, error) {
	return true, nil
}

func (i *BlockFrostIPFS) Upload(ctx context.Context, name string, source string) error {
	api := *i.api

	obj, err := api.Add(ctx, source)
	if err != nil {
		logrus.WithError(err).Debugf("add file %s to ipfs failed", source)
		return err
	}

	pin, err := api.Pin(ctx, obj.IPFSHash)
	if err != nil {
		logrus.WithError(err).Debugf("error occured when pinning file %s", obj.IPFSHash)
		return err
	}

	logrus.Debugf("pin file %s to blockfrost success", pin.IPFSHash)
	return nil
}

func (i *BlockFrostIPFS) Download(ctx context.Context, identity string, dst string) error {
	api := *i.api

	p, err := path.NewPath(identity)
	if err != nil {
		logrus.WithError(err).Debugf("convert string %s to path failed",
			identity)
		return err
	}

	data, err := api.Gateway(ctx, p.String())
	if err != nil {
		logrus.WithError(err).Debugf("get ipfs pinned object failed")
		return err
	}

	err = os.WriteFile(dst, data, 0700)
	if err != nil {
		logrus.WithError(err).Debugf("write byte data to file errored")
		return err
	}

	return nil
}

func (i *BlockFrostIPFS) Delete(ctx context.Context, identity string) error {
	api := *i.api

	_, err := api.Remove(ctx, identity)
	if err != nil {
		logrus.WithError(err).Debugf("remove file %s failed", identity)
		return err
	}

	return nil
}
