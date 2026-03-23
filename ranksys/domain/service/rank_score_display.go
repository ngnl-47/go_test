package service

import "go_test/ranksys/domain/model"

// DisplayScoreFromRedis 将 get_ranks_by_score_* 脚本返回的原始分转为展示分（与 mandarin GetRanksByScore 一致）。
func DisplayScoreFromRedis(sort model.SortOrder, redisScore float64) float64 {
	if sort != model.SortDesc {
		return -redisScore
	}
	return redisScore
}

// RedisAnchorScore 将「展示锚点分」转为传给 anchor 脚本的分数（与 mandarin 一致）。
func RedisAnchorScore(sort model.SortOrder, displayScore float64) float64 {
	if sort != model.SortDesc {
		return -displayScore
	}
	return displayScore
}

// RedisScoreRange 将展示分闭区间 [min,max] 转为 ZRANGEBYSCORE 参数（与 mandarin 一致）。
func RedisScoreRange(sort model.SortOrder, minDisplay, maxDisplay float64) (minRedis, maxRedis float64, ok bool) {
	if minDisplay > maxDisplay {
		return 0, 0, false
	}
	if sort == model.SortDesc {
		return minDisplay, maxDisplay, true
	}
	return -maxDisplay, -minDisplay, true
}
