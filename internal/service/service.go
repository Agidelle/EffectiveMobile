package service

import (
	"context"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
)

const dateForm string = "01-2006"

type SubServiceImpl struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *SubServiceImpl {
	return &SubServiceImpl{repo: repo}
}

func (s *SubServiceImpl) CloseDB() {
	s.repo.CloseDB()
}

func (s *SubServiceImpl) ListSubscriptions(ctx context.Context) ([]*domain.Subscription, error) {
	return nil, nil
}

func (s *SubServiceImpl) GetSubscriptionByID(ctx context.Context, id int) (*domain.Subscription, error) {
	return nil, nil
}

func (s *SubServiceImpl) CreateSubscription(ctx context.Context, input *domain.SubscriptionInput) (int64, error) {
	return 0, nil
}

func (s *SubServiceImpl) UpdateSubscription(ctx context.Context, id int, input *domain.SubscriptionInput) error {
	return nil
}

func (s *SubServiceImpl) DeleteSubscription(ctx context.Context, id int) error {
	return nil
}

func (s *SubServiceImpl) GetSubscriptionsSummary(ctx context.Context) (interface{}, error) {
	return nil, nil
}
