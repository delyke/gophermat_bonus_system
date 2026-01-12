package model

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Balance  float64   `json:"balance"`
}
