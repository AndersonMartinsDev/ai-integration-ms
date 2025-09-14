package repository

import (
	"ai-integration-ms/internal/domain/model"
	"ai-integration-ms/internal/infrastructure/database"
	"context"
	"encoding/json"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	Get(keyId string) (*model.SessionModel, error)
	Save(session *model.SessionModel) (string, error)
	Delete(keyId string) error
}

type RedisRepositoryImpl struct {
}

func NewRedisRepository() *RedisRepositoryImpl {
	return &RedisRepositoryImpl{}
}

// Get busca uma sessão no Redis pelo UUID.
func (r RedisRepositoryImpl) Get(keyId string) (*model.SessionModel, error) {
	cache := database.CACHE
	key := keyId
	ctx := context.Background()
	val, err := cache.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar sessão no Redis: %w", err)
	}

	var session *model.SessionModel
	// json.Unmarshal diretamente para AiSession com json.RawMessage no History
	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar sessão do Redis: %w", err)
	}

	return session, nil
}

// Save salva uma sessão no Redis com um tempo de expiração (TTL).
// Este método substitui a lógica do seu antigo 'InsertSession' para o Redis.
func (r RedisRepositoryImpl) Save(session *model.SessionModel) (string, error) {
	cache := database.CACHE
	key := session.KeyId

	data, err := json.Marshal(session)
	if err != nil {
		return session.KeyId, fmt.Errorf("erro ao serializar sessão para o Redis: %w", err)
	}
	ctx := context.Background()
	err = cache.Set(ctx, key, data, 0).Err()
	if err != nil {
		return session.KeyId, fmt.Errorf("erro ao salvar sessão no Redis: %w", err)
	}

	return session.KeyId, nil
}

// Delete remove uma sessão do Redis.
func (r RedisRepositoryImpl) Delete(keyId string) error {
	cache := database.CACHE
	key := keyId
	ctx := context.Background()
	err := cache.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("erro ao deletar sessão do Redis: %w", err)
	}
	return nil
}
