package repo

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go_test/ranksys/domain/model"
	domainsvc "go_test/ranksys/domain/service"
)

// rowsFromLuaRankObjectJSON 解析 Lua cjson 编码的「对象」榜数据（key -> {key,score,rank,info?}）。
func rowsFromLuaRankObjectJSON(raw string, adjustScore func(float64) float64) ([]*model.RankRow, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" {
		return nil, nil
	}
	var m map[string]*luaRankRow
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("ranksys: decode rank map: %w", err)
	}
	if adjustScore == nil {
		adjustScore = func(f float64) float64 { return f }
	}
	out := make([]*model.RankRow, 0, len(m))
	for k, row := range m {
		if row == nil {
			continue
		}
		sc, err := jsonScoreToFloat64(row.Score)
		if err != nil {
			return nil, err
		}
		rr := &model.RankRow{
			PlayerKey: row.Key,
			Score:     adjustScore(sc),
			Rank:      row.Rank,
			InfoJSON:  infoToString(row.Info),
		}
		if rr.PlayerKey == "" {
			rr.PlayerKey = k
		}
		out = append(out, rr)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Rank < out[j].Rank })
	return out, nil
}

func rowsFromScoreQueryJSON(raw string, sortOrder model.SortOrder) ([]*model.RankRow, error) {
	return rowsFromLuaRankObjectJSON(raw, func(s float64) float64 {
		return domainsvc.DisplayScoreFromRedis(sortOrder, s)
	})
}

// parseMyPageJSON get_my_page：未上榜返回 "[]"，上榜返回与 get_rank 相同的对象 JSON。
func parseMyPageJSON(raw string) ([]*model.RankRow, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "[") {
		return nil, nil
	}
	return rowsFromLuaRankObjectJSON(raw, nil)
}

type luaRankByIDRow struct {
	RoleIdStr string `json:"roleIdStr"`
	Score     any    `json:"score"`
	Rank      int    `json:"rank"`
	Info      any    `json:"info"`
}

// parseRanksByPlayerIDsJSON 与 get_ranks_by_ids.lua 返回的 JSON 数组对齐（含 null）。
func parseRanksByPlayerIDsJSON(raw string, ignoreNil bool) ([]*model.RankRow, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil, fmt.Errorf("ranksys: decode ids array: %w", err)
	}
	if ignoreNil {
		out := make([]*model.RankRow, 0, len(arr))
		for i, rm := range arr {
			s := strings.TrimSpace(string(rm))
			if s == "null" || s == "" {
				continue
			}
			rr, err := decodeRankByIDRow(i, rm)
			if err != nil {
				return nil, err
			}
			out = append(out, rr)
		}
		return out, nil
	}
	out := make([]*model.RankRow, len(arr))
	for i, rm := range arr {
		s := strings.TrimSpace(string(rm))
		if s == "null" || s == "" {
			out[i] = nil
			continue
		}
		rr, err := decodeRankByIDRow(i, rm)
		if err != nil {
			return nil, err
		}
		out[i] = rr
	}
	return out, nil
}

func decodeRankByIDRow(i int, rm json.RawMessage) (*model.RankRow, error) {
	var row luaRankByIDRow
	if err := json.Unmarshal(rm, &row); err != nil {
		return nil, fmt.Errorf("ranksys: decode ids row %d: %w", i, err)
	}
	sc, err := jsonScoreToFloat64(row.Score)
	if err != nil {
		return nil, err
	}
	return &model.RankRow{
		PlayerKey: row.RoleIdStr,
		Score:     sc,
		Rank:      row.Rank,
		InfoJSON:  infoToString(row.Info),
	}, nil
}
