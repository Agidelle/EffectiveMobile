package service

import (
	"context"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"log/slog"
	"time"
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

func (s *SubServiceImpl) Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
	res, err := s.repo.Search(ctx, filter)
	if err != nil {
		slog.Error("Failed to search subscriptions", "error", err)
		return nil, err
	}
	if len(res) == 0 {
		slog.Info("No subscription found", "filter", filter)
		return nil, nil
	}
	slog.Info("Subscription found", "subscription", res)
	return res, nil
}

func (s *SubServiceImpl) CreateSubscription(ctx context.Context, input *domain.Subscription) error {
	err := s.repo.Create(ctx, input)
	if err != nil {
		slog.Error("Failed to create subscription", "error", err)
		return err
	}
	slog.Info("Subscription created successfully", "input", input)
	return nil
}

func (s *SubServiceImpl) UpdateSubscription(ctx context.Context, input *domain.Subscription) error {
	err := s.repo.Update(ctx, input)
	if err != nil {
		slog.Error("Failed to update subscription", "error", err)
		return err
	}
	slog.Info("Subscription updated successfully", "input", input)
	return nil
}

func (s *SubServiceImpl) DeleteSubscription(ctx context.Context, filter *domain.Filter) error {
	err := s.repo.Delete(ctx, filter)
	if err != nil {
		slog.Error("Failed to delete subscription", "error", err)
		return err
	}
	slog.Info("Subscription deleted successfully", "user_id", filter.UserID, "service_name", filter.ServiceName)
	return nil
}

func (s *SubServiceImpl) GetSubscriptionsSummary(ctx context.Context, filter *domain.Filter) (int, error) {
	subs, err := s.repo.GetSubscriptionsForPeriod(ctx, filter)
	if err != nil {
		return 0, err
	}

	filterStart, _ := time.Parse(dateForm, *filter.StartDate)
	filterEnd, _ := time.Parse(dateForm, *filter.EndDate)

	totalPrice := 0
	for _, sub := range subs {
		subStart, _ := time.Parse(dateForm, sub.StartDate)
		subEnd, _ := time.Parse(dateForm, sub.EndDate)

		actualStart := maxTime(filterStart, subStart)
		actualEnd := minTime(filterEnd, subEnd)

		monthsInPeriod := calculateMonthsInPeriodTime(actualStart, actualEnd)
		if monthsInPeriod > 0 {
			totalPrice += sub.Price * monthsInPeriod
		}
	}
	return totalPrice, nil
}

func calculateMonthsInPeriodTime(start, end time.Time) int {
	yearsDiff := end.Year() - start.Year()
	monthsDiff := int(end.Month()) - int(start.Month())
	return yearsDiff*12 + monthsDiff + 1
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}
