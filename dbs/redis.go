package dbs

import (
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// RedisDB redis client
var RedisDB *redis.Client
var RedisClusterDB *redis.ClusterClient
var RedisRingDB *redis.Ring
var CacheRedisType string
var CurCacheRedisType string // 应该智能判断当前模式，以防 Redis 节点差错

const RedisClusterType = "cluster"
const RedisRingType = "ring"
const RedisNormalType = "normal"

type Rediser interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
}

func GetRedisCacheDB() Rediser {
	if CacheRedisType == RedisClusterType {
		return RedisClusterDB
	}
	if CurCacheRedisType == RedisClusterType {
		return RedisClusterDB
	}
	if CurCacheRedisType == RedisRingType {
		return RedisRingDB
	}
	return RedisDB
}

func InitRedisAllDB(c *RedisConfig) {
	if err := InitRedisCacheDBCluster(c.Cluster.Addrs, c.Cluster.Password); err != nil {
		log.Error(err.Error())
	}
	if err := InitRedisCacheRingDB(c.Ring.Addrs, c.Ring.Password); err != nil {
		log.Error(err.Error())
	}
	if err := InitRedisCacheDBClient(c.Redis.Address, c.Redis.Password); err != nil {
		log.Error(err.Error())
	}
}

func InitRedisCacheDB(c *RedisConfig) error {
	var err error
	CacheRedisType = c.Type
	// 先判断 Cluster 模式
	if CacheRedisType == RedisClusterType || len(c.Cluster.Addrs) != 0 {
		err = InitRedisCacheDBCluster(c.Cluster.Addrs, c.Cluster.Password)
		CurCacheRedisType = RedisClusterType
	}
	if err != nil || CacheRedisType == RedisRingType {
		err = InitRedisCacheRingDB(c.Ring.Addrs, c.Ring.Password)
		CurCacheRedisType = RedisRingType
	}
	if err != nil || (CacheRedisType != RedisClusterType && CacheRedisType != RedisRingType) {
		err = InitRedisCacheDBClient(c.Redis.Address, c.Redis.Password)
		CurCacheRedisType = RedisNormalType
	}
	log.WithField("CurType", CurCacheRedisType).WithField("PrevType", c.Type).WithError(err).Info("Init Redis Completed")
	return err
}

func InitRedisDB(c *RedisConfig) error {
	CacheRedisType = c.Type
	var err error
	if CacheRedisType == RedisClusterType {
		err = InitRedisCacheDBCluster(c.Cluster.Addrs, c.Cluster.Password)
	} else {
		err = InitRedisCacheDBClient(c.Redis.Address, c.Redis.Password)
	}

	return err
}

// InitRedisCacheDBClient 初始化 redis 数据库
func InitRedisCacheDBClient(address, password string) error {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:            address,
		Password:        password,
		DialTimeout:     3 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     5 * time.Second,
		MaxRetries:      3,
		MaxRetryBackoff: 3 * time.Second,
	})
	status := RedisDB.Ping()
	if status.Err() != nil {
		log.Error(status.Err())
		return status.Err()
	}
	return nil
}

func InitRedisCacheRingDB(addresses map[string]string, password string) error {
	RedisRingDB = redis.NewRing(&redis.RingOptions{
		Addrs:           addresses,
		Password:        password,
		DialTimeout:     3 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     5 * time.Second,
		MaxRetries:      3,
		MaxRetryBackoff: 3 * time.Second,
	})
	status := RedisRingDB.Ping()
	if status.Err() != nil {
		log.Error(status.Err())
		return status.Err()
	}
	return nil
}

func InitRedisCacheDBCluster(addresses []string, password string) error {
	RedisClusterDB = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           addresses,
		Password:        password,
		DialTimeout:     3 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     5 * time.Second,
		MaxRetries:      3,
		MaxRetryBackoff: 3 * time.Second,
		RouteRandomly:   true,
		ReadOnly:        true,
	})
	status := RedisClusterDB.Ping()
	if status.Err() != nil {
		log.Error(status.Err())
		return status.Err()
	}
	return nil
}
