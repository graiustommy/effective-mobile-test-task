package repository

import (
	"context"

	"effective-task/internal/entity"
	apperrors "effective-task/internal/errors"

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
		return nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to create database pool", err)
	}
	return &Repository{
		pool: pool,
	}, nil
}

func (r *Repository) Exist(ctx context.Context, Request *entity.Subscription) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM subscriptions WHERE service_name = $1 AND price = $2 AND user_id = $3 AND start_date = $4) 
	`
	var existence bool
	err := r.pool.QueryRow(ctx, query, Request.ServiceName, Request.Price, Request.UserID, Request.StartDate).Scan(&existence)
	if err != nil {
		return false, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to check subscription existence", err)
	}
	return existence, nil
}

func (r *Repository) Create(ctx context.Context, Request *entity.Subscription) (uuid.UUID, error) {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		Request.ServiceName,
		Request.Price,
		Request.UserID,
		Request.StartDate,
		Request.EndDate).Scan(&id)
	if err != nil {
		return uuid.Nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to create subscription", err)
	}
	return id, nil
}

func (r *Repository) Update(ctx context.Context, Request *entity.Subscription) error {
	query := `
	UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6
	`
	cmd, err := r.pool.Exec(ctx, query, Request.ServiceName, Request.Price, Request.UserID, Request.StartDate, Request.EndDate, Request.ID)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrKindDatabase, "failed to update subscription", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrKindNotFound, "subscription not found")
	}
	return nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Subscription, error) {
	query := `
	SELECT * FROM subscriptions WHERE user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID.String())
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to list subscriptions", err)
	}
	subs := []*entity.Subscription{}
	for rows.Next() {
		var sub entity.Subscription
		err = rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to scan subscription", err)
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
		return apperrors.Wrap(apperrors.ErrKindDatabase, "failed to delete subscription", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrKindNotFound, "subscription not found")
	}
	return nil
}

func (r *Repository) ReadBySubscriptionID(ctx context.Context, SubscriptionID uuid.UUID) (*entity.Subscription, error) {
	query := `
	SELECT * FROM subscriptions WHERE id = $1
	`
	row, err := r.pool.Query(ctx, query, SubscriptionID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to read subscription", err)
	}
	sub := &entity.Subscription{}
	if !row.Next() {
		return nil, apperrors.New(apperrors.ErrKindNotFound, "subscription not found")
	}
	err = row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to scan subscription", err)
	}
	return sub, nil
}

func (r *Repository) CountByUserID(ctx context.Context, userID uuid.UUID, start, end string) (uint64, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE user_id = $1 AND  
		start_date >= $2
    	AND (end_date <= $3 OR end_date IS NULL)
	`
	var sum uint64
	err := r.pool.QueryRow(ctx, query, userID.String(), start, end).Scan(&sum)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to count subscriptions by user id", err)
	}
	return sum, nil
}

func (r *Repository) CountByServiceName(ctx context.Context, serviceName string, start, end string) (uint64, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE service_name = $1 AND  
		start_date >= $2
    	AND (end_date <= $3 OR end_date IS NULL)
	`
	var sum uint64
	err := r.pool.QueryRow(ctx, query, serviceName, start, end).Scan(&sum)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.ErrKindDatabase, "failed to count subscriptions by service name", err)
	}
	return sum, nil
}
