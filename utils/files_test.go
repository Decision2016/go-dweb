/**
  @author: decision
  @date: 2024/6/14
  @note:
**/

package utils

import (
	"github.com/tj/assert"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestGetFileCidV0(t *testing.T) {
	filepath := "./test/example.jpg"

	cid, err := GetFileCidV0(filepath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "QmbBp5huHazG582ean3eGGwWXkkVKnRc7V5te16ByWNs2N",
		cid.String(), "cid output not equal")
}

func TestCreateDirIndex(t *testing.T) {
	filepath := "/Users/decision/Repos/go-dweb/cmd/dweb-cli/test"

	full, err := CreateDirIndex(filepath)
	if err != nil {
		t.Fatal(err)
	}

	b, err := yaml.Marshal(full)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("full.yml", b, 0700)
	if err != nil {
		t.Fatal(err)
	}
}
