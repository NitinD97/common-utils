package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func Init(pgConfig PgConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
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
	config.ConnConfig.Tracer = &CustomTracer{logger: logger}

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

type CustomTracer struct {
	logger *zap.Logger
}

func (ct *CustomTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	ct.logger.Debug("Query started",
		zap.String("sql", data.SQL),
		zap.Any("args", data.Args),
	)
	return ctx
}

func (ct *CustomTracer) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	if data.Err != nil {
		ct.logger.Error("Query failed",
			zap.Error(data.Err),
		)
	} else {
		ct.logger.Debug("Query succeeded",
			zap.String("commandTag", data.CommandTag.String()),
		)
	}
}
