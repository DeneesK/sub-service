package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type MonthYear struct {
	time.Time
}

const dataLayout = "01-2006" // MM-YYYY

func (my *MonthYear) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse(dataLayout, s)
	if err != nil {
		return err
	}
	my.Time = t
	return nil
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	return []byte(`"` + my.Time.Format(dataLayout) + `"`), nil
}

func (my MonthYear) Value() (driver.Value, error) {
	return my.Time, nil
}

func (my *MonthYear) Scan(value interface{}) error {
	if value == nil {
		*my = MonthYear{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot convert %T to time.Time", value)
	}
	my.Time = t
	return nil
}

// Subscription swagger:model
type Subscription struct {
	ID          string     `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	Price       int        `db:"price" json:"price"`
	UserID      string     `db:"user_id" json:"user_id"`
	StartDate   MonthYear  `db:"start_date" json:"start_date" swaggertype:"string"`
	EndDate     *MonthYear `db:"end_date" json:"end_date,omitempty" swaggertype:"string"`
}

// UpdateSubscription swagger:model
type UpdateSubscription struct {
	ServiceName *string    `db:"service_name" json:"service_name,omitempty"`
	Price       *int       `db:"price" json:"price,omitempty"`
	UserID      *string    `db:"user_id" json:"user_id,omitempty"`
	StartDate   *MonthYear `db:"start_date" json:"start_date,omitempty" swaggertype:"string"`
	EndDate     *MonthYear `db:"end_date" json:"end_date,omitempty" swaggertype:"string"`
}
