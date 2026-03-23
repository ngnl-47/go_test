-- leaderboard_addScore：对已在榜玩家做 ZINCRBY（delta 为 Redis 中存储的分值增量）。
redis.replicate_commands()
local key_rank = KEYS[1]
local key_map = KEYS[2]

local playerKey = ARGV[1]
local delta = tonumber(ARGV[2])

local prefixed = redis.call('hget', key_map, playerKey)
if prefixed == false then
    return cjson.encode({ ok = false, err = 'not_found' })
end

redis.call('zincrby', key_rank, delta, prefixed)
local r0 = redis.call('zrank', key_rank, prefixed)
if r0 == false then
    return cjson.encode({ ok = false, err = 'not_in_zset' })
end

local newScore = redis.call('zscore', key_rank, prefixed)
local rank = r0 + 1
local dispScore = tonumber(newScore)
if string.sub(prefixed, 1, 1) == '-' then
    dispScore = -dispScore
end

return cjson.encode({ ok = true, rank = rank, score = dispScore })
