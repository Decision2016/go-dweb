package plugins

import (
	"context"
	"github.com/gookit/config/v2"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	cfg "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"os"
)

// 以插件的形式出现，对象的实例化基于配置文件
type LocalIPFS struct {
	repoPath string

	api  *iface.CoreAPI
	node *core.IpfsNode
}

var Instance LocalIPFS

func (i *LocalIPFS) Initial(ctx context.Context) error {
	i.repoPath = config.String("deploy.storage.location")

	go i.start(ctx)
	return nil
}

func (i *LocalIPFS) Ping(ctx context.Context) error {
	return nil
}

func (i *LocalIPFS) Exists(ctx context.Context, source string) (bool, error) {

	return true, nil
}

// Upload 将文件存放到 IPFS 下
// filepath 在这个项目下统一使用绝对路径，只有显式硬编码时才可以使用相对路径
func (i *LocalIPFS) Upload(ctx context.Context, name string, source string) error {
	api := *i.api

	fileNode, err := utils.GetUnixFsNode(source)
	if err != nil {
		logrus.WithError(err).Errorln("get file node failed")
		return err
	}

	cid, err := api.Unixfs().Add(ctx, fileNode)
	if err != nil {
		logrus.WithError(err).Errorln("add file failed")
		return err
	}

	logrus.Infof("add file %s with cid %s", name, cid)
	return nil
}

func (i *LocalIPFS) Download(ctx context.Context, identity string, dst string) error {
	api := *i.api

	p, err := path.NewPath(identity)
	if err != nil {
		logrus.WithError(err).Errorf("convert string %s to path failed", identity)
		return err
	}

	fileNode, err := api.Unixfs().Get(ctx, p)
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

func (i *LocalIPFS) Delete(ctx context.Context, identity string) error {
	return nil
}

func (i *LocalIPFS) start(ctx context.Context) {
	logrus.Infof("Using %s as location ipfs repo.", i.repoPath)

	exists, err := utils.CheckDirExistAndEmpty(i.repoPath)

	if !exists {
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Repo directory is not exist, creating new ...")

		c, err := cfg.Init(os.Stdout, 2048)
		if err != nil {
			logrus.WithError(err).Fatalln("local IPFS repo initial failed")
			return
		}

		err = i.load()
		if err != nil {
			logrus.WithError(err).Fatalln("load plugins failed")
		}

		err = fsrepo.Init(i.repoPath, c)
		if err != nil {
			logrus.WithError(err).Fatalln("local IPFS repo initial failed")
			return
		}
	} else {
		err = i.load()
		if err != nil {
			logrus.WithError(err).Fatalln("load plugins failed")
		}
	}

	repo, err := fsrepo.Open(i.repoPath)
	if err != nil {
		logrus.WithError(err).Error("Open repo failed.")
		return
	}

	node, err := core.NewNode(ctx, &core.BuildCfg{
		Repo:   repo,
		Online: true,
	})
	if err != nil {
		logrus.WithError(err).Error("run ipfs node failed.")
		return
	}
	addrs := node.PeerHost.Addrs()
	logrus.Infof("local node address: %s", addrs[0])

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		logrus.WithError(err)
	}
	i.api = &api
	i.node = node

	<-ctx.Done()

	logrus.Infoln("Receive context signal, shutting down IPFS node...")
}

func (i *LocalIPFS) load() error {
	plugins, err := loader.NewPluginLoader(i.repoPath)
	if err != nil {
		return err
	}

	err = plugins.Initialize()
	if err != nil {
		return err
	}

	err = plugins.Inject()
	if err != nil {
		return err
	}

	return nil
}
