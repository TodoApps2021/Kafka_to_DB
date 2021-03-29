package handler

import (
	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
)

func (s *Handler) CreateList(userId int, list message.TodoList) error {
	return s.Repo.TodoList.Create(userId, list) // to postgres sql
}

func (s *Handler) DeleteList(userId, listId int) error {
	return s.Repo.TodoList.Delete(userId, listId) // to postgres sql
}

func (s *Handler) UpdateList(userId, listId int, input message.TodoList) error {
	return s.Repo.TodoList.Update(userId, listId, input) // to postgres sql
}
