/**
  @author: decision
  @date: 2024/8/21
  @note:
**/

package web

import (
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/managers"
	"github.io/decision2016/go-dweb/utils"
	"sync"
	"time"
)

const cacheTimeout = 5 * time.Minute

var (
	observerInst *Observer
	observerOnce sync.Once
)

type Observer struct {
	loader *Loader

	checkChan  chan *utils.Ident
	checkCache map[string]time.Time

	lock sync.RWMutex
}

func ObserverDefault() *Observer {
	observerOnce.Do(func() {
		observerInst = &Observer{
			checkChan:  make(chan *utils.Ident),
			checkCache: make(map[string]time.Time),

			lock: sync.RWMutex{},
		}
	})

	return observerInst
}

func (o *Observer) Run() {
	go o.process()
}

func (o *Observer) SetLoader(loader *Loader) {
	o.loader = loader
}

func (o *Observer) Append(ident *utils.Ident) {
	identStr, err := ident.String()
	if err != nil {
		return
	}

	o.lock.Lock()
	cacheTime, ok := o.checkCache[identStr]
	if !ok {
		o.checkCache[identStr] = time.Now()
		o.lock.Unlock()

		o.checkChan <- ident
		return
	}

	if time.Now().Sub(cacheTime) < cacheTimeout {
		return
	}

	o.lock.Lock()
	o.checkCache[identStr] = time.Now()
	o.lock.Unlock()

	o.checkChan <- ident
	return
}

func (o *Observer) process() {
	cache := managers.CacheDefault()

	for {
		select {
		// 检查更新并且在后台进行更新，需要和原始目录隔离开，在完整下载服务之后再替换文件
		case ident := <-o.checkChan:
			identStr, _ := ident.String()

			indexPath := cache.IndexPath(identStr)
			index, err := utils.LoadIndex(indexPath)
			if err != nil {
				logrus.WithError(err).Debugf(
					"load index from cache directory failed")
				continue
			}

			chain, err := utils.ParseOnChain(ident)
			if err != nil {
				logrus.WithError(err).Debugf(
					"parse on-chain plugin from identity failed")
				continue
			}

			fsIdentStr, err := (*chain).Identity()
			if err != nil {
				logrus.WithError(err).Debugf("obtain on-chain identtiy failed")
				continue
			}

			fsIdent := utils.Ident{}
			err = fsIdent.FromString(fsIdentStr)
			if err != nil {
				logrus.WithError(err).Debugf("parse string to identity failed")
				continue
			}

			if fsIdent.Merkle == index.Root[:8] {
				continue
			}

			task := &LoadTask{
				Ident: ident,
				Type:  TypeUpdate,
			}
			o.loader.AppendTask(task)
		}
	}
}
