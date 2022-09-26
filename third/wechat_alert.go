package third

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseUrl     = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	AlertMsgKey = "3bc2c0f7-d1c6-4c31-9afe-42bd4e103bef"
	// WeChatSwitchOn 企业微信发送开关
	WeChatSwitchOn = true
)

// TextParams text类型消息数据
type TextParams struct {
	MsgType string   `json:"msgtype"`
	Text    TextBody `json:"text"`
}
type TextBody struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

// MarkdownParams markdown类型消息数据
type MarkdownParams struct {
	MsgType  string       `json:"msgtype"`
	MarkDown MarkdownBody `json:"markdown"`
}
type MarkdownBody struct {
	Content string `json:"content"`
}

// SendMsg 发送POST请求，消息格式：text
// key:			 微信机器人key
// msg：         消息内容
// atList：      希望@的人
// mentionedMobileList： 希望通过手机号@的人
func SendMsg(key string, msg string, atList []string, mentionedMobileList []string) string {
	url := baseUrl + key
	client := &http.Client{Timeout: 2 * time.Second}

	textBody := TextBody{
		msg, atList, mentionedMobileList,
	}
	message := TextParams{"text", textBody}
	jsonStr, _ := json.Marshal(message)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

// SendMarkdownMsg 发送POST请求，消息格式：markdown
// key:			 微信机器人key
// msg：         消息内容
// atList：      希望@的人
// mentionedMobileList： 希望通过手机号@的人
func SendMarkdownMsg(key string, msg string) string {
	url := baseUrl + key
	client := &http.Client{Timeout: 2 * time.Second}

	markdownBody := MarkdownBody{
		msg,
	}
	message := MarkdownParams{"markdown", markdownBody}
	jsonStr, _ := json.Marshal(message)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
