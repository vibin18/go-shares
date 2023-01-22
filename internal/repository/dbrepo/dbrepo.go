package dbrepo

import (
	"database/sql"
	"github.com/vibin18/go-shares/internal/config"
	"github.com/vibin18/go-shares/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, q *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: q,
		DB:  conn,
	}
}
