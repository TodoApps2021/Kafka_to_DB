package message

// Auth
type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUser struct {
	Status string `json:"status"`
	User   User   `json:"item"`
}
