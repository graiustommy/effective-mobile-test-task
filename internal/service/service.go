package service

import (
	"context"
	"fmt"

	"effective-task/internal/entity"
	"effective-task/internal/repository"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(context.Context, *entity.Subscription) error
	Read(context.Context, string) *entity.Subscription
	Update(context.Context, *entity.Subscription) error
	Delete(context.Context, string)
	List(context.Context)
	Count(context.Context, string, string)
}

type Service struct {
	repo repository.RepositoryInterface
}

func (s *Service) Create(ctx context.Context, CreateRequest *entity.Subscription) error {
	if CreateRequest.UserID == uuid.Nil {
		return fmt.Errorf("userID required")
	}
	if CreateRequest.ServiceName == "" {
		return fmt.Errorf("serviceName is required")
	}

	return nil
}

func (s *Service) Update(ctx context.Context, Request *entity.Subscription) error {
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
