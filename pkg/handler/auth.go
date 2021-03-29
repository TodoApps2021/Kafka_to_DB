package handler

import (
	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
)

func (h *Handler) CreateUser(user message.User) error {
	return h.Repo.Authorization.CreateUser(user) // to postgres sql
}
