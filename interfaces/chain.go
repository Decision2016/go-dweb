/**
  @author: decision
  @date: 2024/6/5
  @note: 区块链数据访问接口，用于插件的实现
**/

package interfaces

type IChain interface {
	// （读取）FS 静态资源存储相关，获取相关的属性信息
	Identity() (string, error)  // 获取链上存放的索引信息，格式为 /type/subtype/version/cid
	Bootstrap() (string, error) // 获取 P2P 网络连接信息，如果为空则非 P2P 应用

	// 写入数据到链上
	Initial(ident string, url string) error
	SetIdentity(ident string) error
	Join(url string) error // 为 P2P 网络提供到的链上数据扩展

	// 当前类的基本操作
	Setup(address string) error
}
