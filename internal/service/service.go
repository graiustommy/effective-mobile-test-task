package service

import (
	"context"

	"effective-task/internal/entity"
	apperrors "effective-task/internal/errors"
	"effective-task/internal/repository"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(context.Context, *entity.Subscription) (uuid.UUID, error)
	ReadBySubscriptionID(context.Context, string) (*entity.Subscription, error)
	Update(context.Context, *entity.Subscription) error
	Delete(context.Context, string) error
	ListByUserID(context.Context, string) ([]*entity.Subscription, error)
	CountByUserID(context.Context, *entity.CountPriceByUserID) (uint64, error)
	CountByServiceName(context.Context, *entity.CountPriceByServiceName) (uint64, error)
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
		return uuid.Nil, apperrors.Wrap(apperrors.ErrKindValidation, "invalid user id", err)
	}
	if Request.ServiceName == "" {
		return uuid.Nil, apperrors.New(apperrors.ErrKindValidation, "service name is required")
	}
	if err := apperrors.ValidateDateFormat(Request.StartDate); err != nil {
		return uuid.Nil, err
	}
	if err := apperrors.ValidateDateFormat(Request.EndDate); err != nil {
		return uuid.Nil, err
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
		return nil, apperrors.Wrap(apperrors.ErrKindValidation, "invalid subscription id", err)
	}
	sub, err := s.repo.ReadBySubscriptionID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) Update(ctx context.Context, Request *entity.Subscription) error {
	_, err := uuid.Parse(Request.UserID.String())
	if err != nil {
		return apperrors.Wrap(apperrors.ErrKindValidation, "invalid user id", err)
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
		return apperrors.Wrap(apperrors.ErrKindValidation, "invalid subscription id", err)
	}
	err = s.repo.DeleteBySubscriptionID(ctx, subUUID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ListByUserID(ctx context.Context, userID string) ([]*entity.Subscription, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ErrKindValidation, "invalid user id", err)
	}
	subs, err := s.repo.ListByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *Service) CountByUserID(ctx context.Context, req *entity.CountPriceByUserID) (uint64, error) {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.ErrKindValidation, "invalid user id", err)
	}
	if err := apperrors.ValidateDateFormat(req.StartDate); err != nil {
		return 0, err
	}
	if err := apperrors.ValidateDateFormat(req.EndDate); err != nil {
		return 0, err
	}
	sum, err := s.repo.CountByUserID(ctx, userUUID, req.StartDate, req.EndDate)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (s *Service) CountByServiceName(ctx context.Context, req *entity.CountPriceByServiceName) (uint64, error) {
	if req.ServiceName == "" {
		return 0, apperrors.New(apperrors.ErrKindValidation, "service name is required")
	}
	if err := apperrors.ValidateDateFormat(req.StartDate); err != nil {
		return 0, err
	}
	if err := apperrors.ValidateDateFormat(req.EndDate); err != nil {
		return 0, err
	}
	sum, err := s.repo.CountByServiceName(ctx, req.ServiceName, req.StartDate, req.EndDate)
	if err != nil {
		return 0, err
	}
	return sum, nil
}
