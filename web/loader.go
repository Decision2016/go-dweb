package web

type LoadTask struct {
	Type    ChainType
	Address string
}

// Loader 接收任务从链上拉取配置信息，然后获取状态
type Loader struct {
	queue chan LoadTask
}

func (l *Loader) Run() {
	for {
		select {
		case task := <-l.queue:
			task.Type = 1
		}
	}
}

func (l *Loader) AppendTask(task LoadTask) {
	l.queue <- task
}
