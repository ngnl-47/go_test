package model

// SortOrder 排行榜排序方向（与 mandarin SortType 对应）。
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

func (o SortOrder) IsDesc() bool { return o == SortDesc }
