package kfk

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProduct(t *testing.T) {
	w := GetWriter("localhost:9092")
	m := make(map[string]string)
	m["projectCode"] = "114514"
	b, _ := json.Marshal(m)
	w.Send(LogData{
		Topic: "msproject_log",
		Data:  b,
	})
	time.Sleep(3 * time.Second)
}

func TestConsumer(t *testing.T) {
	GetReader([]string{"localhost:9092"}, "group1", "msproject_log")
	for {
	}
}
