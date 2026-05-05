package entity

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

type CountPriceByUserID struct {
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
type CountPriceByServiceName struct {
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
