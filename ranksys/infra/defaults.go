package infra

import (
	"github.com/redis/go-redis/v9"

	"go_test/ranksys/domain/service"
	"go_test/ranksys/infra/repo"
)

// NewRankBoardStack 组装 key 构建器 + Redis 仓储（调用方再包一层 app/service.RankBoardAppService）。
func NewRankBoardStack(appID string, rdb redis.Cmdable) (*service.RankStorageKeyBuilder, *repo.RedisRankBoardRepository, error) {
	kb, err := service.NewRankStorageKeyBuilder(appID)
	if err != nil {
		return nil, nil, err
	}
	return kb, repo.NewRedisRankBoardRepository(rdb), nil
}
