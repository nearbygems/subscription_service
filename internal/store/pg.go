package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nearbygems/subscription-service/internal/model"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Create(sub *model.Subscription) error
	Get(id uuid.UUID) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uuid.UUID) error
	List(limit, offset int, userID *uuid.UUID, serviceName *string) ([]model.Subscription, error)
	Summary(periodFrom, periodTo string, userID *uuid.UUID, serviceName *string) (int, error)
}

type PostgresStore struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewPostgresStore(db *sql.DB, log *logrus.Logger) *PostgresStore {
	return &PostgresStore{db: db, log: log}
}

func (s *PostgresStore) Create(sub *model.Subscription) error {
	_, err := s.db.Exec(`
		insert into subscriptions 
		  (id, service_name, price, user_id, start_date, end_date, created_at)
    	values 
    	  ($1, $2, $3, $4, $5, $6, now())
    	  `, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	return err
}

func (s *PostgresStore) Get(id uuid.UUID) (*model.Subscription, error) {
	row := s.db.QueryRow(`
		select id, 
       		service_name, 
       		price, 
       		user_id, 
       		start_date, 
       		end_date, 
       		created_at 
		from subscriptions 
		where id=$1
		`, id)
	sub := model.Subscription{}
	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *PostgresStore) Update(sub *model.Subscription) error {
	_, err := s.db.Exec(`
		update subscriptions 
		set service_name=$1, 
		    price=$2,
			user_id=$3, 
			start_date=$4, 
			end_date=$5 
		where id=$6
		`, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	return err
}

func (s *PostgresStore) Delete(id uuid.UUID) error {
	_, err := s.db.Exec(`
		delete from subscriptions where id=$1
		`, id)
	return err
}

func (s *PostgresStore) List(limit, offset int, userID *uuid.UUID,
	serviceName *string) ([]model.Subscription, error) {
	query := `
		select id, 
		       service_name, 
		       price, 
		       user_id, 
		       start_date, 
		       end_date,
		       created_at 
		from subscriptions
		`
	clauses := []string{}
	args := []interface{}{}
	if userID != nil {
		clauses = append(clauses, fmt.Sprintf("user_id=$%d", len(args)+1))
		args = append(args, *userID)
	}
	if serviceName != nil {
		clauses = append(clauses, fmt.Sprintf("service_name=$%d", len(args)+1))
		args = append(args, *serviceName)
	}
	if len(clauses) > 0 {
		query += " WHERE " + sqlx.Rebind(sqlx.DOLLAR, clauses[0])
		for _, c := range clauses[1:] {
			query += " AND " + c
		}
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d",
		limit, offset)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []model.Subscription
	for rows.Next() {
		s := model.Subscription{}
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID,
			&s.StartDate, &s.EndDate, &s.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (s *PostgresStore) Summary(periodFrom, periodTo string, userID *uuid.UUID, serviceName *string) (int, error) {
	query := `
		select coalesce(sum(price),0) 
		from subscriptions 
		where start_date >= $1 and (end_date <= $2 or end_date is null)
		`
	args := []interface{}{periodFrom, periodTo}
	if userID != nil {
		query += fmt.Sprintf(" AND user_id=$%d", len(args)+1)
		args = append(args, *userID)
	}
	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name=$%d", len(args)+1)
		args = append(args, *serviceName)
	}
	row := s.db.QueryRow(query, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
