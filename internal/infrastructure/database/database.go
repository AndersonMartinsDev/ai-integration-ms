package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

func conectar() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, erro := sql.Open("postgres", psqlconn)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}
	return db, nil
}

// conectarAdminDB conecta ao PostgreSQL usando um banco de dados padrão (como 'postgres')
func conectarAdminDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD) // Conectando a 'postgres' ou 'template1'
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o DB admin: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao pingar DB admin: %w", err)
	}
	return db, nil
}

func GetConnectionDatabase() *sql.DB {
	banco, erro := conectar()

	if erro != nil {
		panic("Conexão Não estabelecida com o Banco de dados!")
	}
	return banco
}

func InitalStrucuture() {
	banco, erro := conectarAdminDB()
	if erro != nil {
		panic("Problema ao criar tabelas")
	}
	defer banco.Close()
	/**
	ISSO Não funciona porque o banco tenta estabelecer uma conexão só que o banco não existe
	deve-se estabelecer uma conexão com o servico do banco e depois tentar encontrar o banco de dados em si
	*/
	_, err := banco.Exec(fmt.Sprintf("CREATE DATABASE %s", DB_NAME))
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "42P04" {
			slog.Info(fmt.Sprintf("Banco de dados '%s' já existe, prosseguindo.", DB_NAME))
		} else {
			panic(fmt.Sprintf("Erro ao criar banco de dados %s :%v", DB_NAME, err))
		}
	} else {
		slog.Info(fmt.Sprintf("Banco de dados '%s' criado com sucesso.", DB_NAME))
	}
	db, err := conectar() // Use a nova função conectarDBPrincipal
	if err != nil {
		panic(fmt.Sprintf("Problema ao conectar com o banco de dados principal '%s': %v", DB_NAME, err))
	}
	defer db.Close() // Garante que a conexão principal seja fechada

	// 4. Configura o driver do migrate com a conexão principal
	config := &postgres.Config{}
	driver, err := postgres.WithInstance(db, config) // Use 'db' aqui, não 'banco' (que foi renomeado)
	if err != nil {
		slog.Error("Erro ao configurar driver de migração: ", err)
	}

	// Criar uma Variável de ambiente para facilitar esse mapeamento ou mudar isso no dockerfile'
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Use o caminho que funciona para você após o COPY do Dockerfile
		DB_NAME,             // Use DB_NAME aqui
		driver,
	)
	if err != nil {
		slog.Error("Erro ao criar instância de migração: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Erro ao aplicar migrações: ", err)
	}

	slog.Info("Migrações aplicadas com sucesso!")

}
