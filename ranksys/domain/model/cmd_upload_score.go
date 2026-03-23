package model

// UploadScoreInput 已对齐 upload.lua 的 ARGV，由应用层从业务请求编码而来。
type UploadScoreInput struct {
	PlayerKey       string
	PrefixedZMember string  // 时间前缀 + playerKey，作为 ZSET member
	RedisSortScore  float64 // 已按降序规则取负后的 Redis 分数
	InfoJSON        string  // 空串表示无 info
	TimeoutSec      int64   // TTL；0 表示不设过期
	ForceReplace    bool
	SliceKeepTopN   int // arg_slice：保留前 N 名，0 表示不裁剪
	AutoIgnoreInfo  bool
}

// UploadScoreResult upload.lua 返回值。
type UploadScoreResult struct {
	Ok    bool    `json:"ok"`
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
}
