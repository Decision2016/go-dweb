package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/plugins"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logrus.SetFormatter(&utils.CustomFormatter{})
	utils.LoadGlobalConfig("conf.yml")

	ipfs := plugins.LocalIPFS{}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := ipfs.Initial(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	//time.Sleep(5 * time.Second)
	//err = ipfs.Upload(ctx, "test", "./test.txt")
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	time.Sleep(5 * time.Second)
	err = ipfs.Download(ctx,
		"/ipfs/QmRrscxG2pywPmSTBejLPJ16fnQSVPagNQpdK2xaTf92dY",
		"./download.png")
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		select {
		case <-sigs:
			logrus.Info("Receive exit signal, exit...")
			cancel()
			os.Exit(1)
		}
	}

}
