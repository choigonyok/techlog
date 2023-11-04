package model

type Cookie struct {
	Key   string
	Value string
}

type User struct {
	ID       string `json:"id"`
	Password string `json:"pw"`
}
