package message

// List
type TodoList struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CreateList struct {
	Status   string   `json:"status"`
	UserId   int      `json:"user_id"`
	TodoList TodoList `json:"item"`
}

type UpdateList struct {
	Status   string   `json:"status"`
	UserId   int      `json:"user_id"`
	ListId   int      `json:"list_id"`
	TodoList TodoList `json:"item"`
}

type DeleteList struct {
	Status string `json:"status"`
	UserId int    `json:"user_id"`
	ListId int    `json:"list_id"`
}
