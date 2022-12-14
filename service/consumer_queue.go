package service

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"os"
	"sloth/config"
	"sloth/tool"
	"strconv"
	"time"
)

var (
	QueueName string
	rdb       *redis.Client
)

func InitConsumer(c config.Conf) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     c.ConsumerQueue.Host + ":" + c.ConsumerQueue.Port,
		Password: c.ConsumerQueue.Password,
	})
	QueueName = c.ConsumerQueue.QueueName
	_, err := rdb.Ping(rdb.Context()).Result()

	return err
}

func insertQueue(controller, queueName string, params map[string]string) {
	if len(queueName) == 0 {
		queueName = QueueName
	}

	delete(params, "controller")
	// 向消费业务队列写入数据
	params["__q"] = controller
	params["__insert_time"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["__from"], _ = os.Hostname()

	value, err := json.Marshal(params)
	if err != nil {
		return
	}

	// 写入消费队列
	result := rdb.RPush(rdb.Context(), QueueName, value)

	// 记录传入参数、写入数据及返回结果日志
	tool.LogAccess.Infof("Insert Queue:%v#Result:%v", params, result)

	return
}
