package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ozeer/sloth/config"
	"github.com/ozeer/sloth/global"
	"github.com/ozeer/sloth/model/cache"
	"github.com/ozeer/sloth/third"
	"github.com/ozeer/sloth/tool"
)

var (
	// 每个定时器对应一个bucket
	timers []*time.Ticker
	// 存放bucket的channel
	bucket <-chan string
)

// AddTask 添加任务
func AddTask(job cache.Job) error {
	err := cache.AddJob(job.Id, job)
	if err != nil {
		global.Errorf("Add job(%s) to job pool fail：%s", job.Id, err.Error())
		return err
	}

	bucketName := <-bucket
	err = cache.PushToBucket(bucketName, job.Delay, job.Id)
	if err != nil {
		global.Errorf("Add job(%s) to bucket(%s) fail：%s", job.Id, bucketName, err.Error())
		return err
	}

	return nil
}

// Init 初始化定时任务
func InitTimerTask(cf config.Conf) {
	InitTimers(cf)
	bucket = GenerateBucketName(cf)
}

// GenerateBucketName 轮询获取bucket名称, 使任务分布到不同的bucket中, 提高扫描速度
func GenerateBucketName(cf config.Conf) <-chan string {
	bucketChan := make(chan string)
	go func() {
		i := 1
		for {
			bucketChan <- fmt.Sprintf(cf.Core.BucketName, i)
			if i >= cf.Core.BucketSize {
				i = 1
			} else {
				i++
			}
		}
	}()

	return bucketChan
}

// InitTimers 初始化定时器
func InitTimers(cf config.Conf) {
	timers = make([]*time.Ticker, cf.Core.BucketSize)
	var bucketName string
	for i := 0; i < cf.Core.BucketSize; i++ {
		timers[i] = time.NewTicker(1 * time.Second)
		bucketName = fmt.Sprintf(cf.Core.BucketName, i+1)
		go Timer(timers[i], bucketName)
	}
}

func Timer(timer *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-timer.C:
			ScanBucket(t, bucketName)
		}
	}
}

// ScanBucket 扫描bucket, 取出延迟时间小于当前时间的Job
func ScanBucket(t time.Time, bucketName string) {
	for {
		bucketItem, err := cache.GetFromBucket(bucketName)
		if err != nil {
			return
		}

		// 集合为空
		if bucketItem == nil {
			//global.InfoF("%s：no task!", bucketName)
			return
		}

		// 延迟时间未到
		if bucketItem.Timestamp > t.Unix() {
			return
		}

		job, err := cache.GetJob(bucketItem.JobId)
		if err != nil {
			cache.RemoveFromBucket(bucketName, bucketItem.JobId)
			continue
		}

		// job元信息不存在, 从bucket中删除
		if job == nil {
			cache.RemoveFromBucket(bucketName, bucketItem.JobId)
			continue
		}

		// 再次确认元信息中delay是否小于等于当前时间
		if job.Delay > t.Unix() {
			continue
		}

		// 延迟时间小于等于当前时间, 执行定时任务，从bucket中删除任务
		num := cache.RemoveFromBucket(bucketName, bucketItem.JobId)
		if num > 0 {
			var body map[string]string
			err = json.Unmarshal([]byte(job.Body), &body)
			if err != nil {
				global.Errorf("json decode fail：%s(jobId:%s)", err.Error(), job.Id)
				return
			}
			body["__uk"] = job.Id
			insertQueue(body["controller"], body["queueName"], body)

			// 发送消费微信通知
			//if WeChatSwitchOn {
			//	var atList,mentionedMobileList []string
			//	msgMap := []string{
			//		"任务主题：" + job.Topic,
			//		"任务id：" + job.Id,
			//		"任务参数：" + job.Body,
			//		"消费时间：" + tool.CurrentDate(),
			//	}
			//	msg := strings.Join(msgMap, "\n")
			//	SendMsg(AlertMsgKey, msg, atList, mentionedMobileList)
			//}

			// 发送消费钉钉机器人通知
			// var atUserIds, atMobiles []string
			// msgMap := []string{
			// 	"任务主题：" + job.Topic,
			// 	"任务id：" + job.Id,
			// 	"任务参数：" + job.Body,
			// 	"消费时间：" + tool.CurrentDate(),
			// }
			// msg := strings.Join(msgMap, "\n")
			// third.SendTextMsg(third.AccessToken, msg, true, atUserIds, atMobiles)

			msg := fmt.Sprintf("任务主题: %s\\n任务id: %s\\n消费时间: %s", job.Topic, job.Id, tool.CurrentDate())
			third.SendLarkMsg(msg)
		}
	}
}
