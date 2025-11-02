package message

// InputMessage representa o formato da mensagem que o ai-integration-ms consome.
type InputMessage struct {
	OriginMessage string `json:"origin_message"`
	SessionKey    string `json:"session_key"`
	AgentID       uint   `json:"agent_id"`
	Message       string `json:"message"`
	FontNumber    string `json:"font_number"`
	MessageType   string `json:"message_type"`
	MessageID     string `json:"message_id,omitempty"`
	MediaUrl      string `json:"media_url,omitempty"`
}
