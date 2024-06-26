/**
  @author: decision
  @date: 2024/6/5
  @note:
**/

package main

import (
	"context"
	"github.io/decision2016/go-dweb/utils"
	"github.io/decision2016/go-dweb/web"
	"os"
	"os/signal"
	"syscall"
)

//var example = "/chain/evm/0x7313AA3d928479Deb7debaC9E9d38286496D542e"

func main() {
	_, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	service := web.DefaultService{}
	service.Run()

	utils.Waiting(func() {
		cancel()
	})
}
