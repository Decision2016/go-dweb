/**
  @author: decision
  @date: 2024/6/13
  @note:
**/

package plugins

import (
	"context"
	"github.com/gookit/config/v2"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	http "net/http"
)

type InfuraIPFS struct {
	baseUrl string
	id      string
	secret  string

	node *rpc.HttpApi
}

func (i *InfuraIPFS) Initial(ctx context.Context) error {
	i.baseUrl = config.String("deploy.storage.baseurl")
	i.id = config.String("deploy.storage.id")
	i.secret = config.String("deploy.storage.secret")

	client := i.newClient(i.id, i.secret)
	node, err := rpc.NewURLApiWithClient(i.baseUrl, client)
	if err != nil {
		logrus.WithError(err).Fatalln("http api load failed")
		return err
	}

	i.node = node

	return nil
}

func (i *InfuraIPFS) Ping(ctx context.Context) error {
	return nil
}
func (i *InfuraIPFS) Exists(ctx context.Context, identity string) (bool, error) {
	return true, nil
}

func (i *InfuraIPFS) Upload(ctx context.Context, name string, source string) error {
	fileNode, err := utils.GetUnixFsNode(source)
	if err != nil {
		logrus.WithError(err).Errorln("get file node failed")
		return err
	}

	cid, err := i.node.Unixfs().Add(ctx, fileNode)
	if err != nil {
		logrus.WithError(err).Errorln("add file failed")
		return err
	}

	logrus.Infof("add file %s with cid %s", name, cid)
	return nil

	return nil
}

func (i *InfuraIPFS) Download(ctx context.Context, identity string, dst string) error {
	p, err := path.NewPath(identity)
	if err != nil {
		logrus.WithError(err).Errorf("convert string %s to path failed", identity)
		return err
	}

	fileNode, err := i.node.Unixfs().Get(ctx, p)
	if err != nil {
		logrus.WithError(err).Errorf("get file with identity %s failed", identity)
		return err
	}

	err = files.WriteTo(fileNode, dst)
	if err != nil {
		logrus.WithError(err).Errorf("write binary to %s failed", dst)
		return err
	}
	logrus.Debugf("download file %s to %s success", identity, dst)

	return nil
}

func (i *InfuraIPFS) Delete(ctx context.Context, identity string) error {
	return nil
}

func (i *InfuraIPFS) newClient(id, secret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     id,
			ProjectSecret: secret,
		},
	}
}

type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}
