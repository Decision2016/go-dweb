/**
  @author: decision
  @date: 2024/9/2
  @note:
**/

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.io/decision2016/go-dweb/utils"
	"io/fs"
	"path/filepath"
)

var testMerkleDir string

var testCmd = &cobra.Command{
	Use:              "test",
	Short:            "DWeb test commands",
	Long:             "DWeb test commands",
	TraverseChildren: true,
}

var testMerkleCalculate = &cobra.Command{
	Use:   "merkle",
	Short: "Calculate merkle root hash for directory",
	Long:  "Calculate merkle root hash for directory",
	Run: func(cmd *cobra.Command, args []string) {
		cids := make([]string, 0)

		err := filepath.Walk(testMerkleDir, func(path string, info fs.FileInfo,
			err error) error {
			if info.IsDir() {
				return nil
			}
			filename := filepath.Base(path)
			if filename == ".gitignore" {
				return nil
			}

			cid, err := utils.GetFileCidV0(path)
			if err != nil {
				return err
			}

			cids = append(cids, cid.String())
			logrus.Infof("current file: %s cid: %s", filename, cid)
			return nil
		})

		if err != nil {
			logrus.WithError(err).Errorln("walk to obtain file cids failed")
			return
		}

		merkle := utils.MerkleRoot(cids)
		logrus.Infof("directory merkle hash: %s", merkle)
	},
}

func init() {
	testMerkleCalculate.Flags().StringVarP(&testMerkleDir, "dir", "d", "", "source directory path")
	testMerkleCalculate.MarkFlagRequired("dir")

	testCmd.AddCommand(testMerkleCalculate)
}
