package configuration

import (
	"ai-integration-ms/internal/infrastructure/ai/gemini"
	"ai-integration-ms/internal/infrastructure/commons/logger"
	"ai-integration-ms/internal/infrastructure/database"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	Porta  = 8080
	Origin = ""
)

func LoadEnv() {
	if erro := godotenv.Load(); erro != nil {
		panic("Error ao carregar as variáveis de ambiente!")
	}
	Origin = os.Getenv("ORIGINS")
	slog.Info("Variáveis de ambiente carregadas com sucesso!")
}

// LoadLogger apenas para carregar logs personalizados
func LoadLogger() {
	custom_log := slog.New(logger.NewHandler(nil))
	slog.SetDefault(custom_log)
	slog.Info("Logger Carregado com sucesso!")
}

func LoadServer(routers http.Handler) {
	slog.Info(fmt.Sprintf("Servidor iniciado na porta %d", Porta))
	if erro := http.ListenAndServe(fmt.Sprintf(":%d", Porta), routers); erro != nil {
		panic(fmt.Sprintf("Error ao iniciar servidor %s", erro.Error()))
	}
}

func LoadRedis() {
	redis_pass_secret, err := GetSecret("redis_password")
	if err != nil {
		panic(fmt.Sprintf("Error ao tentar recuperar Gemini Key: %e", err))
	}
	redis_address := os.Getenv("REDIS_ADDR")
	redis_port := os.Getenv("REDIS_PORT")
	redis_pass := redis_pass_secret

	database.SetRedisEnv(redis_address, redis_port, redis_pass)
	database.InitializeCache()
}

func ConfigGenerativeGemini() {
	gemini_key_secret, err := GetSecret("gemini_api_key")
	if err != nil {
		panic(fmt.Sprintf("Error ao tentar recuperar Gemini Key: %e", err))
	}
	gemini_config := gemini.GeminiConfig{}
	gemini_config.Url = os.Getenv("GEMINI_URL")
	gemini_config.Version = os.Getenv("GEMINI_VERSION")
	gemini_config.Model = os.Getenv("GEMINI_MODEL")
	gemini_config.Key = gemini_key_secret
	gemini_config.SetConfig()
}
