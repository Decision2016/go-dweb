/**
  @author: decision
  @date: 2024/6/29
  @note:
**/

package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/interfaces"
	"github.io/decision2016/go-dweb/utils"
	"os"
	"time"
)

type progress struct {
	Commit     string            `json:"commit"`
	Version    int               `json:"version"`
	UpdateTime int64             `json:"update_time"`
	Files      []string          `json:"files"`
	Index      *utils.FullStruct `json:"index"`
}

type Uploader struct {
	storage *interfaces.IFileStorage
	index   *utils.FullStruct

	files  []string
	commit string
}

func NewUploader() *Uploader {
	return nil
}

// upload task
func (u *Uploader) Process(ctx context.Context) error {
	total := len(u.files)
	bar := progressbar.Default(int64(total))

	for idx, file := range u.files {
		err := u.save(u.files[idx+1:])
		if err != nil {
			logrus.WithError(err).Debugln("save progress to disk failed")
			return err
		}

		ident, err := (*u.storage).Upload(ctx, file, file)
		if err != nil {
			logrus.WithError(err).Debugln("error occurred when uploading file")
			return err
		}

		u.index.Paths[file] = ident

		err = bar.Add(1)
		if err != nil {
			logrus.WithError(err).Debugln("add progress failed")
			return err
		}
	}

	return nil
}

func (u *Uploader) save(files []string) error {
	p := progress{
		Commit:     u.commit,
		Version:    0,
		UpdateTime: time.Now().Unix(),
		Files:      files,
		Index:      u.index,
	}

	progressBytes, err := json.Marshal(p)
	if err != nil {
		logrus.WithError(err).Debugln("marshal progress to json failed")
		return err
	}

	err = os.WriteFile(".deploy", progressBytes, 0700)
	if err != nil {
		logrus.WithError(err).Debugln("save deploy progress failed")
		return err
	}

	return nil
}

func (u *Uploader) load() (*progress, error) {
	_, err := os.Stat(".deploy")
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	cachedBytes, err := os.ReadFile(".deploy")
	if err != nil {
		logrus.WithError(err).Debugln("read cached progress file failed")
		return nil, err
	}

	var p progress
	err = json.Unmarshal(cachedBytes, &p)
	if err != nil {
		logrus.WithError(err).Debugln("unmarshal json to progress failed")
		return nil, err
	}

	return &p, nil
}

func (u *Uploader) Setup(index *utils.FullStruct, storage *interfaces.IFileStorage) error {
	u.index = index
	u.commit = index.Commit
	u.storage = storage

	p, err := u.load()
	if err != nil {
		return err
	}
	var opt string
	uploadList := make([]string, 0)

	if p != nil && p.Commit == u.commit {
		fmt.Printf("found existing progress cache,  override? (Y/N): ")
		_, err = fmt.Scan(&opt)
		if err != nil {
			logrus.WithError(err).Debugln("read option failed")
			return err
		}
	}

	if opt == "N" {
		u.index = p.Index
		uploadList = p.Files
		return nil
	}

	for k, _ := range index.Paths {
		if k == "" {
			uploadList = append(uploadList, k)
		}
	}
	u.files = uploadList
	return nil
}
