--leaderboard_get_ranks_by_score_anchor
redis.replicate_commands()
local key = KEYS[1]
local key_rank = key .. ":rank"
local key_info = key .. ":info"
local arg_score = tonumber(ARGV[1])
local arg_high = tonumber(ARGV[2])
local arg_low = tonumber(ARGV[3])
local arg_includeInfo = tonumber(ARGV[4])

local results = {}

local dir = -1 --分数命中方向
local rs = redis.call('ZREVRANGEBYSCORE', key_rank, arg_score, "-inf", "WITHSCORES", "LIMIT", 0, 1)
if #rs == 0 then
    rs = redis.call('ZRANGEBYSCORE', key_rank, arg_score, "+inf", "WITHSCORES", "LIMIT", 0, 1)
    if #rs == 0 then
        return cjson.encode(results)
    end
    dir = 1
end
if tonumber(rs[2]) ~= arg_score then
    if dir < 0 then
        arg_high = arg_high - 1
    else
        arg_low = arg_low - 1
    end
end
local rank = redis.call('ZRANK', key_rank, rs[1])
local high = rank - arg_high
if high < 0 then
    high = 0
end
rs = redis.call('ZRANGE', key_rank, high, rank + arg_low, "WITHSCORES")
rank = high+1
for i = 1, #rs, 2 do
    local result = {
        key = string.sub(rs[i], 16),
        score = rs[i + 1],
        rank = rank
    }
    results[result.key] = result
    rank = rank + 1
end

if arg_includeInfo == 1 and next(results) ~= nil then
    local keys = {}
    for k, _ in pairs(results) do
        table.insert(keys, k)
    end
    local step = 5000
    local max = #keys
    for i = 1, max, step do
        local rs2 = redis.call('hmget', key_info, unpack(keys, i, math.min(i + step - 1, max)))
        local keysIdx = i
        for j = 1, #rs2 do
            if rs2[j] ~= false then
                results[keys[keysIdx]].info = rs2[j]
            end
            keysIdx = keysIdx + 1
        end
    end
end
return cjson.encode(results)
