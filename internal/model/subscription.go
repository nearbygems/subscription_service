package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ServiceName string    `db:"service_name" json:"service_name"`
	Price       int       `db:"price" json:"price"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	StartDate   string    `db:"start_date" json:"start_date"`
	EndDate     *string   `db:"end_date" json:"end_date,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
