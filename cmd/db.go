package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/moyai-network/build/internal/app"
	"github.com/moyai-network/flyway/storage/cache"
	flywayports "github.com/moyai-network/flyway/storage/ports"
	"github.com/moyai-network/flyway/storage/postgres"
)

func loadStorage(conf app.Config) (flywayports.Store, *sql.DB, error) {
	db, err := newDB(context.Background(), conf.Database)
	if err != nil {
		return nil, nil, err
	}

	var store flywayports.Store = postgres.NewStore(db)
	store = cache.NewCachedStorage(context.Background(), store, conf.Database.DSN())

	return store, db, nil
}

func newDB(ctx context.Context, cfg app.DatabaseConfig) (*sql.DB, error) {
	if cfg.Host == "" || cfg.Name == "" || cfg.User == "" {
		return nil, fmt.Errorf("database config is incomplete: host/name/user are required")
	}

	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}

	if cfg.MaxConns > 0 {
		db.SetMaxOpenConns(int(cfg.MaxConns))
		db.SetMaxIdleConns(int(cfg.MaxConns))
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
