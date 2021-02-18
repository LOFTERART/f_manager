/**
 * @Author: LOFTER
 * @Description:
 * @File:  cache
 * @Date: 2021/2/8 12:47 上午
 */
package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisOp struct {
	redisPool *redis.Pool
}

/**
检查是否已满
*/
func (rop *RedisOp) ChunkIsFull(fileMd5 string) int {
	cnt, e2 := redis.Int(rop.redisPool.Get().Do("ZCARD", fileMd5))
	errorHanler(e2)
	return cnt
}
func (rop *RedisOp) GetMem(fileMd5 string) []string {
	res, e := redis.Strings(rop.redisPool.Get().Do("ZRANGE", fileMd5, 0, -1)) // 无需分数 （顺序）
	errorHanler(e)
	return res
}
func (rop *RedisOp) GetMemWithScore(fileMd5 string) map[string]string {
	res, e := redis.StringMap(rop.redisPool.Get().Do("ZRANGE", fileMd5, 0, -1, "WITHSCORES")) // 无需分数 （顺序）
	errorHanler(e)
	return res
}
func (rop *RedisOp) ClearSet(fileMd5 string) {
	_, e := rop.redisPool.Get().Do("DEL", fileMd5)
	errorHanler(e)
}

/**
添加块
*/
func (rop *RedisOp) ChunkAdd(idx int, fileMd5 string, pieceMd5 string) {
	rop.redisPool.Get().Do("ZADD", fileMd5, idx, pieceMd5)
}

/**
文件信息存储、查询
*/
func (rop *RedisOp) FileInfo(fileMd5 string, filename string, filepath string, filetype string) string {
	res, e := redis.String(rop.redisPool.Get().Do("HMSET", "hashmap_"+fileMd5,
		"filename", filename,
		"filepath", filepath,
		"filetype", filetype,
	))
	errorHanler(e)
	return res
}

func (rop *RedisOp) GetFileinfo(fileMd5 string) map[string]string {
	res, e := redis.StringMap(rop.redisPool.Get().Do("HGETALL", "hashmap_"+fileMd5))
	errorHanler(e)
	return res
}

//删除hash

func (rop *RedisOp) DelFileInfo(fileMd5 string) int {
	res, e := redis.Int(rop.redisPool.Get().Do("HDEL", "hashmap_"+fileMd5, "filename", "filepath", "filetype"))
	errorHanler(e)
	return res
}

/**
合并加锁
*/
func (rop *RedisOp) Merging(fileMd5 string) int {
	res, e := redis.Int(rop.redisPool.Get().Do("SADD", "merging_list", fileMd5))
	errorHanler(e)
	return res
}

/**
是否正在合并
返回int
*/
func (rop *RedisOp) ISMerging(fileMd5 string) int {
	res, e := redis.Int(rop.redisPool.Get().Do("SISMEMBER", "merging_list", fileMd5))
	errorHanler(e)
	return res
}

/**
合并完删除
*/
func (rop *RedisOp) DelMerging(fileMd5 string) int {
	res, e := redis.Int(rop.redisPool.Get().Do("SREM", "merging_list", fileMd5))
	errorHanler(e)
	return res
}

//newPool
func (rop *RedisOp) newPool() *redis.Pool {
	server := ":6379"
	return &redis.Pool{
		MaxIdle:     10000,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
func (rop *RedisOp) RedisInit() *RedisOp {
	rop.redisPool = rop.newPool()
	return rop
}

// 全局
func errorHanler(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
