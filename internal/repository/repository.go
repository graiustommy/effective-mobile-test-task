package repository

import (
	"context"
	"fmt"

	"effective-task/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	Exist(context.Context, *entity.Subscription) (bool, error)
	Create(context.Context, *entity.Subscription) (uuid.UUID, error)
	Update(context.Context, *entity.Subscription) error
	ListByUserID(context.Context, uuid.UUID) ([]*entity.Subscription, error)
	DeleteBySubscriptionID(context.Context, uuid.UUID) error
	ReadBySubscriptionID(context.Context, uuid.UUID) (*entity.Subscription, error)
	CountByUserID(context.Context, uuid.UUID, string, string) (uint64, error)
	CountByServiceName(context.Context, string, string, string) (uint64, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(connStr string) (RepositoryInterface, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create new repository: %v", err)
	}
	return &Repository{
		pool: pool,
	}, nil
}

func (r *Repository) Exist(ctx context.Context, Request *entity.Subscription) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM subscriptions WHERE service_name = $1 AND price = $2 AND user_id = $3 AND start_date = $4) 
	`
	var existince bool
	err := r.pool.QueryRow(ctx, query, Request.ServiceName, Request.Price, Request.UserID, Request.StartDate).Scan(&existince)
	if err != nil {
		return false, fmt.Errorf("failed to check existince: %v", err)
	}
	return existince, nil
}

func (r *Repository) Create(ctx context.Context, Request *entity.Subscription) (uuid.UUID, error) {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		Request.ServiceName,
		Request.Price,
		Request.UserID,
		Request.StartDate,
		Request.EndDate).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create subscription: %v", err)
	}
	return id, nil
}

func (r *Repository) Update(ctx context.Context, Request *entity.Subscription) error {
	query := `
	UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $5
	`
	cmd, err := r.pool.Exec(ctx, query, Request.ServiceName, Request.Price, Request.UserID, Request.StartDate, Request.EndDate)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %v", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("failed to update subscription: 0 rows affected")
	}
	return nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Subscription, error) {
	query := `
	SELECT * FROM subscriptions WHERE user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list: %v", err)
	}
	subs := []*entity.Subscription{}
	for rows.Next() {
		var sub entity.Subscription
		err = rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to list: %v", err)
		}
		subs = append(subs, &sub)
	}
	return subs, nil
}

func (r *Repository) DeleteBySubscriptionID(ctx context.Context, SubscriptionID uuid.UUID) error {
	query := `
	DELETE FROM subscriptions WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query, SubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %v", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete subscription: subscription doesn't exist")
	}
	return nil
}

func (r *Repository) ReadBySubscriptionID(ctx context.Context, SubscriptionID uuid.UUID) (*entity.Subscription, error) {
	query := `
	SELECT * FROM subscriptions WHERE id = $1
	`
	row, err := r.pool.Query(ctx, query, SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to Read By SubscriptionID: %v", err)
	}
	sub := &entity.Subscription{}
	err = row.Scan(&sub)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscription by ID: %v", err)
	}
	return sub, nil
}

func (r *Repository) CountByUserID(ctx context.Context, userID uuid.UUID, start, end string) (uint64, error) {
	query := ` 
		SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE user_id = $1 AND  
		TO_CHAR(start_date, 'MM-YYYY') >= $2
    	AND (TO_CHAR(end_date, 'MM-YYYY') <= $3 OR end_date IS NULL)
	`
	var sum uint64
	err := r.pool.QueryRow(ctx, query, userID.String(), start, end).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to count sum by user id: %v", err)
	}
	return sum, nil
}

func (r *Repository) CountByServiceName(ctx context.Context, serviceName string, start, end string) (uint64, error) {
	query := ` 
		SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE service_name = $1 AND  
		TO_CHAR(start_date, 'MM-YYYY') >= $2
    	AND (TO_CHAR(end_date, 'MM-YYYY') <= $3 OR end_date IS NULL)
	`
	var sum uint64
	err := r.pool.QueryRow(ctx, query, serviceName, start, end).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to count sum by user id: %v", err)
	}
	return sum, nil
}
