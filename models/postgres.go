package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "lenslocked",
		Password: "lenslocked",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

// migrate usado para executar migrations que foram embutidas no binário com embed
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	// In case the dir is an empty string, they probably meant the current directory and goose wants a period for that.
	if dir == "" {
		dir = "."
	}
	// define que o diretório base do fs dos migrations é o diretório
	// onde colocamos o fs.go dos migrations = migrations/fs.go
	// assim ele vai considerar esse caminho quando for executar as
	// migrations passando o local "."
	goose.SetBaseFS(migrationsFS)
	defer func() {
		// Ensure that we remove the FS on the off chance some other part of our app uses goose for migrations and doesn't want to use our FS.
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
