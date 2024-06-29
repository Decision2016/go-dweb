/**
  @author: decision
  @date: 2024/6/7
  @note:
**/

package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/managers"
	"github.io/decision2016/go-dweb/utils"
	"net/http"
	"path/filepath"
)

type DefaultService struct {
	port   int
	router *gin.Engine
	loader *Loader

	loaded map[string]bool
	tried  *lru.Cache

	//lock sync.Mutex			// 同步锁， TODO：后续考虑如果高并发的场景下需限制读写
}

func NewDWebService(ctx context.Context) (*DefaultService, error) {
	// todo: cache 应该根据配置文件的方式来配置，并且还需要进一步检查工作目录是否存在
	managers.cache.Initial()

	c, err := lru.New(300)
	if err != nil {
		logrus.WithError(err).Debugf("create lru cache failed")
		return nil, err
	}

	port := config.Int("web.port", 8080)

	service := DefaultService{
		port:   port,
		router: gin.Default(),

		loaded: make(map[string]bool),
		tried:  c,
	}

	loader := NewLoader(ctx, service.loadCallback)
	loader.Run(ctx)
	service.loader = loader

	return &service, nil
}

func (s *DefaultService) Run() {
	go s.process()
}

func (s *DefaultService) Clean() error {
	return nil
}

func (s *DefaultService) process() {
	p := fmt.Sprintf(":%d", s.port)
	logrus.Infof("process web services on port %d", s.port)

	s.router.Use(s.middleware)
	s.router.GET("/*path", s.handle)
	err := s.router.Run(p)
	if err != nil {
		logrus.WithError(err).Fatalln("error occurred when running web service")
	}
}

func (s *DefaultService) middleware(c *gin.Context) {
	path := c.Request.URL.Path
	ident, err := utils.URLPathToChainIdent(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, "request path invalid")
		logrus.WithError(err).Debugf("error occurred when convert url to path")
		return
	}

	uid := managers.cache.uid(ident)
	if !s.loaded[uid] {
		_, ok := s.tried.Get(uid)
		if !ok {
			logrus.Debugf("append task %s to loader", uid)
			s.loader.AppendTask(ident)
			s.tried.Add(uid, nil)
		}

		c.JSON(http.StatusNotFound, "app not exist or waiting for load")
		return
	}

	c.Set("ident", ident)
	c.Set("uid", uid)
	c.Next()
}

func (s *DefaultService) handle(c *gin.Context) {
	path := c.Param("path")

	filePath, err := utils.ExtractFilePath(path)
	if err != nil {
		logrus.WithError(err).Debugf("file path not exist")
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	ident, ok := c.Get("ident")
	if !ok {
		c.JSON(http.StatusInternalServerError, nil)
	}

	location := managers.cache.Path(ident.(string))
	absPath := filepath.Join(location, filePath)
	c.File(absPath)
}

func (s *DefaultService) loadCallback(uid string) {
	s.loaded[uid] = true
}
