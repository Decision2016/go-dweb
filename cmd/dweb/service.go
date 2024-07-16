/**
  @author: decision
  @date: 2024/7/16
  @note:
**/

package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/utils"
	"github.io/decision2016/go-dweb/web"
	"net/http"
	"path/filepath"
)

var (
	localFilePath string
	localIdentity string
)

var serviceCmd = &cobra.Command{
	Use:              "service",
	Short:            "Decentralized web application service",
	Long:             "Decentralized web application service",
	TraverseChildren: true,
}

var serviceRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run DWeb service",
	Long:  "Run DWeb service",
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.LoadGlobalConfig("./config.yml")
		if err != nil {
			logrus.WithError(err).Errorln("load config file failed")
			return
		}

		ctx := context.Background()
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
	},
}

var serviceLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Run DWeb local service",
	Long:  "Run DWeb local service",
	Run: func(cmd *cobra.Command, args []string) {
		port := config.Int("web.port", 8080)
		p := fmt.Sprintf(":%d", port)

		r := gin.Default()
		r.GET("/*path", func(c *gin.Context) {
			path := c.Param("path")

			ident, err := utils.URLPathToChainIdent(path)
			if err != nil {
				c.JSON(http.StatusBadRequest, "request path invalid")
				logrus.WithError(err).Debugf("error occurred when convert url to path")
				return
			}

			if ident != localIdentity {
				c.JSON(http.StatusBadRequest, "identity not equal to option")
				return
			}

			filePath, err := utils.ExtractFilePath(path)
			if err != nil {
				logrus.WithError(err).Debugf("file path not exist")
				c.JSON(http.StatusBadRequest, nil)
				return
			}

			absPath := filepath.Join(localFilePath, filePath)
			c.File(absPath)
		})

		err := r.Run(p)
		if err != nil {
			logrus.WithError(err).Fatalln("error occurred when running web service")
		}

		utils.Waiting(nil)
	},
}

func init() {
	serviceLocalCmd.Flags().StringVarP(&localFilePath, "path", "p", "",
		"static web service path (required)")
	serviceLocalCmd.Flags().StringVarP(&localIdentity, "identity", "i", "",
		"local service identity (required)")
	serviceLocalCmd.MarkFlagRequired("path")
	serviceLocalCmd.MarkFlagRequired("identity")
}
