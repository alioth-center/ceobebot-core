package log

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

const timeFormat = "2006.01.02-15:04:05.000Z07:00"

type jsonFormatPayload struct {
	Time     string         `json:"time"`
	Level    string         `json:"level"`
	File     string         `json:"file"`
	Function string         `json:"function"`
	Payload  map[string]any `json:"data"`
}

type jsonFormatter struct {
	original *logrus.JSONFormatter
}

func (f jsonFormatter) callParent(entry *logrus.Entry) {
	_, _ = f.original.Format(entry)
}

func (f jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	go f.callParent(entry)

	data := map[string]any{}
	for key, value := range entry.Data {
		switch key {
		case "file_name", "func_name":
			continue
		default:
			data[key] = value
		}
	}
	payload := jsonFormatPayload{
		Time:     time.Now().Format(timeFormat),
		Level:    entry.Level.String(),
		File:     entry.Data["file_name"].(string),
		Function: entry.Data["func_name"].(string),
		Payload:  data,
	}

	result, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	} else {
		return append(result, '\n'), nil
	}
}
