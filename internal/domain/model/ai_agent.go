package model

import (
	"fmt"
)

type AIAgent struct {
	Id                 uint64 `json:"id"`
	Name               string `json:"nome"`
	ModelId            uint   `json:"modelo"`
	CompanyName        string `json:"nomeEmpresa"`
	CompanyDescription string `json:"descricaoEmpresa"`
	BehaviourIa        string `json:"comportamentoIa"`
	CompanyUrl         string `json:"urlInformacoes"`
	Instructions       string
	UUID_user          string
}

// SetInstructions should use a pointer receiver to modify the original struct.
func (bot *AIAgent) SetInstructions(instructions string) {
	bot.Instructions = instructions
}

// GetSystemInstructions can use a value receiver, as it only reads the data.
func (bot AIAgent) GetSystemInstructions() string {
	return fmt.Sprintf(
		`
            Seu nome de agente: %s
            Nome de quem você representa: %s
            Se comporte de maneira: %s
            Esse é o site que você terá informações sobre quem voce tá representando: %s,
            Essa é uma breve descrição de quem você representa: %s

            Essas são suas instruções de personalidade: %s
        `, bot.Name,
		bot.CompanyName,
		bot.BehaviourIa,
		bot.CompanyUrl,
		bot.CompanyDescription,
		bot.Instructions,
	)
}
