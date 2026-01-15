package model

import "github.com/google/uuid"

type OrderJob struct {
	UserID      uuid.UUID
	OrderUUID   uuid.UUID
	OrderNumber string
	OrderStatus OrderStatus
	Attempts    int
}
