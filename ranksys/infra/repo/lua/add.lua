--leaderboard_add
redis.replicate_commands()
local key_rank = KEYS[1]
local key_map = KEYS[2]
local key_info = KEYS[3]
local arg_id = ARGV[1]
local arg_prefixedId = ARGV[2]
local arg_scoreAdd = tonumber(ARGV[3])
local arg_info = ARGV[4]
local arg_timeout = tonumber(ARGV[5])
local arg_slice = tonumber(ARGV[6])

local result = {}
local myRank = 0
local oldPrefixedId = redis.call('hget', key_map, arg_id)
local oldScore = 0
if oldPrefixedId ~= false then
    local oldScore2 = redis.call('zscore', key_rank, oldPrefixedId)
    if oldScore2 ~= false then
        oldScore = tonumber(oldScore2)
    end
end
local score = oldScore + arg_scoreAdd
if arg_info ~= '' then
    redis.call('hset', key_info, arg_id, arg_info)
end
if oldPrefixedId ~= false then
    redis.call('zrem', key_rank, oldPrefixedId)
end
redis.call('hset', key_map, arg_id, arg_prefixedId)
redis.call('zadd', key_rank, score, arg_prefixedId)
myRank = redis.call('zrank', key_rank, arg_prefixedId)
result.rank = myRank + 1
result.score = score
result.ok = true
if arg_timeout > 0 then
    redis.call('EXPIRE', key_rank, arg_timeout)
    redis.call('EXPIRE', key_map, arg_timeout)
    redis.call('EXPIRE', key_info, arg_timeout)
end
if string.sub(arg_prefixedId, 1, 1) == '-' then
    result.score = -result.score
end

if arg_slice > 0 then
    local rs = redis.call('zrange', key_rank, arg_slice,  -1)
    local max = #rs
    if max > 0 then
        redis.call('ZREMRANGEBYRANK', key_rank, arg_slice,  -1)
        local toDel = {}
        for _, v in ipairs(rs) do
            table.insert(toDel, string.sub(v, 16))
        end
        local step = 100
        for i = 1, max, step do
            redis.call('hdel', key_map, unpack(toDel, i, math.min(i + step - 1, max)))
            redis.call('hdel', key_info, unpack(toDel, i, math.min(i + step - 1, max)))
        end
    end
end

return cjson.encode(result)
