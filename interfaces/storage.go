/**
  @author: decision
  @date: 2024/6/5
  @note: 分布式数据存取接口，用于静态资源存储，插件接口的定义
**/

package interfaces

type FSInterface interface {
	Upload(name string, data []byte) error
	Download(identity string) ([]byte, error)
	// 在实际的应用下 IPFS 中不能保证完全删除文件，所以通常需要进行增量更新
	Delete(identity string) error
}
