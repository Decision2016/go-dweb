/**
  @author: decision
  @date: 2024/6/24
  @note:
**/

package utils

import (
	"github.com/tj/assert"
	"testing"
)

func TestMerkleRoot(t *testing.T) {
	data1 := []string{"test1", "test", "test2", "111", " ", "2"}
	merkle := MerkleRoot(data1)

	assert.Equal(t,
		"6d95630b3e1d85e0abe2bda91167b0e07570767912f7e4af0559e9c77d696716",
		merkle, "merkle root wrong")
}
