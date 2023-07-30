package third

import (
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	WebHook = "https://open.feishu.cn/open-apis/bot/v2/hook/b1e8ca8e-154a-49c7-8efe-2e25360da46a"
)

func SendLarkMsg(msg string) string {
	client := &http.Client{Timeout: 5 * time.Second}

	sendData := `{
		"msg_type": "text",
		"content": {"text": "` + "sloth消息通知: \\n" + msg + `"}
	  }`

	resp, err := client.Post(WebHook, "application/json", strings.NewReader(sendData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result)
}
