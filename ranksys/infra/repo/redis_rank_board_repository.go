package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/redis/go-redis/v9"

	"go_test/ranksys/domain/facade"
	"go_test/ranksys/domain/model"
)

// RedisRankBoardRepository 使用 go-redis Eval 执行三元组 Lua（与 mandarin leaderboard 行为对齐）。
type RedisRankBoardRepository struct {
	rdb redis.Cmdable
}

var _ facade.RankBoardRepository = (*RedisRankBoardRepository)(nil)

func NewRedisRankBoardRepository(rdb redis.Cmdable) *RedisRankBoardRepository {
	return &RedisRankBoardRepository{rdb: rdb}
}

func (r *RedisRankBoardRepository) UploadScore(ctx context.Context, keys model.RankStorageKeyTriplet, in model.UploadScoreInput) (model.UploadScoreResult, error) {
	force := "0"
	if in.ForceReplace {
		force = "1"
	}
	auto := "0"
	if in.AutoIgnoreInfo {
		auto = "1"
	}
	raw, err := r.rdb.Eval(ctx, scriptUpload, keys.EvalKeys(),
		in.PlayerKey,
		in.PrefixedZMember,
		in.RedisSortScore,
		in.InfoJSON,
		in.TimeoutSec,
		force,
		in.SliceKeepTopN,
		auto,
	).Text()
	if err != nil {
		return model.UploadScoreResult{}, err
	}
	var out model.UploadScoreResult
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return model.UploadScoreResult{}, fmt.Errorf("ranksys: decode upload result: %w", err)
	}
	return out, nil
}

type luaRankRow struct {
	Key   string `json:"key"`
	Score any    `json:"score"`
	Rank  int    `json:"rank"`
	Info  any    `json:"info"`
}

func jsonScoreToFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case string:
		return strconv.ParseFloat(x, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unexpected score type %T", v)
	}
}

func infoToString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func (r *RedisRankBoardRepository) GetRanks(ctx context.Context, keys model.RankStorageKeyTriplet, q model.RankRangeQuery) ([]*model.RankRow, *model.RankRow, error) {
	inc := 0
	if q.IncludeInfo {
		inc = 1
	}
	raw, err := r.rdb.Eval(ctx, scriptGetRank, keys.EvalKeys(),
		q.MinRank,
		q.MaxRank,
		q.MyPlayerKey,
		inc,
	).Text()
	if err != nil {
		return nil, nil, err
	}
	var m map[string]*luaRankRow
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, nil, fmt.Errorf("ranksys: decode get_rank: %w", err)
	}

	var me *model.RankRow
	rows := make([]*model.RankRow, 0, len(m))
	for k, row := range m {
		if row == nil {
			continue
		}
		sc, err := jsonScoreToFloat64(row.Score)
		if err != nil {
			return nil, nil, err
		}
		rr := &model.RankRow{
			PlayerKey: row.Key,
			Score:     sc,
			Rank:      row.Rank,
			InfoJSON:  infoToString(row.Info),
		}
		if rr.PlayerKey == "" {
			rr.PlayerKey = k
		}
		if k == q.MyPlayerKey {
			cp := *rr
			me = &cp
		}
		if rr.Rank >= q.MinRank && rr.Rank <= q.MaxRank {
			cp := *rr
			rows = append(rows, &cp)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Rank < rows[j].Rank })
	return rows, me, nil
}

func (r *RedisRankBoardRepository) RemovePlayer(ctx context.Context, keys model.RankStorageKeyTriplet, playerKey string) error {
	_, err := r.rdb.Eval(ctx, scriptRemove, keys.EvalKeys(), playerKey).Result()
	return err
}

func (r *RedisRankBoardRepository) Clear(ctx context.Context, keys model.RankStorageKeyTriplet) error {
	_, err := r.rdb.Eval(ctx, scriptClear, keys.EvalKeys()).Result()
	return err
}
