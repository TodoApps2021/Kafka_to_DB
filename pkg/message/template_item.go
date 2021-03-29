package message

// Item
type TodoItem struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type CreateItem struct {
	Status   string   `json:"status"`
	ListId   int      `json:"list_id"`
	TodoItem TodoItem `json:"item"`
}

type UpdateItem struct {
	Status   string   `json:"status"`
	UserId   int      `json:"user_id"`
	ItemId   int      `json:"item_id"`
	TodoItem TodoItem `json:"item"`
}

type DeleteItem struct {
	Status string `json:"status"`
	UserId int    `json:"user_id"`
	ItemId int    `json:"item_id"`
}
