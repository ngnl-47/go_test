package model

// RankRangeQuery getRank.lua：按 1-based 名次区间拉榜。
type RankRangeQuery struct {
	MinRank     int // 含
	MaxRank     int // 含
	MyPlayerKey string
	IncludeInfo bool
}
