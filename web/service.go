/**
  @author: decision
  @date: 2024/6/7
  @note:
**/

package web

import "github.com/gin-gonic/gin"

type DefaultService struct {
	router *gin.Engine
}

func (s *DefaultService) Run() {

}

func (s *DefaultService) Clean() error {
	return nil
}
