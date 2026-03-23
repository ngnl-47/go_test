package model

// RankRow 榜上一行（查询结果），对应 Lua getRank 返回的每条记录。
type RankRow struct {
	PlayerKey string  // 业务玩家主键（字符串，可与 roleId 一致）
	Score     float64 // 已按 SortOrder 还原的展示分数（降序存储时 Lua/仓储会取反）
	Rank      int     // 1-based
	InfoJSON  string  // 可选，扩展信息 JSON 串；无则空串
}
