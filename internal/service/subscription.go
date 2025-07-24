package service

import (
	"fmt"
	"strings"
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

func (s *SubscriptionService) List(userID string) ([]model.Subscription, error) {
	var subs []model.Subscription
	var err error

	if userID == "" {
		query := `SELECT * FROM subscriptions ORDER BY start_date DESC`
		err = s.db.Select(&subs, query)
	} else {
		query := `SELECT * FROM subscriptions WHERE user_id = $1 ORDER BY start_date DESC`
		err = s.db.Select(&subs, query, userID)
	}

	return subs, err
}

func (s *SubscriptionService) Update(id string, upd *model.UpdateSubscription) error {
	setClauses := []string{}
	args := map[string]interface{}{"id": id}

	if upd.ServiceName != nil {
		setClauses = append(setClauses, "service_name=:service_name")
		args["service_name"] = *upd.ServiceName
	}
	if upd.Price != nil {
		setClauses = append(setClauses, "price=:price")
		args["price"] = *upd.Price
	}
	if upd.UserID != nil {
		setClauses = append(setClauses, "user_id=:user_id")
		args["user_id"] = *upd.UserID
	}
	if upd.StartDate != nil {
		setClauses = append(setClauses, "start_date=:start_date")
		args["start_date"] = *upd.StartDate
	}
	if upd.EndDate != nil {
		setClauses = append(setClauses, "end_date=:end_date")
		args["end_date"] = *upd.EndDate
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf(`UPDATE subscriptions SET %s WHERE id=:id`, strings.Join(setClauses, ", "))
	_, err := s.db.NamedExec(query, args)
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
