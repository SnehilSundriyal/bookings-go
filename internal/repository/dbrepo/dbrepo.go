package dbrepo

import (
	"database/sql"
	"github.com/SnehilSundriyal/bookings-go/internal/config"
	"github.com/SnehilSundriyal/bookings-go/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
