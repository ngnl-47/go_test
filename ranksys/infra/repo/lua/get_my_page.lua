--leaderboard_get_my_page
redis.replicate_commands()
local key = KEYS[1]
local key_rank = key .. ":rank"
local key_map = key .. ":map"
local key_info = key .. ":info"
local arg_id = ARGV[1]
local arg_pageSize = tonumber(ARGV[2])
local arg_includeInfo = tonumber(ARGV[3])

local prefixedId = redis.call('hget', key_map, arg_id)
if prefixedId == false then
    return "[]"
end

local rank = redis.call('zrank', key_rank, prefixedId)
if rank == false then
    return "[]"
end

local results = {}

local min = math.floor(rank / arg_pageSize) * arg_pageSize
local max = min + arg_pageSize - 1

local rank = min
local rs = redis.call('zrange', key_rank, min, max, 'WITHSCORES')
for i = 1, #rs, 2 do
    rank = rank+1
    local result = {
        key = string.sub(rs[i], 16),
        score = rs[i + 1],
        rank = rank
    }
    if string.sub(rs[i], 1, 1) == '-' then
        result.score = -result.score
    end
    results[result.key] = result
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
