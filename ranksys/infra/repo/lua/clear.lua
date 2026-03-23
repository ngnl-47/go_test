--leaderboard_clear
redis.replicate_commands()
redis.call('del', KEYS[1])
redis.call('del', KEYS[2])
redis.call('del', KEYS[3])
return 'ok'
