package model

import "encoding/json"

type RedisSessionModel struct {
	KeyId   string
	AgentId uint
	History json.RawMessage
}
