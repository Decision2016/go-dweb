package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/storage"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.LoadGlobalConfig("conf.yml")

	ipfs := storage.LocalIPFS{}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := ipfs.Initial(ctx)
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
