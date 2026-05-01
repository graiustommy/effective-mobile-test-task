package service

import (
	"context"
	"fmt"

	"effective-task/internal/entity"
	"effective-task/internal/repository"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(context.Context, *entity.Subscription) (uuid.UUID, error)
	ReadBySubscriptionID(context.Context, string) (*entity.Subscription, error)
	Update(context.Context, *entity.Subscription) error
	Delete(context.Context, string) error
	ListByUserID(context.Context, string) ([]*entity.Subscription, error)
	CountByUserID(context.Context, string, string, string) (uint64, error)
}

type Service struct {
	repo repository.RepositoryInterface
}

func NewService(repo repository.RepositoryInterface) ServiceInterface {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, Request *entity.Subscription) (uuid.UUID, error) {
	_, err := uuid.Parse(Request.UserID.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid userID")
	}
	if Request.ServiceName == "" {
		return uuid.Nil, fmt.Errorf("serviceName is required")
	}
	id, err := s.repo.Create(ctx, Request)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) ReadBySubscriptionID(ctx context.Context, userID string) (*entity.Subscription, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read subscription: %v", err)
	}
	sub, err := s.repo.ReadBySubscriptionID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to read subscription by ID: %v", err)
	}
	return sub, nil
}

func (s *Service) Update(ctx context.Context, Request *entity.Subscription) error {
	_, err := uuid.Parse(Request.UserID.String())
	if err != nil {
		return fmt.Errorf("invalid userID")
	}
	exist, err := s.repo.Exist(ctx, Request)
	if err != nil {
		return fmt.Errorf("failed to update: %v", err)
	}
	if !exist {
		return fmt.Errorf("subscription doesn't exist")
	}
	err = s.repo.Update(ctx, Request)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, SubscriptionID string) error {
	subUUID, err := uuid.Parse(SubscriptionID)
	if err != nil {
		return fmt.Errorf("invalid SubscriptionID")
	}
	err = s.repo.DeleteBySubscriptionID(ctx, subUUID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %v", err)
	}
	return nil
}

func (s *Service) ListByUserID(ctx context.Context, userID string) ([]*entity.Subscription, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID")
	}
	subs, err := s.repo.ListByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list by user id: %v", err)
	}
	return subs, nil
}

func (s *Service) CountByUserID(ctx context.Context, userID string, start, end string) (uint64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count sum: %v", err)
	}
	sum, err := s.repo.CountByUserID(ctx, userUUID, start, end)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (s *Service) CountByServiceName(ctx context.Context, serviceName string, start, end string) (uint64, error) {
	if serviceName != "" {
		return 0, fmt.Errorf("failed to count sum: invalid service name")
	}
	sum, err := s.repo.CountByServiceName(ctx, serviceName, start, end)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
