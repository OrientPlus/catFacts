package entity

import "encoding/json"

type Fact struct {
	Message string `json:"fact"`
	Length  int32  `json:"length"`
}

type Request struct {
	CorrelationID string          `json:"correlation_id"`
	ReplyTopic    string          `json:"reply_topic"`
	Payload       json.RawMessage `json:"payload"`
}
