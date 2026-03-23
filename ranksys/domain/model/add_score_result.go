package model

// AddScoreResult add.lua 返回值。
type AddScoreResult struct {
	Ok    bool    `json:"ok"`
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
	Err   string  `json:"err,omitempty"`
}
