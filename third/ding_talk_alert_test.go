package third

import (
	"log"
	"testing"
)

// go test -v ding_talk_alert_test.go ding_talk_alert.go

const AlertToken = ""

func TestSendTextMsg(t *testing.T) {
	var atUsers, atMobiles []string
	res := SendTextMsg(AlertToken, "测试2021", true, atUsers, atMobiles)
	log.Println(res)
}

func TestSendMarkdownMsg(t *testing.T) {
	var atUsers, atMobiles []string
	res := SendMarkdownMessage(AlertToken, "标题", "内容", true, atUsers, atMobiles)
	log.Println(res)
}
