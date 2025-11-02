package services

import (
	"ai-integration-ms/internal/domain/model"
	"context"
)

type FileServiceInterface interface {
	DownloadAudioFile(ctx context.Context, aiRequestModel model.AiRequestModel) ([]byte, error)
}
