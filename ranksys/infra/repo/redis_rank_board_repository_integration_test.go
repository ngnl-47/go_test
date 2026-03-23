// Redis 需支持 Lua cjson（与线上 mandarin 环境一致）；未设置 RANKSYS_REDIS_ADDR 时跳过。
package repo

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"go_test/ranksys/domain/model"
	domainsvc "go_test/ranksys/domain/service"
)

func redisIntegrationAddr(t *testing.T) string {
	addr := os.Getenv("RANKSYS_REDIS_ADDR")
	if addr == "" {
		t.Skip("set RANKSYS_REDIS_ADDR (e.g. 127.0.0.1:6379) for Redis+Lua+cjson integration tests")
	}
	return addr
}

// ascMember 升序榜 ZSET member：+ + 13 位毫秒 + ':' + playerKey（Lua 从第 16 字节截取 id）。
func ascMember(playerKey string, now time.Time) string {
	ms := now.UnixMilli()
	p := strconv.FormatInt(ms, 10)
	for len(p) < 13 {
		p = "0" + p
	}
	if len(p) > 13 {
		p = p[len(p)-13:]
	}
	return "+" + p + ":" + playerKey
}

func TestRedisRankBoardRepository_MigratedLua_Integration(t *testing.T) {
	addr := redisIntegrationAddr(t)
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { _ = rdb.Close() })
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis ping: %v", err)
	}

	kb, err := domainsvc.NewRankStorageKeyBuilder("inttest")
	if err != nil {
		t.Fatal(err)
	}
	board, err := model.ParseBoardName("migrate_lua")
	if err != nil {
		t.Fatal(err)
	}
	keys := kb.Build(board)
	repo := NewRedisRankBoardRepository(rdb)

	if err := repo.Clear(ctx, keys); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Clear(ctx, keys) })

	now := time.Unix(1_700_000_000, 0).UTC()
	p1 := ascMember("p1", now)
	p2 := ascMember("p2", now.Add(time.Millisecond))
	p3 := ascMember("p3", now.Add(2*time.Millisecond))

	_, err = repo.UploadScore(ctx, keys, model.UploadScoreInput{
		PlayerKey: "p1", PrefixedZMember: p1, RedisSortScore: 100, InfoJSON: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.UploadScore(ctx, keys, model.UploadScoreInput{
		PlayerKey: "p2", PrefixedZMember: p2, RedisSortScore: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.UploadScore(ctx, keys, model.UploadScoreInput{
		PlayerKey: "p3", PrefixedZMember: p3, RedisSortScore: 75,
	})
	if err != nil {
		t.Fatal(err)
	}

	byID, err := repo.GetRanksByPlayerIDs(ctx, keys, model.RanksByPlayerIDsInput{
		PlayerKeys: []string{"p1", "p2", "missing"}, IncludeInfo: false, IgnoreNil: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(byID) != 3 || byID[2] != nil {
		t.Fatalf("by id: len=%d third=%v", len(byID), byID[2])
	}

	byScore, err := repo.GetRanksByScoreRange(ctx, keys, model.RankScoreRangeInput{
		MinRedisScore: 50,
		MaxRedisScore: 80,
		Skip:          0,
		Limit:         0,
		IncludeInfo:   false,
		Sort:          model.SortAsc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(byScore) != 2 {
		t.Fatalf("by score want 2 got %d %+v", len(byScore), byScore)
	}

	p2new := ascMember("p2", now.Add(time.Hour))
	res, err := repo.AddScore(ctx, keys, model.AddScoreInput{
		PlayerKey: "p2", PrefixedZMember: p2new, RedisScoreDelta: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Ok || res.Score != 55 {
		t.Fatalf("add score: %+v", res)
	}

	pageRows, err := repo.GetMyPage(ctx, keys, model.MyPageQuery{
		MyPlayerKey: "p3", PageSize: 10, IncludeInfo: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pageRows) != 3 {
		t.Fatalf("page rows want 3 got %d %+v", len(pageRows), pageRows)
	}

	anchorRows, err := repo.GetRanksByScoreAnchor(ctx, keys, model.RankScoreAnchorInput{
		RedisAnchorScore: 75,
		High:             2,
		Low:              2,
		IncludeInfo:      false,
		Sort:             model.SortAsc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(anchorRows) < 1 {
		t.Fatalf("anchor rows: %+v", anchorRows)
	}
}
