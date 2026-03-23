package facade

import (
	"context"

	"go_test/ranksys/domain/model"
)

// RankBoardRepository 排行榜三元组 + Lua 脚本端口（领域不感知 Redis 类型细节）。
type RankBoardRepository interface {
	UploadScore(ctx context.Context, keys model.RankStorageKeyTriplet, in model.UploadScoreInput) (model.UploadScoreResult, error)
	GetRanks(ctx context.Context, keys model.RankStorageKeyTriplet, q model.RankRangeQuery) (rows []*model.RankRow, me *model.RankRow, err error)
	RemovePlayer(ctx context.Context, keys model.RankStorageKeyTriplet, playerKey string) error
	Clear(ctx context.Context, keys model.RankStorageKeyTriplet) error
}
