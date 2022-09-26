package service

import (
	"testing"
)

func TestInsert(t *testing.T) {
	data := make(map[string]string)
	insertQueue("queue/delayqueue/test", QueueName, data)
}
