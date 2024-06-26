package web

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/interfaces"
	"github.io/decision2016/go-dweb/metrics"
	"github.io/decision2016/go-dweb/utils"
	"path/filepath"
	"sync"
	"time"
)

// Loader 接收任务从链上拉取配置信息，然后获取状态
type Loader struct {
	ctx context.Context

	queue chan string
	mp    sync.Map
}

var (
	loaderDownloadTimeout = 30 * time.Second
	once                  sync.Once
	loader                *Loader
)

func LoaderInstance() *Loader {
	once.Do(func() {
		loader = &Loader{
			ctx:   context.TODO(),
			queue: make(chan string),
			mp:    sync.Map{},
		}
	})

	return loader
}

func (l *Loader) Run(ctx context.Context) {
	l.ctx = ctx
	go l.processTask()
}

func (l *Loader) AppendTask(ident string) {
	_, ok := l.mp.Load(ident)
	if ok {
		return
	}

	metrics.LoaderTaskCountInQueue.Inc()
	l.queue <- ident
}

func (l *Loader) downloadApp(chainIdent string, index *utils.FullStruct,
	fs *interfaces.IFileStorage) error {
	total := len(index.Paths)
	count := 0
	parentDir := cache.Path(chainIdent)
	errored := false

	metrics.LoaderCurrentTaskProgress.Set(0)

	for p, ident := range index.Paths {
		dst := filepath.Join(parentDir, p)
		ctx, cancel := context.WithTimeout(l.ctx, loaderDownloadTimeout)
		err := (*fs).Download(ctx, ident, dst)

		cancel()
		if err != nil {
			logrus.WithError(err).Debugf("download file from storage failed")
			errored = true
			break
		}

		count += 1
		metrics.LoaderCurrentTaskProgress.Set(100.0 * float64(
			count) / float64(total))
	}

	if errored {
		err := cache.Delete(chainIdent)
		if err != nil {
			logrus.WithError(err).Debugf("error occurred when removing file cache")
			return err
		}
		logrus.Infoln("download cache removed")
	}

	return nil
}

func (l *Loader) processTask() {
	for {
		select {
		case ident := <-l.queue:
			metrics.LoaderTaskCountInQueue.Desc()

			// 对链上唯一标识进行解析
			chain, err := utils.ParseOnChain(ident)
			if err != nil {
				logrus.WithError(err).Debugf(
					"parser onchain identity %s failed", ident)
				return
			}

			// 获取链上存放的 FS 索引信息
			fsIdent, err := (*chain).Identity()
			if err != nil {
				logrus.WithError(err).Debug("load on-chain identity failed")
				return
			}

			// 根据去中心化索引解析得到所需要的 identity
			indexIdent, fs, err := utils.ParseFileStorage(l.ctx, fsIdent)
			if err != nil {
				logrus.WithError(err).Debugf("load filestorage failed")
				return
			}

			// App 的索引信息拉取
			dst := cache.IndexPath(indexIdent)
			err = (*fs).Download(l.ctx, indexIdent, dst)
			if err != nil {
				logrus.WithError(err).Debugf("donwload index file failed")
				return
			}

			index, err := utils.LoadIndex(dst)
			if err != nil {
				logrus.WithError(err).Debugf("load index from file failed")
				return
			}

			err = l.downloadApp(ident, index, fs)
			if err != nil {
				logrus.WithError(err).Errorf("download dapp %d failed", ident)
			}
		}
	}
}
