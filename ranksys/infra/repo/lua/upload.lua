--leaderboard_upload
redis.replicate_commands()
local key_rank = KEYS[1]
local key_map = KEYS[2]
local key_info = KEYS[3]
local arg_id = ARGV[1]
local arg_prefixedId = ARGV[2]
local arg_score = tonumber(ARGV[3])
local arg_info = ARGV[4]
local arg_timeout = tonumber(ARGV[5])
local arg_forceReplace = ARGV[6]
local arg_slice = tonumber(ARGV[7])
local arg_autoIgnoreInfo = ARGV[8]

local result = {}
local myRank = 0
local ignore = false
local oldPrefixedId = redis.call('hget', key_map, arg_id)
local oldScore
if oldPrefixedId ~= false then
    oldScore = redis.call('zscore', key_rank, oldPrefixedId)
    if oldScore ~= false then
        oldScore = tonumber(oldScore)
        ignore = oldScore == arg_score or arg_forceReplace ~= '1' and arg_score >= tonumber(oldScore)
    end
end

if ignore then
    local forceUpdateInfo = arg_autoIgnoreInfo ~= "1"
    if (forceUpdateInfo and arg_info ~= '') then
        redis.call('hset', key_info, arg_id, arg_info)
    end
    myRank = redis.call('zrank', key_rank, oldPrefixedId)
    result.rank = myRank + 1
    result.score = oldScore
    result.ok = false
else
    if arg_info ~= '' then
        redis.call('hset', key_info, arg_id, arg_info)
    end
    if oldPrefixedId ~= false then
        redis.call('zrem', key_rank, oldPrefixedId)
    end
    redis.call('hset', key_map, arg_id, arg_prefixedId)
    redis.call('zadd', key_rank, arg_score, arg_prefixedId)
    myRank = redis.call('zrank', key_rank, arg_prefixedId)
    result.rank = myRank + 1
    result.score = arg_score
    result.ok = true
end
if arg_timeout > 0 then
    redis.call('EXPIRE', key_rank, arg_timeout)
    redis.call('EXPIRE', key_map, arg_timeout)
    redis.call('EXPIRE', key_info, arg_timeout)
end
if string.sub(arg_prefixedId, 1, 1) == '-' then
    result.score = -result.score
end

if arg_slice > 0 then
    local rs = redis.call('zrange', key_rank, arg_slice, -1)
    local max = #rs
    if max > 0 then
        redis.call('ZREMRANGEBYRANK', key_rank, arg_slice, -1)
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
