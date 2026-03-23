package model

import (
	"errors"
	"strings"
)

// BoardName 排行榜逻辑名（与 mandarin leaderboard.Name 一致），进入 Redis 后包在
// appId:leaderboard:{BoardName} 的 hash tag 中，保证 rank/map/info 三 key 同槽。
type BoardName string

var (
	ErrEmptyBoardName = errors.New("ranksys: board name is empty")
)

// ParseBoardName 校验并规范化排行榜名（去首尾空白）。
func ParseBoardName(s string) (BoardName, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", ErrEmptyBoardName
	}
	return BoardName(t), nil
}

func (n BoardName) String() string { return string(n) }
