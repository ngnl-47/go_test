--leaderboard_get_ranks_by_score_range
redis.replicate_commands()
local key = KEYS[1]
local key_rank = key .. ":rank"
local key_info = key .. ":info"
local arg_min = tonumber(ARGV[1])
local arg_max = tonumber(ARGV[2])
local arg_skip = tonumber(ARGV[3])
local arg_limit = tonumber(ARGV[4])
local arg_includeInfo = tonumber(ARGV[5])

local results = {}

local rs
if arg_limit > 0 and arg_skip >= 0 then
    rs = redis.call('ZRANGEBYSCORE', key_rank, arg_min, arg_max, "WITHSCORES", "limit", arg_skip, arg_limit)
else
    rs = redis.call('ZRANGEBYSCORE', key_rank, arg_min, arg_max, "WITHSCORES")
end

if #rs == 0 then
    return cjson.encode(results)
end
local rank = redis.call('ZRANK', key_rank, rs[1])
rank = rank + 1
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
