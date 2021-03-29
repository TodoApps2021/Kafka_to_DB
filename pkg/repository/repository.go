package repository

import (
	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Authorization interface {
	CreateUser(user message.User) error
}

type TodoList interface {
	Create(userId int, list message.TodoList) error
	Delete(userId, listId int) error
	Update(userId, listId int, input message.TodoList) error
}

type TodoItem interface {
	Create(listId int, item message.TodoItem) error
	Delete(userId, itemId int) error
	Update(userId, itemId int, input message.TodoItem) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(pool),
		TodoList:      NewTodoListPostgres(pool),
		TodoItem:      NewTodoItemPostgres(pool),
	}
}
