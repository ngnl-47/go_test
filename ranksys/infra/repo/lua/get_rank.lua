--leaderboard_getRank
redis.replicate_commands()
local arg_min = tonumber(ARGV[1])
local arg_max = tonumber(ARGV[2])
local arg_myId = ARGV[3]
local arg_includeInfo = tonumber(ARGV[4])

local results = {}

local rs = redis.call('zrange', KEYS[1], arg_min - 1, arg_max - 1, 'WITHSCORES')
local rank = arg_min
for i = 1, #rs, 2 do
    local result = {
        key = string.sub(rs[i], 16),
        score = rs[i + 1],
        rank = rank
    }
    if string.sub(rs[i], 1, 1) == '-' then
        result.score = -result.score
    end
    results[result.key] = result
    rank = rank + 1
end
if arg_myId ~= '' and results[arg_myId] == nil then
    local myPrefixedId = redis.call('hget', KEYS[2], arg_myId)
    if myPrefixedId ~= false then
        local myRank = redis.call('zrank', KEYS[1], myPrefixedId)
        if myRank ~= false then
            local myScore = redis.call('ZSCORE', KEYS[1], myPrefixedId)
            local me = {
                key = arg_myId,
                score = myScore,
                rank = myRank + 1
            }
            if string.sub(myPrefixedId, 1, 1) == '-' then
                me.score = -me.score
            end
            results[me.key] = me
        end
    end
end
if arg_includeInfo == 1 and next(results) ~= nil then
    local keys = {}
    for k, _ in pairs(results) do
        table.insert(keys, k)
    end
    local step = 5000
    local max = #keys
    for i = 1, max, step do
        local rs2 = redis.call('hmget', KEYS[3], unpack(keys, i, math.min(i + step - 1, max)))
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
