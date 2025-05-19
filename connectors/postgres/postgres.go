package postgres

import (
	"context"
	"fmt"
	"github.com/NitinD97/common-utils/enums"
	"github.com/NitinD97/common-utils/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse database URL: %s", dbURL))
	}

	config.MinConns = int32(pgConfig.PoolMinConnections)
	config.MinIdleConns = int32(pgConfig.PoolMinConnections)
	config.ConnConfig.Tracer = &CustomTracer{logger: logger}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to connect to database: %s", dbURL))
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
	dbQueryId := uuid.NewString()
	ctx = context.WithValue(ctx, enums.DbQueryId, dbQueryId)
	ct.logger.Info("Query started",
		zap.String("sql", data.SQL),
		zap.Any("args", data.Args),
		zap.Any(enums.RequestId, ctx.Value(enums.RequestId)),
		zap.String(enums.DbQueryId, dbQueryId),
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
			zap.Any(enums.RequestId, ctx.Value(enums.RequestId)),
			zap.String(enums.DbQueryId, ctx.Value(enums.DbQueryId).(string)),
		)
	} else {
		ct.logger.Info("Query succeeded",
			zap.String("commandTag", data.CommandTag.String()),
			zap.Any(enums.RequestId, ctx.Value(enums.RequestId)),
			zap.String(enums.DbQueryId, ctx.Value(enums.DbQueryId).(string)),
		)
	}
}
