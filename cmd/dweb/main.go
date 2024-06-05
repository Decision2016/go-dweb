/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(static.Serve("/",
		static.LocalFile("/Users/decision/Repos/go-dweb/static", false),
	))
	router.Run(":8080")
}
