package handler

import (
	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
)

func (s *Handler) CreateItem(listId int, item message.TodoItem) error {
	return s.Repo.TodoItem.Create(listId, item) // to postgres sql
}

func (s *Handler) DeleteItem(userId, itemId int) error {
	return s.Repo.TodoItem.Delete(userId, itemId) // to postgres sql
}

func (s *Handler) UpdateItem(userId, itemId int, input message.TodoItem) error {
	return s.Repo.TodoItem.Update(userId, itemId, input) // to postgres sql
}
