package deploy

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func AppInitial() {
	var path string

	fmt.Println("This utility will walk you through creating a dweb application work directory.\n")
	fmt.Println("Press ^C at any time to quit.\n")
	fmt.Print("Work directory path:")
	_, err := fmt.Scan(&path)
	if err != nil {
		logrus.Error("read path string error")
		return
	}

	err = appInitial(path)
	if err != nil {
		logrus.Error(err)
	}
}
