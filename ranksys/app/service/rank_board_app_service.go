package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_test/ranksys/domain/facade"
	"go_test/ranksys/domain/model"
	domainsvc "go_test/ranksys/domain/service"
)

// UploadScoreRequest 业务层上传分数请求（应用服务负责编码为 model.UploadScoreInput）。
type UploadScoreRequest struct {
	PlayerID       int64
	PlayerKey      string // 非空时优先于 PlayerID
	Sort           model.SortOrder
	Score          float64
	Info           any   // nil 或 json.Marshaler；写入 info Hash
	ExpiredAtUnix  int64 // 0 表示不过期；秒级时间戳
	ForceReplace   bool
	SliceKeepTopN  int
	AutoIgnoreInfo bool
}

func (r UploadScoreRequest) resolvedPlayerKey() string {
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

// UploadScore 上传或更新分数（语义与 mandarin leaderboard.UploadScoreEx 对齐）。
func (s *RankBoardAppService) UploadScore(ctx context.Context, board model.BoardName, req UploadScoreRequest, now time.Time) (model.UploadScoreResult, error) {
	in, err := s.buildUploadInput(req, now)
	if err != nil {
		return model.UploadScoreResult{}, err
	}
	keys := s.keys.Build(board)
	return s.repo.UploadScore(ctx, keys, in)
}

// GetRanks 按名次区间查询。
func (s *RankBoardAppService) GetRanks(ctx context.Context, board model.BoardName, q model.RankRangeQuery) ([]*model.RankRow, *model.RankRow, error) {
	keys := s.keys.Build(board)
	return s.repo.GetRanks(ctx, keys, q)
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
