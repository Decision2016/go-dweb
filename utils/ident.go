/**
  @author: decision
  @date: 2024/7/1
  @note:
**/

package utils

import (
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
