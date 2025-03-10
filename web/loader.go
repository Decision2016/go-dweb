package web

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/interfaces"
	"github.io/decision2016/go-dweb/managers"
	"github.io/decision2016/go-dweb/metrics"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TaskType int8

const (
	TypeLoad   TaskType = 0x01
	TypeUpdate TaskType = 0x02
)

type LoadTask struct {
	Ident *utils.Ident
	Type  TaskType
}

// Loader 接收任务从链上拉取配置信息，然后获取状态
type Loader struct {
	ctx context.Context

	queue chan *LoadTask
	mp    sync.Map

	updateFunc func(uid string)
}

const (
	loaderDownloadTimeout = 30 * time.Second
)

func NewLoader(ctx context.Context, callback func(uid string)) *Loader {
	loader := &Loader{
		ctx:        ctx,
		queue:      make(chan *LoadTask, 30),
		mp:         sync.Map{},
		updateFunc: callback,
	}

	return loader
}

func (l *Loader) Run(ctx context.Context) {
	l.ctx = ctx
	go l.processTask()
}

func (l *Loader) AppendTask(task *LoadTask) {
	metrics.LoaderTaskCountInQueue.Inc()
	l.queue <- task
}

// AppendTaskByString 添加加载任务
func (l *Loader) AppendTaskByString(identStr string) {
	// 检查队列中是否存在对应的 ident
	_, ok := l.mp.Load(identStr)
	if ok {
		return
	}

	ident := utils.Ident{}
	err := ident.FromString(identStr)
	if err != nil {
		return
	}

	task := &LoadTask{
		Ident: &ident,
		Type:  TypeLoad,
	}
	metrics.LoaderTaskCountInQueue.Inc()
	l.queue <- task
}

// downloadApp 根据 index 中的映射关系表下载 DWApp
func (l *Loader) downloadApp(chainIdent string, index *utils.Index,
	fs *interfaces.IFileStorage, parentDir string) error {
	total := len(index.Paths)
	count := 0

	cache := managers.CacheDefault()
	errored := false

	uid := cache.Uid(chainIdent)

	metrics.LoaderCurrentTaskProgress.Set(0)

	// 遍历并调用 IStorage 类型的插件下载
	for p, ident := range index.Paths {
		dst := filepath.Join(parentDir, p)
		//ctx, cancel := context.WithTimeout(l.ctx, loaderDownloadTimeout)
		//ctx, cancel := context.WithCancel(l.ctx)
		err := (*fs).Download(l.ctx, ident, dst)

		//cancel()
		if err != nil {
			logrus.
				WithError(err).
				WithField("dst", dst).
				Debugf("download file from storage failed")
			errored = true
			break
		}

		count += 1
		progress := 100.0 * float64(count) / float64(total)
		metrics.LoaderCurrentTaskProgress.Set(progress)
		logrus.Debugf("current download progress: %f", progress)
	}

	if errored {
		err := cache.Delete(chainIdent)
		if err != nil {
			logrus.WithError(err).Debugf("error occurred when removing file cache")
			return err
		}
		logrus.Infoln("download cache removed")
	} else {
		logrus.Infof("app with uid %s downloaded", uid)
	}
	l.updateFunc(uid)

	return nil
}

// processTask DWApp 加载器的运行协程
func (l *Loader) processTask() {
	cache := managers.CacheDefault()
	cache.Initial()

	for {
		select {
		case task := <-l.queue:
			metrics.LoaderTaskCountInQueue.Desc()

			// 对链上唯一标识进行解析
			ident := task.Ident
			chain, err := utils.ParseOnChain(ident)
			if err != nil {
				logrus.WithError(err).Debugf(
					"parser onchain identity %s failed", ident)
				continue
			}

			strIdent, _ := ident.String()
			// 获取链上存放的 FS 索引信息
			fsIdent, err := (*chain).Identity()
			if err != nil {
				logrus.WithError(err).Debug("load on-chain identity failed")
				continue
			}

			// 根据去中心化索引解析得到所需要的 identity
			indexIdent, fs, err := utils.ParseFileStorage(l.ctx, fsIdent)
			if err != nil {
				logrus.
					WithError(err).
					WithField("ident", fsIdent).
					Debugf("load filestorage failed")
				continue
			}

			err = (*fs).Initial(l.ctx)
			if err != nil {
				logrus.WithError(err).Debugf("file storage initial failed")
				continue
			}

			var dst, appDir string
			switch task.Type {
			case TypeLoad:
				dst = cache.IndexPath(strIdent)
				appDir = cache.Path(strIdent)
			case TypeUpdate:
				updateDir := cache.UpdatePath()
				err = cache.Clean(updateDir)
				if err != nil {
					logrus.WithError(err).Debugf("clean update directory failed")
					continue
				}

				dst = filepath.Join(updateDir, "index")
				appDir = filepath.Join(updateDir, "app")
			}

			if err = cache.RemoveIfExists(dst); err != nil {
				logrus.WithError(err).Debugf("error occurred when deleting index file")
				continue
			}

			logrus.Debugf("download index %s to destintation: %s",
				indexIdent, dst)
			err = (*fs).Download(l.ctx, indexIdent, dst)
			if err != nil {
				logrus.WithError(err).Debugf("donwload index file failed")
				continue
			}

			index, err := utils.LoadIndex(dst)
			if err != nil {
				logrus.WithError(err).Debugf("load index from file failed")
				continue
			}
			logrus.Debugln("load index file from storage success")

			err = l.downloadApp(strIdent, index, fs, appDir)
			if err != nil {
				logrus.WithError(err).Errorf("download dapp %d failed", ident)
			}

			if task.Type == TypeUpdate {
				realIndexPath := cache.IndexPath(strIdent)
				realAppDir := cache.Path(strIdent)

				if err = os.Rename(dst, realIndexPath); err != nil {
					logrus.WithError(err).Debugf(
						"move temporary index to workdir failed")
					continue
				}

				if err = os.Rename(appDir, realAppDir); err != nil {
					logrus.WithError(err).Debugf(
						"move temporary app to workdir failed")
					continue
				}
			}
			l.mp.Delete(strIdent)
		}
	}
}
