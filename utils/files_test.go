/**
  @author: decision
  @date: 2024/6/14
  @note:
**/

package utils

import (
	"github.com/tj/assert"
	"testing"
)

func TestGetFileCidV0(t *testing.T) {
	filepath := "./test/example.jpg"

	cid, err := GetFileCidV0(filepath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n",
		cid.String(), "cid output not equal")
}
