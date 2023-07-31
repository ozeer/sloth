package model

import (
	"github.com/go-redis/redis/v8"
	"github.com/ozeer/sloth/config"
	"github.com/ozeer/sloth/global"
)

// BucketItem bucket中的元素
type BucketItem struct {
	Timestamp int64
	JobId     string
}

// Bucket存储结构：Zset；key为dq_bucket_{num}，score为延迟任务执行时的时间戳，member为job_id

// PushToBucket 添加JobId到bucket中
func PushToBucket(key string, timestamp int64, jobId string) error {
	_, err := config.Rdb.Do(config.Ctx, "ZADD", key, timestamp, jobId).Result()
	return err
}

// GetFromBucket 从bucket中获取延迟时间最小的JobId
func GetFromBucket(key string) (*BucketItem, error) {
	value, err := config.Rdb.ZRangeByScoreWithScores(config.Ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, nil
	}

	if len(value) == 0 {
		return nil, nil
	}

	item := &BucketItem{}
	item.Timestamp = int64(value[0].Score)
	item.JobId = value[0].Member.(string)

	return item, nil
}

// RemoveFromBucket 从bucket中删除JobId
func RemoveFromBucket(bucket string, jobId string) int64 {
	num, err := config.Rdb.ZRem(config.Ctx, bucket, jobId).Result()
	if err != nil {
		global.Errorf("delete job(%s) fail from bucket(%s)：%s", jobId, bucket, err.Error())
	}
	return num
}
