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
)

type DefaultService struct {
	port   int
	router *gin.Engine
}

func (s *DefaultService) process() {
	p := fmt.Sprintf(":%d", s.port)
	logrus.Infof("process web services on port %d", s.port)

	err := s.router.Run(p)
	if err != nil {
		logrus.WithError(err).Errorln("error occurred when running web service")
	}
}

func (s *DefaultService) Run() {
	s.router = gin.Default()
	s.port = 8080
}

func (s *DefaultService) Clean() error {
	return nil
}

func (s *DefaultService) page(c *gin.Context) {
	path := c.Request.URL.Path

	println(path)
}
