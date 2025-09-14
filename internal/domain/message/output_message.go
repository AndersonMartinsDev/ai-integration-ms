package message

// InputMessage representa o formato da mensagem que o ai-integration-ms consome.
type OutPutMessage struct {
	FontNumber  string `json:"font_number"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}
