package repository

import (
	"context"
	"fmt"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthPostgres struct {
	pool *pgxpool.Pool
}

func NewAuthPostgres(pool *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{pool: pool}
}

func (ap *AuthPostgres) CreateUser(user message.User) error {
	ctx := context.Background()
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	conn, err := ap.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
