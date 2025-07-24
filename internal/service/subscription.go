package service

import (
	"fmt"
	"time"

	"github.com/DeneesK/sub-service/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	db  *sqlx.DB
	log *zap.SugaredLogger
}

func NewSubscriptionService(db *sqlx.DB, log *zap.SugaredLogger) *SubscriptionService {
	return &SubscriptionService{db: db, log: log}
}

func (s *SubscriptionService) Create(sub *model.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return s.db.QueryRow(
		query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&sub.ID)
}

func (s *SubscriptionService) Get(id string) (*model.Subscription, error) {
	var sub model.Subscription
	err := s.db.Get(&sub, "SELECT * FROM subscriptions WHERE id=$1", id)
	return &sub, err
}

func (s *SubscriptionService) List() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := s.db.Select(&subs, "SELECT * FROM subscriptions ORDER BY start_date DESC")
	return subs, err
}

func (s *SubscriptionService) Update(id string, upd *model.Subscription) error {
	_, err := s.db.NamedExec(
		`UPDATE subscriptions SET
            service_name=:service_name,
            price=:price,
            user_id=:user_id,
            start_date=:start_date,
            end_date=:end_date
         WHERE id=:id`,
		map[string]interface{}{
			"id":           id,
			"service_name": upd.ServiceName,
			"price":        upd.Price,
			"user_id":      upd.UserID,
			"start_date":   upd.StartDate,
			"end_date":     upd.EndDate,
		},
	)
	return err
}

func (s *SubscriptionService) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM subscriptions WHERE id=$1", id)
	return err
}

func (s *SubscriptionService) Aggregate(
	from, to time.Time, userID, service string,
) (int, error) {
	q := `SELECT COALESCE(SUM(price),0) FROM subscriptions
          WHERE start_date >= $1 AND start_date <= $2`
	args := []interface{}{from, to}
	if userID != "" {
		q += " AND user_id = $" + fmt.Sprint(len(args)+1)
		args = append(args, userID)
	}
	if service != "" {
		q += " AND service_name = $" + fmt.Sprint(len(args)+1)
		args = append(args, service)
	}
	var sum int
	err := s.db.Get(&sum, q, args...)
	return sum, err
}
