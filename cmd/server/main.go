package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"time"
)

func main() {
	_ = utils.LoadGlobalConfig("config.yml")
	r := gin.Default()

	r.POST("/new", func(context *gin.Context) {
		var req utils.Ident
		err := context.ShouldBindBodyWithJSON(&req)

		if err != nil {
			return
		}

		go worker(&req)
	})

	r.Run(":47777")
}

func worker(ident *utils.Ident) {
	strIdent, _ := ident.String()
	indexIdent, fs, err := utils.ParseFileStorage(context.TODO(), strIdent)
	if err != nil {
		return
	}

	err = (*fs).Initial(context.TODO())
	if err != nil {
		logrus.WithError(err).Debugf("file storage initial failed")
		return
	}

	for {
		err := (*fs).Download(context.TODO(), indexIdent, "./index")
		if err != nil {
			logrus.Infof("file %s download t_3 = %d", ident.Address, time.Now().UnixMilli())
		}
	}
}
