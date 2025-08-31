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

	totalPrice := 0
	for _, sub := range subs {
		// Определяем фактический период подписки в рамках выбранного диапазона
		actualStart := maxDate(*filter.StartDate, sub.StartDate)
		actualEnd := minDate(*filter.EndDate, sub.EndDate)

		monthsInPeriod := calculateMonthsInPeriod(actualStart, actualEnd)
		if monthsInPeriod > 0 {
			totalPrice += sub.Price * monthsInPeriod
		}
	}
	return totalPrice, nil
}

func calculateMonthsInPeriod(start, end string) int {
	startTime, _ := time.Parse(dateForm, start)
	endTime, _ := time.Parse(dateForm, end)

	yearsDiff := endTime.Year() - startTime.Year()
	monthsDiff := int(endTime.Month()) - int(startTime.Month())

	return yearsDiff*12 + monthsDiff + 1 // +1, чтобы включить текущий месяц
}

func maxDate(date1, date2 string) string {
	if date1 > date2 {
		return date1
	}
	return date2
}

func minDate(date1 string, date2 *string) string {
	if date2 == nil || date1 < *date2 {
		return date1
	}
	return *date2
}
