--leaderboard_remove
redis.replicate_commands()
local arg_id = ARGV[1]

local oldPrefixedId = redis.call('hget', KEYS[2], arg_id)
if oldPrefixedId ~= false then
    redis.call('zrem', KEYS[1], oldPrefixedId)
end
redis.call('hdel', KEYS[2], arg_id)
redis.call('hdel', KEYS[3], arg_id)
return 'ok'
