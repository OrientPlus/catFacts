package entity

import (
	"encoding/json"
	"time"
)

type Fact struct {
	Message   string    `json:"fact"`
	Length    int32     `json:"length"`
	TimePoint time.Time `json:"time_point"`
}

type Request struct {
	CorrelationID string          `json:"correlation_id"`
	ReplyTopic    string          `json:"reply_topic"`
	Payload       json.RawMessage `json:"payload"`
}
