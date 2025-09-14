package model

import "encoding/json"

type SessionModel struct {
	KeyId       string
	AgentId     uint
	PhoneNumber string
	History     json.RawMessage
}
