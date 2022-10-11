package third

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sloth/config"
	"time"
)

// 钉钉机器人发送开关
var (
	SwitchOn    bool
	BaseUrl     string
	AccessToken string
)

// TextMsgParams text类型消息数据
type TextMsgParams struct {
	MsgType string      `json:"msgtype"`
	Text    TextMsgBody `json:"text"`
	At      AtBody      `json:"at"`
}

type TextMsgBody struct {
	Content string `json:"content"`
}

type AtBody struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `bool:"isAtAll"`
}

// MarkdownMsgParams markdown类型消息数据
type MarkdownMsgParams struct {
	MsgType  string          `json:"msgtype"`
	MarkDown MarkdownMsgBody `json:"markdown"`
	At       AtBody          `json:"at"`
}
type MarkdownMsgBody struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func InitDingTalk(c config.Conf) error {
	BaseUrl = c.DingTalk.BaseUrl
	AccessToken = c.DingTalk.AccessToken
	SwitchOn = c.DingTalk.SwitchOn

	if BaseUrl == "" || AccessToken == "" {
		return errors.New("dingTalk configuration parameter error")
	}

	return nil
}

// SendTextMsg 发送POST请求，消息格式：text
// accessToken:	     机器人access_token
// msg：         消息内容
// isAtAll:      是否@所有人
// atUserIds：      希望通过user id@的人
// atMobiles：   希望通过手机号@的人
func SendTextMsg(accessToken string, msg string, isAtAll bool, atUserIds []string, atMobiles []string) string {
	if !SwitchOn {
		return ""
	}

	url := BaseUrl + accessToken
	client := &http.Client{Timeout: 2 * time.Second}

	textBody := TextMsgBody{
		msg,
	}
	atBody := AtBody{
		atMobiles, atUserIds, isAtAll,
	}
	message := TextMsgParams{"text", textBody, atBody}
	jsonStr, _ := json.Marshal(message)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result)
}

// SendMarkdownMessage 发送POST请求，消息格式：markdown
// accessToken:	     机器人access_token
// title：       消息标题
// text:         消息正文
// isAtAll:      是否@所有人
// atUserIds：   希望通过user id@的人
// atMobiles：   希望通过手机号@的人
func SendMarkdownMessage(accessToken string, title string, text string, isAtAll bool, atUserIds []string, atMobiles []string) string {
	if !SwitchOn {
		return ""
	}

	url := BaseUrl + accessToken
	client := &http.Client{Timeout: 2 * time.Second}

	markdownBody := MarkdownMsgBody{
		title, text,
	}
	atBody := AtBody{
		atMobiles, atUserIds, isAtAll,
	}
	message := MarkdownMsgParams{"markdown", markdownBody, atBody}
	jsonStr, _ := json.Marshal(message)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result)
}
