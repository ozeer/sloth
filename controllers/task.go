package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ozeer/sloth/global"
	"github.com/ozeer/sloth/model/cache"
	"github.com/ozeer/sloth/service"
	"github.com/ozeer/sloth/third"
	"github.com/ozeer/sloth/tool"
)

// AddTask 添加定时任务
func AddTask(c *gin.Context) {
	var job cache.Job

	id := c.PostForm("id")
	topic := c.PostForm("topic")
	delay := tool.StringToInt64(c.PostForm("delay"))
	body := c.PostForm("body")

	// 校验参数
	if id == "" {
		tool.Fail(c, "id不能为空")
		return
	}

	// 校验参数
	if topic == "" {
		tool.Fail(c, "topic不能为空")
		return
	}

	if delay < 0 {
		tool.Fail(c, "delay参数不合法")
		return
	}

	job.Id = id
	job.Topic = topic
	job.Body = body
	job.Delay = delay

	var temp = job

	job.Delay = time.Now().Unix() + job.Delay
	err := service.AddTask(job)

	if err != nil {
		global.Error("Add task fail：", err.Error())
		tool.Fail(c, err.Error())
		return
	} else {
		global.InfoF("Insert Queue:%v", body)

		// 创建任务消息通知
		//if service.WeChatSwitchOn {
		//	var atList,mentionedMobileList []string
		//	msgMap := []string{
		//		"任务主题：" + job.Topic,
		//		"任务id：" + job.Id,
		//		"延迟时间(秒)：" + tool.Int64ToString(delay),
		//		"任务参数：" + job.Body,
		//		"创建时间：" + tool.CurrentDate(),
		//		"执行时间：" + tool.TimestampToDate(job.Delay),
		//	}
		//	msg := strings.Join(msgMap, "\n")
		//	service.SendMsg(service.AlertMsgKey, msg, atList, mentionedMobileList)
		//}

		// DingTalk消息
		// var atUserIds, atMobiles []string
		// msgMap := []string{
		// 	"任务主题：" + job.Topic,
		// 	"任务id：" + job.Id,
		// 	"任务参数：" + job.Body,
		// 	"延迟时间(秒)：" + tool.Int64ToString(delay),
		// 	"创建时间：" + tool.CurrentDate(),
		// 	"执行时间：" + tool.TimestampToDate(job.Delay),
		// }
		// msg := strings.Join(msgMap, "\n")
		// third.SendTextMsg(third.AccessToken, msg, true, atUserIds, atMobiles)

		// 飞书消息
		msg := fmt.Sprintf("任务主题: %s\\n任务id: %s\\n延迟时间: %s(s)\\n创建时间: %s\\n执行时间: %s", job.Topic, job.Id, tool.Int64ToString(delay), tool.CurrentDate(), tool.TimestampToDate(job.Delay))
		third.SendLarkMsg(msg)

		tool.Success(c, "添加成功", temp)
		return
	}
}

// GetTask 查询定时任务详情
func GetTask(c *gin.Context) {
	id := c.Query("id")
	job, err := cache.GetJob(id)
	if err != nil {
		global.Errorf("Get task#%s# fail：%s", id, err.Error())
		tool.Fail(c, err.Error())
	}

	if job == nil {
		global.Error("Task not exist：", id)
		tool.Fail(c, "Task not exist")
	}

	var body interface{}
	err = json.Unmarshal([]byte(job.Body), &body)
	if err != nil {
		tool.Fail(c, err.Error())
	}
	tool.Success(c, "获取成功", job)
}

// DelTask 删除任务
func DelTask(c *gin.Context) {
	id := c.PostForm("id")
	err := cache.RemoveJob(id)

	if err != nil {
		tool.Fail(c, "delete task fail")
	} else {
		var data interface{}
		tool.Success(c, "取消成功", data)
	}
}

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to sloth！",
	})
}
