-- leaderboard_getRanksByPlayerIds：按玩家 id 批量查名次与分。
redis.replicate_commands()
local key_rank = KEYS[1]
local key_map = KEYS[2]
local key_info = KEYS[3]

local arg_includeInfo = tonumber(ARGV[1])
local ids = cjson.decode(ARGV[2])

local results = {}
for _, pid in ipairs(ids) do
    local prefixed = redis.call('hget', key_map, pid)
    if prefixed ~= false then
        local r0 = redis.call('zrank', key_rank, prefixed)
        if r0 ~= false then
            local sc = redis.call('zscore', key_rank, prefixed)
            local disp = tonumber(sc)
            if string.sub(prefixed, 1, 1) == '-' then
                disp = -disp
            end
            results[pid] = {
                key = pid,
                score = disp,
                rank = r0 + 1
            }
        end
    end
end

if arg_includeInfo == 1 and next(results) ~= nil then
    local keys = {}
    for k, _ in pairs(results) do
        table.insert(keys, k)
    end
    table.sort(keys)
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
