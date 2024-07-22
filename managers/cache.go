/**
  @author: decision
  @date: 2024/6/7
  @note: cache 文件夹的管理器
**/

package managers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.io/decision2016/go-dweb/utils"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type CacheManager struct {
	indexPath string // index 文件存储全局路径
	appPath   string // app 的一系列文件的全局存储路径
}

var instance *CacheManager
var cacheOnce sync.Once

func CacheDefault() *CacheManager {
	cacheOnce.Do(func() {
		instance = &CacheManager{}
	})

	return instance
}

func (c *CacheManager) Initial() {
	ex, err := os.Executable()
	if err != nil {
		logrus.WithError(err).Fatalln("get current dir path failed")
		return
	}

	path := filepath.Dir(ex)
	c.indexPath = filepath.Join(path, ".service", "index")
	c.appPath = filepath.Join(path, ".service", "app")
	// todo: 如果 service 不存在则进行初始化
	logrus.Infof("service cache manager inited")
}

// Uid 通过哈希求出指定 identity 的唯一标识符
func (c *CacheManager) Uid(identity string) string {
	sha2 := sha256.New()
	sha2.Write([]byte(identity))

	digest := hex.EncodeToString(sha2.Sum(nil))
	return digest[:8]
}

// Validate 检查指定 identity 的目录文件是否正确，完整性校验
func (c *CacheManager) Validate(identity string) (bool, error) {
	exists, err := c.Exists(identity)
	if !exists {
		return false, err
	}

	uid := c.Uid(identity)
	indexPath := filepath.Join(c.indexPath, uid)
	index, err := utils.LoadIndex(indexPath)
	if err != nil {
		return false, err
	}

	merkle := index.Root
	appDir := filepath.Join(c.appPath, uid)
	cids := make([]string, 0)

	err = filepath.Walk(appDir, func(path string, info fs.FileInfo, err error) error {
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
		return nil
	})

	if err != nil {
		return false, err
	}
	dirMerkle := utils.MerkleRoot(cids)
	if dirMerkle != merkle {
		return false, fmt.Errorf("mekrle root hash not equal")
	}

	return true, nil
}

// Exists 检查目录下是否存在对应 identity 的目录
func (c *CacheManager) Exists(identity string) (bool, error) {
	uid := c.Uid(identity)
	index := filepath.Join(c.indexPath, uid)

	_, err := os.Stat(index)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	appDir := filepath.Join(c.appPath, uid)
	_, err = os.Stat(appDir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Path 获取到本地的存储路径-工作目录
func (c *CacheManager) Path(identity string) string {
	uid := c.Uid(identity)
	appDir := filepath.Join(c.appPath, uid)

	return appDir
}

func (c *CacheManager) IndexPath(identity string) string {
	uid := c.Uid(identity)
	appDir := filepath.Join(c.indexPath, uid)

	return appDir
}

// Delete 删除指定 identity 目录
func (c *CacheManager) Delete(identity string) error {
	path := c.Path(identity)

	err := os.RemoveAll(path)
	return err
}
