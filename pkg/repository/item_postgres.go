package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoItemPostgres struct {
	pool *pgxpool.Pool
}

func NewTodoItemPostgres(pool *pgxpool.Pool) *TodoItemPostgres {
	return &TodoItemPostgres{pool: pool}
}

func (t *TodoItemPostgres) Create(listId int, item message.TodoItem) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(context.Background(), createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			return err
		}
		return err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(context.Background(), createListItemsQuery, listId, itemId)
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			return err
		}
		return err
	}

	return tx.Commit(context.Background())
}

func (t *TodoItemPostgres) Delete(userId, itemId int) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err = conn.Exec(context.Background(), query, userId, itemId)

	return err
}

func (t *TodoItemPostgres) Update(userId, itemId int, input message.TodoItem) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err = conn.Exec(context.Background(), query, args...)
	return err
}
