package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		delay := rand.Intn(99) + 1
		data := `{
			"topic": "order",
			"delay": %d,
			"body": "{"controller":"queue/delayqueue/test","queueName":"hzQDoraemon""
		}`
		data = fmt.Sprintf(data, delay)
		request, _ := http.NewRequest("POST", "http://localhost:9288/v1/add_task", strings.NewReader(data))
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("error：%v\n", err)
		} else {
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("success, delay %d, response data：%v\n", delay, string(respBody))
		}
	}
}
