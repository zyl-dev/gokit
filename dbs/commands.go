package dbs

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// RedisDBHSet 这个方法的调用，通常建议写多个 Redis，以防不同的客户端请求，
// 在写入前需要通过 InitRedisAllDB 进行初始化，这样的例子大多可以在 Job 中看到
func RedisDBHSet(key, field string, value interface{}) *redis.IntCmd {
	var result *redis.IntCmd
	if RedisClusterDB != nil {
		clusterResult := RedisClusterDB.HSet(context.TODO(), key, field, value)
		if clusterResult.Err() == nil {
			result = clusterResult
		}
	}
	if RedisRingDB != nil {
		ringResult := RedisRingDB.HSet(context.TODO(), key, field, value)
		if ringResult.Err() == nil {
			result = ringResult
		}
	}
	if RedisDB != nil {
		clientResult := RedisDB.HSet(context.TODO(), key, field, value)
		if clientResult.Err() == nil {
			result = clientResult
		}
	}
	return result
}

// RedisDBHGet 对 redis 的 hget 操作进行封装，避免直接调用
func RedisDBHGet(key, field string) *redis.StringCmd {
	if CurCacheRedisType == RedisClusterType {
		return RedisClusterDB.HGet(context.TODO(), key, field)
	}
	if CurCacheRedisType == RedisRingType {
		return RedisRingDB.HGet(context.TODO(), key, field)
	}
	return RedisDB.HGet(context.TODO(), key, field)
}

func RedisDBHMSet(key string, fields map[string]interface{}) *redis.BoolCmd {
	if CacheRedisType == RedisClusterType {
		return RedisClusterDB.HMSet(context.TODO(), key, fields)
	}
	if CacheRedisType == RedisRingType {
		return RedisRingDB.HMSet(context.TODO(), key, fields)
	}
	return RedisDB.HMSet(context.TODO(), key, fields)
}

func RedisDBHMGet(key string, fields ...string) *redis.SliceCmd {
	if CacheRedisType == RedisClusterType {
		return RedisClusterDB.HMGet(context.TODO(), key, fields...)
	}
	if CacheRedisType == RedisRingType {
		return RedisRingDB.HMGet(context.TODO(), key, fields...)
	}
	return RedisDB.HMGet(context.TODO(), key, fields...)
}
