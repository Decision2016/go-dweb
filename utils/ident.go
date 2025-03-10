/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

type Ident struct {
	Type    string
	SubType string
	Merkle  string
	Address string
}

// String ident 转换为字符串类型
func (ident *Ident) String() (string, error) {
	if ident.Type == "chain" {
		return filepath.Join("/", ident.Type, ident.SubType, ident.Address), nil
	} else if ident.Type == "storage" {
		return filepath.Join("/", ident.Type, ident.SubType,
			ident.Merkle, ident.Address), nil
	} else {
		return "", fmt.Errorf("ident invalid")
	}
}

// FromString 从字符串转换得到 ident
func (ident *Ident) FromString(path string) error {
	items := strings.Split(path, "/")
	if len(items) < 4 || (items[1] != "storage" && items[1] != "chain") {
		return fmt.Errorf("ident path invalid")
	}

	ident.Type = items[1]
	ident.SubType = items[2]
	if ident.Type == "storage" {
		ident.Merkle = items[3]
		ident.Address = items[4]
	} else {
		ident.Address = items[3]
	}
	return nil
}

// Uid 获取 cid 的哈希
func (ident *Ident) Uid() (string, error) {
	identStr, err := ident.String()
	if err != nil {
		return "", err
	}

	sha2 := sha256.New()
	sha2.Write([]byte(identStr))

	digest := hex.EncodeToString(sha2.Sum(nil))
	return digest[:8], nil

}
