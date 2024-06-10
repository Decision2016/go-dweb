package storage

import (
	"context"
	"github.com/gookit/config/v2"
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

func (i *LocalIPFS) Ping() bool {
	return true
}

func (i *LocalIPFS) start(ctx context.Context) {
	logrus.Infof("Using %s as location ipfs repo.", i.repoPath)

	exists, err := utils.CheckDirExistAndEmpty(i.repoPath)
	if err != nil {
		logrus.Fatalln(err)
	}

	if !exists {
		logrus.Info("Repo directory is not exist, creating new ...")

		plugins, err := loader.NewPluginLoader(i.repoPath)
		if err != nil {
			logrus.WithError(err).Fatalln("load repo path failed")
		}

		err = plugins.Initialize()
		if err != nil {
			logrus.WithError(err).Fatalln("plugin initial failed")
		}

		err = plugins.Inject()
		if err != nil {
			logrus.WithError(err).Fatalln("plugin inject failed")
		}

		c, err := cfg.Init(os.Stdout, 2048)
		if err != nil {
			logrus.WithError(err).Fatalln("Local IPFS repo initial failed")
			return
		}

		err = fsrepo.Init(i.repoPath, c)
		if err != nil {
			logrus.WithError(err).Fatalln("Local IPFS repo initial failed")
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("Check directory failed.")
		return
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
		logrus.WithError(err).Error("Run ipfs node failed.")
		return
	}
	addrs := node.PeerHost.Addrs()
	logrus.Infof("Local node address: %s", addrs[0])

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		logrus.WithError(err)
	}
	i.api = &api
	i.node = node

	<-ctx.Done()

	logrus.Infoln("Receive context signal, shutting down IPFS node...")
}
