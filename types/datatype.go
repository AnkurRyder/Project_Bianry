package types

import "github.com/google/uuid"

// Data struct for storing data to the DBß
type Data struct {
	ID    uuid.UUID
	Value bool   `json:"value"`
	Key   string `json:"key"`
}

// User struct for storing user credentials to DB
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
