package model

import "strings"

// RankStorageKeyTriplet 排行榜在 Redis 上的三元存储：ZSET（带时间前缀的 member）、
// Hash（playerKey -> zset member）、Hash（playerKey -> info JSON）。
// 三者的 key 字符串共享同一段 leaderboard:{name}，满足 Redis Cluster 单槽约束。
type RankStorageKeyTriplet struct {
	rankZSetKey string // KEYS[1] — 有序集合
	idToZMember string // KEYS[2] — 玩家 id -> ZSET member 映射
	infoHashKey string // KEYS[3] — 玩家 id -> 展示信息 JSON
}

// NewRankStorageKeyTriplet 由基础设施或测试构造；业务侧应通过 RankStorageKeyBuilder 生成。
func NewRankStorageKeyTriplet(rankZSet, idToZMember, infoHash string) RankStorageKeyTriplet {
	return RankStorageKeyTriplet{
		rankZSetKey: rankZSet,
		idToZMember: idToZMember,
		infoHashKey: infoHash,
	}
}

func (t RankStorageKeyTriplet) RankZSetKey() string { return t.rankZSetKey }

func (t RankStorageKeyTriplet) IDToZMemberMapKey() string { return t.idToZMember }

func (t RankStorageKeyTriplet) InfoHashKey() string { return t.infoHashKey }

// EvalKeys 三元组 Lua（upload/getRank/remove/clear）的 KEYS 顺序。
func (t RankStorageKeyTriplet) EvalKeys() []string {
	return []string{t.rankZSetKey, t.idToZMember, t.infoHashKey}
}

// LuaSingleKeyScriptsBaseKey 与 mandarin 单 KEYS[1] 脚本一致：appId:leaderboard:{name}（Lua 内再拼 :rank / :map / :info）。
func (t RankStorageKeyTriplet) LuaSingleKeyScriptsBaseKey() string {
	rk := t.rankZSetKey
	if strings.HasSuffix(rk, ":rank") {
		return strings.TrimSuffix(rk, ":rank")
	}
	return rk
}

// BaseKey 同 LuaSingleKeyScriptsBaseKey，便于测试与调用方命名。
func (t RankStorageKeyTriplet) BaseKey() string { return t.LuaSingleKeyScriptsBaseKey() }
