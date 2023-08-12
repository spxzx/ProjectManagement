package config

import "github.com/spxzx/project-common/kfk"

var kw *kfk.KafkaWriter

func InitKafkaWriter() func() {
	kw = kfk.GetWriter("localhost:9092")
	return kw.Close
}

func SendLog(data []byte) {
	kw.Send(kfk.LogData{
		Topic: "msproject_log",
		Data:  data,
	})
}
