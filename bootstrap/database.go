package bootstrap

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
	Mu   sync.RWMutex
}

func NewPostgresDb(env *Env) *PostgresClient {
	pool, err := pgxpool.New(context.Background(), env.PostgresConn)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	if err := pool.Ping(context.Background()); err != nil {
		logrus.Fatalf(err.Error())
	}
	return &PostgresClient{
		Pool: pool,
		Mu:   sync.RWMutex{},
	}
}
