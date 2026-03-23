package model

// RankScoreAnchorInput get_ranks_by_score_anchor.lua：以某分为锚点向上下扩展名次。
type RankScoreAnchorInput struct {
	RedisAnchorScore float64 // 已按降序规则取负后的锚点分
	High             int     // 比锚点更高名次侧条数
	Low              int     // 比锚点更低名次侧条数
	IncludeInfo      bool
	Sort             SortOrder // 返回展示分时的符号修正（与 mandarin 一致）
}

// RankScoreRangeInput get_ranks_by_score_range.lua：按 Redis ZSET 分数闭区间拉榜。
type RankScoreRangeInput struct {
	MinRedisScore float64
	MaxRedisScore float64
	Skip          int
	Limit         int // 0 表示 lua 中不 limit
	IncludeInfo   bool
	Sort          SortOrder
}

// RanksByPlayerIDsInput get_ranks_by_ids.lua。
type RanksByPlayerIDsInput struct {
	PlayerKeys  []string
	IncludeInfo bool
	IgnoreNil   bool // true 时结果只含在榜上的行，且顺序为追加顺序非输入顺序
}
