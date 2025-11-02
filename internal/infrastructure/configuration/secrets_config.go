package configuration

import (
	"os"
	"strings"
)

const secretsPath = "/run/secrets/"

// GetSecret le o valor de um arquivo Secret montado pelo Docker
func GetSecret(secretName string) (string, error) {
	filePath := secretsPath + secretName

	// Verifica se o arquivo existe (ou fallback para env)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Se nao for um Secret montado, tenta ler como Variavel de Ambiente (Fallback)
		return os.Getenv(secretName), nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Remove quebras de linha ou espacos em branco que o Docker possa adicionar
	return strings.TrimSpace(string(content)), nil
}
