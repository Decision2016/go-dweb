package web

type LoadTask struct {
	Type    ChainType
	Address string
}

// Loader 接收任务从链上拉取配置信息，然后获取状态
type Loader struct {
}
