package service

import (
	"errors"
	"strings"

	"go_test/ranksys/domain/model"
)

// RankStorageKeyBuilder 根据 appId + 排行榜逻辑名构造三元 Redis key（与 mandarin redisKeys 一致）。
type RankStorageKeyBuilder struct {
	appID string
}

var ErrEmptyAppID = errors.New("ranksys: app id is empty")

func NewRankStorageKeyBuilder(appID string) (*RankStorageKeyBuilder, error) {
	if strings.TrimSpace(appID) == "" {
		return nil, ErrEmptyAppID
	}
	return &RankStorageKeyBuilder{appID: strings.TrimSpace(appID)}, nil
}

// MustNewRankStorageKeyBuilder 用于示例/测试；生产环境建议用 New 并处理错误。
func MustNewRankStorageKeyBuilder(appID string) *RankStorageKeyBuilder {
	b, err := NewRankStorageKeyBuilder(appID)
	if err != nil {
		panic(err)
	}
	return b
}

// Build 生成 appId:leaderboard:{name}:rank | :map | :info。
func (b *RankStorageKeyBuilder) Build(name model.BoardName) model.RankStorageKeyTriplet {
	base := b.appID + ":leaderboard:{" + name.String() + "}"
	return model.NewRankStorageKeyTriplet(
		base+":rank",
		base+":map",
		base+":info",
	)
}
