/**
  @author: decision
  @date: 2024/6/24
  @note:
**/

package utils

import (
	"gopkg.in/yaml.v2"
	"os"
)

type FullStruct struct {
	Commit string            `yaml:"commit"`
	Root   string            `yaml:"root"`
	Paths  map[string]string `yaml:"paths"`
}

func LoadIndex(filepath string) (*FullStruct, error) {
	dataBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var full FullStruct
	err = yaml.Unmarshal(dataBytes, &full)
	if err != nil {
		return nil, err
	}

	return &full, nil
}

func (f *FullStruct) MerkleRoot() {
	cids := make([]string, len(f.Paths))

	for _, path := range f.Paths {
		cids = append(cids, path)
	}

	f.Root = MerkleRoot(cids)
}
