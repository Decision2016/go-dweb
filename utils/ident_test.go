/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package utils

import (
	"github.com/tj/assert"
	"testing"
)

func TestIdent_String(t *testing.T) {
	i1 := Ident{
		Type:    "chain",
		SubType: "evm",
		Merkle:  "",
		Address: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
	}

	i2 := Ident{
		Type:    "storage",
		SubType: "ipfs",
		Merkle:  "a2e24ff4",
		Address: "QmXViwQ1frFwabQHtmpt18SUPhnpcRzhWayt9rnTJ8GTay",
	}

	i1str, err := i1.String()
	if err != nil {
		t.Fatal(err)
	}
	i2str, err := i2.String()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "/chain/evm/0xc0ffee254729296a45a3885639AC7E10F9d54979",
		i1str, "ident to string not expected")
	assert.Equal(t, "/storage/ipfs/a2e24ff4/QmXViwQ1frFwabQHtmpt18SUPhnpcRzhWayt9rnTJ8GTay",
		i2str, "ident to string not expected")

	i3 := Ident{
		Type:    "chain",
		SubType: "evm",
		Merkle:  "a2e24ff4",
		Address: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
	}

	i3str, err := i3.String()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "/chain/evm/0xc0ffee254729296a45a3885639AC7E10F9d54979",
		i3str, "ident to string not expected")

	i4 := Ident{
		Type:    "blockchain",
		SubType: "evm",
		Merkle:  "a2e24ff4",
		Address: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
	}
	i4str, err := i4.String()
	if err == nil {
		t.Fatal("ident to string errored")
	}
	assert.Equal(t, "", i4str, "ident to string not expected")
}
