--leaderboard_get_ranks_by_ids
redis.replicate_commands()
local key = KEYS[1]
local key_rank = key .. ":rank"
local key_map = key .. ":map"
local key_info = key .. ":info"
local arg_ids = cjson.decode(ARGV[1])
local arg_includeInfo = ARGV[2]

local results = {}

local step = 5000
local max = #arg_ids
local sort = 0
for i = 1, max, step do
    local rs = redis.call('hmget', key_map, unpack(arg_ids, i, math.min(i + step - 1, max)))
    for j, prefixedId in ipairs(rs) do
        if prefixedId == false then
            table.insert(results, cjson.null)
        else
            local rank = redis.call('zrank', key_rank, prefixedId)
            if rank == false then
                table.insert(results, cjson.null)
            else
                local score = redis.call('ZSCORE', key_rank, prefixedId)
                if sort == 0 then
                    if string.sub(prefixedId, 1, 1) == '-' then
                        sort = -1
                    else
                        sort = 1
                    end
                end
                score = score * sort
                table.insert(results, {
                    roleIdStr = arg_ids[i + j - 1],
                    score = score,
                    rank = rank + 1
                })
            end
        end
    end
    if arg_includeInfo == "1" then
        local rs2 = redis.call('hmget', key_info, unpack(arg_ids, i, math.min(i + step - 1, max)))
        for j, info in ipairs(rs2) do
            local result = results[i + j - 1]
            if result ~= cjson.null then
                result.info = info
            end
        end
    end
end
return cjson.encode(results)
