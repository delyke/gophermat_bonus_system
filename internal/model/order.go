package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderNew        OrderStatus = "NEW"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderInvalid    OrderStatus = "INVALID"
	OrderProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	UUID       uuid.UUID   `json:"uuid"`
	OrderID    string      `json:"order_id"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual"`
	UserUUID   uuid.UUID   `json:"user_uuid"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
