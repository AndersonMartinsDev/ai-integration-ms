package gateway

import "ai-integration-ms/internal/domain/model"

type WebHookProcessorClient interface {
	GetSessionData(session_id string) *model.SessionModel
}
