/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"context"
	"github.com/gookit/config/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"github.io/decision2016/go-dweb/web"
	"net/http"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})

	ctx := context.Background()
	err := utils.LoadGlobalConfig("./conf.yml")
	if err != nil {
		logrus.WithError(err).Errorln("load config file failed")
		return
	}

	enableMetrics := config.Bool("web.metrics.enable", false)
	if enableMetrics {
		port := config.String("web.metrics.port", "9090")
		metricPort := ":" + port
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe(metricPort, nil)
		logrus.Infof("metrics server start on localhost%s", metricPort)
	}

	service, err := web.NewDWebService(ctx)
	if err != nil {
		logrus.WithError(err).Errorln("create new dweb service failed")
		return
	}
	service.Run()

	utils.Waiting(nil)
}
