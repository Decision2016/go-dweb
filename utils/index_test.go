/**
  @author: decision
  @date: 2024/6/27
  @note:
**/

package utils

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLoadIndex(t *testing.T) {
	path := "/Users/decision/Repos/go-dweb/cmd/dweb/.service/index/dd5bf960"
	_, err := LoadIndex(path)
	if err != nil {
		logrus.Fatal(err)
	}

}
