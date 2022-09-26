package third

import (
	"log"
	"testing"
)

// go test -v wechat_alert_test.go wechat_alert.go

func TestSendMsg(t *testing.T) {
	var atList, mentionedMobileList []string
	atList = append(atList, "zhouyang")
	res := SendMsg(AlertMsgKey, "测试", atList, mentionedMobileList)
	log.Println(res)
}
