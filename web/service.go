/**
  @author: decision
  @date: 2024/6/7
  @note:
**/

package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"path/filepath"
)

type DefaultService struct {
	port   int
	router *gin.Engine

	loaded map[string]bool
}

func (s *DefaultService) Run() {
	s.router = gin.Default()
	s.port = 8080
	s.loaded = make(map[string]bool)

	go s.process()
}

func (s *DefaultService) Clean() error {
	return nil
}

func (s *DefaultService) page(c *gin.Context) {
	path := c.Request.URL.Path

	println(path)
}

func (s *DefaultService) process() {
	p := fmt.Sprintf(":%d", s.port)
	logrus.Infof("process web services on port %d", s.port)

	s.router.Use(s.middleware)
	s.router.GET("/*path", s.handle)
	err := s.router.Run(p)
	if err != nil {
		logrus.WithError(err).Errorln("error occurred when running web service")
	}
}

func (s *DefaultService) middleware(c *gin.Context) {
	path := c.Request.URL.Path
	ident, err := utils.URLPathToChainIdent(path)
	if err != nil {
		c.JSON(400, "request path invalid")
		logrus.WithError(err).Debugf("error occurred when convert url to path")
		c.Next()
	}

	uid := cache.uid(ident)
	if !s.loaded[uid] {
		c.JSON(404, "app not exist or waiting for load")
		return
	}

	c.Keys["ident"] = ident
	c.Keys["uid"] = uid
	c.Next()
}

func (s *DefaultService) handle(c *gin.Context) {
	path := c.Param("path")

	filePath, err := utils.ExtractFilePath(path)
	if err != nil {
		logrus.WithError(err).Debugf("file path not exist")
		c.JSON(400, nil)
		return
	}

	ident, ok := c.Get("ident")
	if !ok {
		c.JSON(500, nil)
	}

	location := cache.Path(ident.(string))
	absPath := filepath.Join(location, filePath)
	c.File(absPath)
}
