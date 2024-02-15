package models

import (
	"time"
)

//go:generate easyjson -all models.go
type StorageRecord struct {
	ID        int64     `json:"id"`
	Data      []byte    `json:"data"`
	Deleted   bool      `json:"deleted"`
	UpdatedAt time.Time `json:"updated_at"`
}
