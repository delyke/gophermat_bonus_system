package model

import (
	"time"

	"github.com/google/uuid"
)

type Withdrawal struct {
	UUID        uuid.UUID
	UserUUID    uuid.UUID
	OrderID     string
	Amount      float64
	ProcessedAt time.Time
}
