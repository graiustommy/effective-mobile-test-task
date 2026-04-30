package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       uint
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
}
