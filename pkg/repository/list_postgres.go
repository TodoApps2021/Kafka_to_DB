package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoListPostgres struct {
	pool *pgxpool.Pool
}

func NewTodoListPostgres(pool *pgxpool.Pool) *TodoListPostgres {
	return &TodoListPostgres{pool: pool}
}

func (t *TodoListPostgres) Create(userId int, list message.TodoList) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(context.Background(), createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			return e
		}
		return err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(context.Background(), createUsersListQuery, userId, id)
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			return e
		}
		return err
	}

	return tx.Commit(context.Background())
}

func (t *TodoListPostgres) Delete(userId, listId int) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)

	_, err = conn.Exec(context.Background(), query, userId, listId)

	return err
}

func (t *TodoListPostgres) Update(userId, listId int, input message.TodoList) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	_, err = conn.Exec(context.Background(), query, args...)
	return err
}
