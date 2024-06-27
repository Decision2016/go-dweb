/**
  @author: decision
  @date: 2024/6/26
  @note:
**/

package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Waiting(callback func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			logrus.Info("Receive exit signal, exit...")
			if callback != nil {
				callback()
			}
			os.Exit(1)
		}
	}
}
