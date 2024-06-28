/**
  @author: decision
  @date: 2024/6/28
  @note:
**/

package utils

import "testing"

func TestCreateWorkDir(t *testing.T) {
	err := CreateWorkDir()
	if err != nil {
		t.Fatal(err)
	}
}
