-- leaderboard_getRanksByScore：ZRANGEBYSCORE，返回展示分与 1-based 名次。
redis.replicate_commands()
local key_rank = KEYS[1]
local key_map = KEYS[2]
local key_info = KEYS[3]

local minS = tonumber(ARGV[1])
local maxS = tonumber(ARGV[2])
local arg_includeInfo = tonumber(ARGV[3])

local rs = redis.call('zrangbyscore', key_rank, minS, maxS, 'WITHSCORES')
local results = {}

for i = 1, #rs, 2 do
    local member = rs[i]
    local sc = tonumber(rs[i + 1])
    local r0 = redis.call('zrank', key_rank, member)
    if r0 ~= false then
        local rank = r0 + 1
        local key = string.sub(member, 16)
        local disp = sc
        if string.sub(member, 1, 1) == '-' then
            disp = -disp
        end
        results[key] = {
            key = key,
            score = disp,
            rank = rank
        }
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
