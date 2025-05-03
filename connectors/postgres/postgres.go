package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"net/url"
)

var db *pgxpool.Pool

type PgConfig struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	User               string `json:"user"`
	Password           string `json:"password"`
	DbName             string `json:"db_name"`
	SslMode            string `json:"ssl_mode"`
	PoolMinConnections int    `json:"pool_min_conns"`
}

func Init(pgConfig PgConfig) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pgConfig.User,
		url.PathEscape(pgConfig.Password),
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.DbName,
		pgConfig.SslMode,
	)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, errors.Errorf("failed to parse database URL: %s, error: %w", dbURL, err)
	}

	config.MinConns = int32(pgConfig.PoolMinConnections)
	config.MinIdleConns = int32(pgConfig.PoolMinConnections)

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %s, error: %w", dbURL, err)
	}
	return conn, nil
}

func Disconnect() {
	if db != nil {
		db.Close()
	}
}
