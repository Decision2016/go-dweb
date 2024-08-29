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
	port    int
	router  *gin.Engine
	loader  *Loader
	checker *Checker

	loaded map[string]bool
	tried  *lru.Cache

	//lock sync.Mutex			// 同步锁， TODO：后续考虑如果高并发的场景下需限制读写
}

func NewDWebService(ctx context.Context) (*DefaultService, error) {
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

	checker := CheckerDefault()
	checker.SetLoader(loader)
	checker.Run()
	service.checker = checker

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
	// TODO: 路径的处理
	identPath, err := utils.URLPathToChainIdent(path)
	if err != nil {
		logrus.WithError(err).Debugf("error occurred when convert url to path")
		c.JSON(http.StatusBadRequest, SimpleMsg(errRequestPathInvalid))
		return
	}

	cache := managers.CacheDefault()
	// 判断对应的 dapp web 是否已经加载到本地
	uid := cache.Uid(identPath)

	if !s.loaded[uid] {
		// 先校验本地的文件是否存在以及是否完整, 如果存在且完整则设置 cache 为已加载
		// todo: 这里的检验方法需要修复一下
		valid, err := cache.Validate(identPath)
		if valid {
			s.loaded[uid] = true
		} else {
			logrus.
				WithError(err).
				Debugf("application %s invalid on disk", uid)
			s.loader.AppendTaskByString(identPath)
			s.tried.Add(uid, nil)
			c.JSON(http.StatusInternalServerError, "application not valid, "+
				"waiting for reload...")
			return
		}
	}

	ident := &utils.Ident{}
	err = ident.FromString(identPath)
	if err != nil {
		logrus.WithError(err).Debugf("parse %s to identity failed", identPath)
		c.JSON(http.StatusBadRequest, SimpleMsg(errRequestPathInvalid))
		return
	}
	s.checker.Append(ident)

	c.Set("ident", identPath)
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

	// 根据 ident 获取本地的路径信息
	cache := managers.CacheDefault()
	location := cache.Path(ident.(string))
	absPath := filepath.Join(location, filePath)
	c.File(absPath)
}

func (s *DefaultService) loadCallback(uid string) {
	s.loaded[uid] = true
}
