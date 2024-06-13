/**
  @author: decision
  @date: 2024/6/5
  @note: 分布式数据存取接口，用于静态资源存储，插件接口的定义
**/

package interfaces

import "context"

type IFileStorage interface {
	Initial(ctx context.Context) error                               // 基于配置文件进行初始化并启动实例
	Ping(ctx context.Context) error                                  // 测试接口的可用性
	Exists(ctx context.Context, source string) (bool, error)         // 检查文件在 FS 中是否存在
	Upload(ctx context.Context, name string, source string) error    // 上传文件
	Download(ctx context.Context, identity string, dst string) error // 加载文件
	Delete(ctx context.Context, identity string) error               // 在实际的应用下 IPFS 中不能保证完全删除文件，所以通常需要进行增量更新
	// Start(ctx context.Context)
	// fixed: 由接口的实现类自行实现内置 start 方法
	// 实例运行协程，例如本地 FS 模式下的 IPFS 内置节点
}
