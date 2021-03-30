package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TodoApps2021/Kafka_to_DB/pkg/message"
	"github.com/TodoApps2021/Kafka_to_DB/pkg/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func (h *Handler) Handle(ctx context.Context, key, value []byte, topic string, partition int32) error {
	log.Printf("key: <%s>, value: <%s>", string(key), string(value))

	// add handler
	switch topic {
	case "auth":
		var input message.CreateUser
		if err := json.Unmarshal(value, &input); err != nil {
			return err
		}
		return h.CreateUser(input.User)
	case "todo_list":
		if partition == 0 {
			var input message.CreateList

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.CreateList(input.UserId, input.TodoList)
		} else if partition == 1 {
			var input message.DeleteList

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.DeleteList(input.UserId, input.ListId)
		} else if partition == 2 {
			var input message.UpdateList

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.UpdateList(input.UserId, input.ListId, input.TodoList)
		}
	case "todo_item":
		if partition == 0 {
			var input message.CreateItem

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.CreateItem(input.ListId, input.TodoItem)
		} else if partition == 1 {
			var input message.DeleteItem

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.DeleteItem(input.UserId, input.ItemId)
		} else if partition == 2 {
			var input message.UpdateItem

			if err := json.Unmarshal(value, &input); err != nil {
				return err
			}

			return h.UpdateItem(input.UserId, input.ItemId, input.TodoItem)
		}
	}

	return nil
}
