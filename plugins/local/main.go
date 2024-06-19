package main

import (
	"context"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	cfg "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"sync"
)

// 以插件的形式出现，对象的实例化基于配置文件
type LocalIPFS struct {
	repoPath string

	api  *iface.CoreAPI
	node *core.IpfsNode
}

var Instance LocalIPFS

var peers = []string{
	// 8api.sh
	"/ip4/78.46.108.24/tcp/4001/p2p/12D3KooWGASC2jm3pmohEJXUhuStkxDitPgzvs4qMuFPaiD9x1BA",
	"/ip4/65.109.19.136/tcp/4001/p2p/12D3KooWRbWZN3GvLf9CHmozq4vnTzDD4EEoiqtRJxg5FV6Gfjmm",
	"/ip4/120.226.39.189/tcp/4001/p2p/12D3KooWHRPErPuUPJZ6RRcEzsmgex93fx9V6h6vVCfcSpgyYGGS",
	"/ip4/120.226.39.187/tcp/4001/p2p/12D3KooWJYRM6GpGnusWxde1s2qL1HMJ7rWM7f5x9UTouhuVVwnA",

	// IPFS Bootstrapper nodes.
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

	// IPFS Cluster Pinning nodes
	"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
}

// check 检查配置项是否存在
func (i *LocalIPFS) check() bool {
	var cfgs = []string{}

	for _, item := range cfgs {
		if !config.Exists(item) {
			logrus.Errorf("config item %s not exist")
			return false
		}
	}

	return true
}

func (i *LocalIPFS) Initial(ctx context.Context) error {
	if !i.check() {
		return fmt.Errorf("required configuration not exist")
	}

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
	//api := *i.api

	return nil
}

func (i *LocalIPFS) start(ctx context.Context) {
	logrus.Infof("using %s as location ipfs repo.", i.repoPath)

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
	go i.connectPeers(ctx)

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

func (i *LocalIPFS) connectPeers(ctx context.Context) {
	var wg sync.WaitGroup
	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))

	for _, addr := range peers {
		multiAddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			logrus.WithError(err).Debugf("peer address error, skip...")
			continue
		}

		pii, err := peer.AddrInfoFromP2pAddr(multiAddr)
		if err != nil {
			logrus.WithError(err).Debugf("multiaddr to addr info error, skip...")
			continue
		}

		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()

			err := i.node.PeerHost.Connect(ctx, *peerInfo)
			if err != nil {
				logrus.WithError(err).Debugf("connect to %s failed", peerInfo.ID)
			}

			logrus.Debugf("connect to node %s success", peerInfo.ID)
		}(peerInfo)
	}

	wg.Wait()
}
