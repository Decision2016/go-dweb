package web

type AppIndex struct {
	Commit string            `json:"commit"` // 应用工作目录下 git 的 commit 哈希值
	Maps   map[string]string `json:"maps"`   // 相对目录 -> cid 的索引信息，后续可以考虑进行优化
}
