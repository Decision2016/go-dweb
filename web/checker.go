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
	checkerInst *Checker
	checkerOnce sync.Once
)

type Checker struct {
	loader *Loader

	checkChan  chan *utils.Ident
	checkCache map[string]time.Time

	lock sync.RWMutex
}

func CheckerDefault() *Checker {
	checkerOnce.Do(func() {
		checkerInst = &Checker{
			checkChan:  make(chan *utils.Ident, 100),
			checkCache: make(map[string]time.Time),

			lock: sync.RWMutex{},
		}
	})

	return checkerInst
}

func (c *Checker) Run() {
	go c.process()
}

func (c *Checker) SetLoader(loader *Loader) {
	c.loader = loader
}

func (c *Checker) Append(ident *utils.Ident) {
	identStr, err := ident.String()
	if err != nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	cacheTime, ok := c.checkCache[identStr]
	if !ok {
		c.checkCache[identStr] = time.Now()

		c.checkChan <- ident
		return
	}

	if time.Now().Sub(cacheTime) < cacheTimeout {
		return
	}
	c.checkCache[identStr] = time.Now()

	c.checkChan <- ident
	return
}

func (c *Checker) process() {
	cache := managers.CacheDefault()

	for {
		select {
		// 检查更新并且在后台进行更新，需要和原始目录隔离开，在完整下载服务之后再替换文件
		case ident := <-c.checkChan:
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
			c.loader.AppendTask(task)
		}
	}
}
