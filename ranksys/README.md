# ranksys — 排行榜三元 Redis Key（DDD 提炼）

从 `mandarin/v3/modules/leaderboard` 中提炼 **ZSET + id→member Hash + info Hash** 的命名规则与 Lua，按 [ddd-package-practices.md](../../thoughtnotes/golang/ddd/ddd-package-practices.md) 组织（与 `go_test` 同属 `workspace_other` 时相对路径如此）。

## 领域含义

| 概念 | 位置 | 说明 |
|------|------|------|
| `model.BoardName` | `domain/model` | 排行榜逻辑名，写入 `leaderboard:{name}` 的 **hash tag** |
| `model.RankStorageKeyTriplet` | `domain/model` | 三元 key；`BaseKey()` 供单 KEY 脚本（内部拼 `:rank`/`:map`/`:info`） |
| `service.RankStorageKeyBuilder` | `domain/service` | 构造三元组 |
| `service.DisplayScoreFromRedis` 等 | `domain/service` | 展示分 ↔ Redis 分（与 mandarin 一致） |
| `facade.RankBoardRepository` | `domain/facade` | 仓储端口 |
| `app/service.RankBoardAppService` | `app/service` | 用例编排（时间前缀、分数、TTL、区间换算） |
| `infra/repo` | `infra/repo` | embed Lua + `RedisRankBoardRepository` |

## 已迁移 Lua

| 脚本 | 用途 |
|------|------|
| `upload.lua` | 上传/覆盖分 |
| `get_rank.lua` | 名次区间 |
| `add.lua` | 累加分 |
| `get_ranks_by_score_anchor.lua` | 按分锚点上下扩展 |
| `get_ranks_by_score_range.lua` | 按分数区间 |
| `get_ranks_by_ids.lua` | 批量 id |
| `get_my_page.lua` | 按 PageSize 分页的一页 |
| `remove.lua` / `clear.lua` | 删人 / 清空 |

## 依赖

- `github.com/redis/go-redis/v9`（仅 `infra/repo`）

## 使用示例

```go
kb, rrepo, err := infra.NewRankBoardStack("myapp", rdb)
appSvc := appsvc.NewRankBoardAppService(kb, rrepo)

name, _ := model.ParseBoardName("rank:server:level:{1}")
_, _ = appSvc.UploadScore(ctx, name, appsvc.UploadScoreRequest{
    PlayerID: 1001, Sort: model.SortDesc, Score: 99.5,
}, time.Now())

_, _ = appSvc.AddScore(ctx, name, appsvc.AddScoreRequest{
    PlayerID: 1001, Sort: model.SortDesc, ScoreAdd: 1,
}, time.Now())

_, _ = appSvc.GetRanksByScoreAnchor(ctx, name, 100, 5, 5, true, model.SortDesc)
_, _ = appSvc.GetRanksByScoreRange(ctx, name, 10, 200, 0, 10, true, model.SortDesc)
_, _ = appSvc.GetRanksByPlayerIDs(ctx, name, model.RanksByPlayerIDsInput{
    PlayerKeys: []string{"1", "2"}, IncludeInfo: true, IgnoreNil: false,
})
_, _ = appSvc.GetMyPage(ctx, name, model.MyPageQuery{MyPlayerKey: "1", PageSize: 20, IncludeInfo: true})
```

## 测试

- **单元**：`domain/service/*_test.go`（key、展示分换算）、`domain/model/*_test.go`、`infra/repo/parse_test.go`（Lua JSON 解析）、`app/service/*_test.go`（请求编码 + mock 仓储）。
- **集成**：设置 `RANKSYS_REDIS_ADDR`（需支持 Lua **cjson**，与线上一致）后运行 `go test ./ranksys/infra/repo/ -run Integration`。
