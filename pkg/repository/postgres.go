package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Config struct {
	DB_URL string
}

func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), cfg.DB_URL)
	if err != nil {
		return nil, err
	}
	logrus.Info("Connected to DB!")

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return pool, nil
}
