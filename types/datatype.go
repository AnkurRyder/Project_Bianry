package types

import (
	"github.com/google/uuid"
)

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

// TokenDetails to store token data
type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

// TokenMeta to store data for cross verification
type TokenMeta struct {
	ID      uint64
	AccesID string
	UserID  uint64
	ExpTime int64
}

// AccessDetails to cross check token
type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}
