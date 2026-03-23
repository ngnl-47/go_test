package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_test/ranksys/domain/facade"
	"go_test/ranksys/domain/model"
	domainsvc "go_test/ranksys/domain/service"
)

var ErrInvalidScoreRange = errors.New("ranksys: min display score > max display score")

// UploadScoreRequest 业务层上传分数请求（应用服务负责编码为 model.UploadScoreInput）。
type UploadScoreRequest struct {
	PlayerID       int64
	PlayerKey      string // 非空时优先于 PlayerID
	Sort           model.SortOrder
	Score          float64
	Info           any // nil 或 json.Marshaler；写入 info Hash
	ExpiredAtUnix  int64 // 0 表示不过期；秒级时间戳
	ForceReplace   bool
	SliceKeepTopN  int
	AutoIgnoreInfo bool
}

// AddScoreRequest 累加分（与 mandarin AddScoreArgs 对齐）。
type AddScoreRequest struct {
	PlayerID      int64
	PlayerKey     string
	Sort          model.SortOrder
	ScoreAdd      float64
	Info          any
	ExpiredAtUnix int64
	SliceKeepTopN int
}

func (r UploadScoreRequest) resolvedPlayerKey() string {
	if strings.TrimSpace(r.PlayerKey) != "" {
		return strings.TrimSpace(r.PlayerKey)
	}
	return strconv.FormatInt(r.PlayerID, 10)
}

func (r AddScoreRequest) resolvedPlayerKey() string {
	if strings.TrimSpace(r.PlayerKey) != "" {
		return strings.TrimSpace(r.PlayerKey)
	}
	return strconv.FormatInt(r.PlayerID, 10)
}

// RankBoardAppService 编排：构造三元 key、编码 Lua 参数、调用仓储。
type RankBoardAppService struct {
	keys *domainsvc.RankStorageKeyBuilder
	repo facade.RankBoardRepository
}

func NewRankBoardAppService(keys *domainsvc.RankStorageKeyBuilder, repo facade.RankBoardRepository) *RankBoardAppService {
	return &RankBoardAppService{keys: keys, repo: repo}
}

func (s *RankBoardAppService) buildUploadInput(req UploadScoreRequest, now time.Time) (model.UploadScoreInput, error) {
	playerKey := req.resolvedPlayerKey()
	score := req.Score
	if req.Sort.IsDesc() {
		score = -score
	}
	prefix := now.UnixMilli()
	prefixS := fmt.Sprintf("%d", prefix)
	if len(prefixS) < 13 {
		prefixS = strings.Repeat("0", 13-len(prefixS)) + prefixS
	}
	if req.Sort.IsDesc() {
		prefixS = "-" + prefixS
	} else {
		prefixS = "+" + prefixS
	}
	prefixedMember := prefixS + ":" + playerKey

	var infoJSON string
	if req.Info != nil {
		b, err := json.Marshal(req.Info)
		if err != nil {
			return model.UploadScoreInput{}, fmt.Errorf("ranksys: marshal info: %w", err)
		}
		infoJSON = string(b)
	}

	timeout := int64(0)
	if req.ExpiredAtUnix != 0 {
		sec := int64(time.Until(time.Unix(req.ExpiredAtUnix, 0)).Seconds())
		if sec <= 0 {
			sec = 1
		}
		timeout = sec
	}

	return model.UploadScoreInput{
		PlayerKey:       playerKey,
		PrefixedZMember: prefixedMember,
		RedisSortScore:  score,
		InfoJSON:        infoJSON,
		TimeoutSec:      timeout,
		ForceReplace:    req.ForceReplace,
		SliceKeepTopN:   req.SliceKeepTopN,
		AutoIgnoreInfo:  req.AutoIgnoreInfo,
	}, nil
}

func (s *RankBoardAppService) buildAddScoreInput(req AddScoreRequest, now time.Time) (model.AddScoreInput, error) {
	playerKey := req.resolvedPlayerKey()
	delta := req.ScoreAdd
	if req.Sort.IsDesc() {
		delta = -delta
	}
	prefix := now.UnixMilli()
	prefixS := fmt.Sprintf("%d", prefix)
	if len(prefixS) < 13 {
		prefixS = strings.Repeat("0", 13-len(prefixS)) + prefixS
	}
	if req.Sort.IsDesc() {
		prefixS = "-" + prefixS
	} else {
		prefixS = "+" + prefixS
	}
	prefixedMember := prefixS + ":" + playerKey

	var infoJSON string
	if req.Info != nil {
		b, err := json.Marshal(req.Info)
		if err != nil {
			return model.AddScoreInput{}, fmt.Errorf("ranksys: marshal info: %w", err)
		}
		infoJSON = string(b)
	}

	timeout := int64(0)
	if req.ExpiredAtUnix != 0 {
		sec := int64(time.Until(time.Unix(req.ExpiredAtUnix, 0)).Seconds())
		if sec <= 0 {
			sec = 1
		}
		timeout = sec
	}

	return model.AddScoreInput{
		PlayerKey:       playerKey,
		PrefixedZMember: prefixedMember,
		RedisScoreDelta: delta,
		InfoJSON:        infoJSON,
		TimeoutSec:      timeout,
		SliceKeepTopN:   req.SliceKeepTopN,
	}, nil
}

// UploadScore 上传或更新分数（语义与 mandarin leaderboard.UploadScoreEx 对齐）。
func (s *RankBoardAppService) UploadScore(ctx context.Context, board model.BoardName, req UploadScoreRequest, now time.Time) (model.UploadScoreResult, error) {
	in, err := s.buildUploadInput(req, now)
	if err != nil {
		return model.UploadScoreResult{}, err
	}
	keys := s.keys.Build(board)
	return s.repo.UploadScore(ctx, keys, in)
}

// AddScore 累加分数（语义与 mandarin leaderboard.AddScore 对齐）。
func (s *RankBoardAppService) AddScore(ctx context.Context, board model.BoardName, req AddScoreRequest, now time.Time) (model.AddScoreResult, error) {
	in, err := s.buildAddScoreInput(req, now)
	if err != nil {
		return model.AddScoreResult{}, err
	}
	keys := s.keys.Build(board)
	return s.repo.AddScore(ctx, keys, in)
}

// GetRanks 按名次区间查询。
func (s *RankBoardAppService) GetRanks(ctx context.Context, board model.BoardName, q model.RankRangeQuery) ([]*model.RankRow, *model.RankRow, error) {
	keys := s.keys.Build(board)
	return s.repo.GetRanks(ctx, keys, q)
}

// GetRanksByScoreAnchor 以展示分为锚点上下扩展（展示分 / Sort 与 mandarin 一致）。
func (s *RankBoardAppService) GetRanksByScoreAnchor(ctx context.Context, board model.BoardName, displayScore float64, high, low int, includeInfo bool, sort model.SortOrder) ([]*model.RankRow, error) {
	keys := s.keys.Build(board)
	in := model.RankScoreAnchorInput{
		RedisAnchorScore: domainsvc.RedisAnchorScore(sort, displayScore),
		High:             high,
		Low:              low,
		IncludeInfo:      includeInfo,
		Sort:             sort,
	}
	return s.repo.GetRanksByScoreAnchor(ctx, keys, in)
}

// GetRanksByScoreRange 按展示分闭区间查询。
func (s *RankBoardAppService) GetRanksByScoreRange(ctx context.Context, board model.BoardName, minDisplay, maxDisplay float64, skip, limit int, includeInfo bool, sort model.SortOrder) ([]*model.RankRow, error) {
	minR, maxR, ok := domainsvc.RedisScoreRange(sort, minDisplay, maxDisplay)
	if !ok {
		return nil, ErrInvalidScoreRange
	}
	keys := s.keys.Build(board)
	in := model.RankScoreRangeInput{
		MinRedisScore: minR,
		MaxRedisScore: maxR,
		Skip:          skip,
		Limit:         limit,
		IncludeInfo:   includeInfo,
		Sort:          sort,
	}
	return s.repo.GetRanksByScoreRange(ctx, keys, in)
}

// GetRanksByPlayerIDs 按玩家 id 批量查榜上行（顺序与入参一致；IgnoreNil 为 true 时去掉未上榜项）。
func (s *RankBoardAppService) GetRanksByPlayerIDs(ctx context.Context, board model.BoardName, in model.RanksByPlayerIDsInput) ([]*model.RankRow, error) {
	keys := s.keys.Build(board)
	return s.repo.GetRanksByPlayerIDs(ctx, keys, in)
}

// GetMyPage 返回当前玩家所在「页」内的排行（PageSize 与 mandarin 一致）。
func (s *RankBoardAppService) GetMyPage(ctx context.Context, board model.BoardName, q model.MyPageQuery) ([]*model.RankRow, error) {
	keys := s.keys.Build(board)
	return s.repo.GetMyPage(ctx, keys, q)
}

// RemovePlayer 从榜中移除玩家。
func (s *RankBoardAppService) RemovePlayer(ctx context.Context, board model.BoardName, playerKey string) error {
	keys := s.keys.Build(board)
	return s.repo.RemovePlayer(ctx, keys, playerKey)
}

// Clear 清空整张榜（三元 key 全部删除）。
func (s *RankBoardAppService) Clear(ctx context.Context, board model.BoardName) error {
	keys := s.keys.Build(board)
	return s.repo.Clear(ctx, keys)
}
