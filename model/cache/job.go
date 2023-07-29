package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/ozeer/sloth/model/storage"
	"github.com/vmihailenco/msgpack"
)

// Job 使用msgpack序列化后保存到Redis,减少内存占用
// Job存储结构：String；key为Job id，value为job元信息
type Job struct {
	// Job类型
	Topic string `json:"topic" msgpack:"1"`
	// job唯一标识ID，需确保Job ID唯一
	Id string `json:"id" msgpack:"2"`
	// Job需要延迟的时间, unix时间戳，单位：秒
	Delay int64 `json:"delay" msgpack:"3"`
	// Job的内容，供消费者做具体的业务处理，如果是json格式需转义
	Body string `json:"body" msgpack:"4"`
}

// GetJob 获取Job
func GetJob(key string) (*Job, error) {
	value, err := storage.Rdb.Get(storage.Ctx, key).Result()
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		return nil, err
	} else {
		byteValue := []byte(value)
		job := &Job{}
		err = msgpack.Unmarshal(byteValue, job)
		if err != nil {
			return nil, err
		}

		return job, nil
	}
}

// AddJob 添加Job
func AddJob(key string, job Job) error {
	value, err := msgpack.Marshal(job)
	if err != nil {
		return err
	}
	_, err = storage.Rdb.Set(storage.Ctx, key, value, 0).Result()
	return err
}

// RemoveJob 删除Job
func RemoveJob(key string) error {
	_, err := storage.Rdb.Del(storage.Ctx, key).Result()
	return err
}
