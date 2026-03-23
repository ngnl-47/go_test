# ranksys — 排行榜三元 Redis Key（DDD 提炼）

从 `mandarin/v3/modules/leaderboard` 中提炼 **ZSET + id→member Hash + info Hash** 的命名规则与核心 Lua，按 [ddd-package-practices.md](../../thoughtnotes/golang/ddd/ddd-package-practices.md) 组织（与 `go_test` 同属 `workspace_other` 时相对路径如此）。

## 领域含义

| 概念 | 位置 | 说明 |
|------|------|------|
| `model.BoardName` | `domain/model` | 排行榜逻辑名，写入 `leaderboard:{name}` 的 **hash tag**，保证三元 key 同槽 |
| `model.RankStorageKeyTriplet` | `domain/model` | 三个 Redis key 值对象 |
| `service.RankStorageKeyBuilder` | `domain/service` | `appId:leaderboard:{name}:rank|:map|:info` |
| `facade.RankBoardRepository` | `domain/facade` | 上传/查榜/删人/清空（无 Redis 类型泄漏） |
| `app/service.RankBoardAppService` | `app/service` | 用例编排：编码时间前缀、分数符号、TTL |
| `infra/repo` | `infra/repo` | `embed` 的 `upload/get_rank/remove/clear` Lua + `RedisRankBoardRepository` |

## 依赖

- `github.com/redis/go-redis/v9`（仅 `infra/repo`）

## 使用示例

```go
kb, rrepo, err := infra.NewRankBoardStack("myapp", rdb)
if err != nil { ... }
appSvc := appsvc.NewRankBoardAppService(kb, rrepo)

name, _ := model.ParseBoardName("rank:server:level:{1}")
res, err := appSvc.UploadScore(ctx, name, appsvc.UploadScoreRequest{
    PlayerID: 1001,
    Sort:     model.SortDesc,
    Score:    99.5,
}, time.Now())
```

## 与 mandarin 的差异

- 未迁移：`GetRanksByScore`、`AddScore`、`GetRanksByPlayerIds`、`GetMyPage` 等脚本；需要时可按同样方式增加 `domain/facade` 方法与 `infra/repo/lua`。
