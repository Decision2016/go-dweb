/**
  @author: decision
  @date: 2024/6/29
  @note:
**/

package managers

import "github.io/decision2016/go-dweb/interfaces"

type Uploader struct {
	storage *interfaces.IFileStorage

	files []string
}

func NewUploader() *Uploader {
	return nil
}

// upload task
func (u *Uploader) process() {

}

func (u *Uploader) save() error {

}

func (u *Uploader) load() error {

}

func (u *Uploader) Setup(files []string) error {

}
