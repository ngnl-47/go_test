package model

// AddScoreInput 对齐 add.lua 的 ARGV（分数增量已按 SortOrder 换算为 Redis 存储分）。
type AddScoreInput struct {
	PlayerKey         string
	PrefixedZMember   string
	RedisScoreDelta   float64
	InfoJSON          string
	TimeoutSec        int64
	SliceKeepTopN     int
}
