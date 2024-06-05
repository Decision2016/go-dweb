/**
  @author: decision
  @date: 2024/6/5
  @note: 区块链数据访问接口，用于插件的实现
**/

package interfaces

type ChainInterface interface {
	Upload(name string, data []byte) error
	Download(identity string) ([]byte, error)
	// 在实际的应用下 IPFS 中不能保证完全删除文件，所以通常需要进行增量更新
	Delete(identity string) error
}
