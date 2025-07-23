package model

import "time"

// Subscription swagger:model
type Subscription struct {
	ID          string    `db:"id" json:"id"`
	ServiceName string    `db:"service_name" json:"service_name"`
	Price       int       `db:"price" json:"price"`
	UserID      string    `db:"user_id" json:"user_id"`
	StartDate   time.Time `db:"start_date" json:"start_date"` // YYYY-MM-DD
	EndDate     time.Time `db:"end_date,omitempty" json:"end_date,omitempty"`
}
