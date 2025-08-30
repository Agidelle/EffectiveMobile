package domain

import "context"

type Subscription struct {
	UserID      string  `json:"user_id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionInput struct {
	UserID      *string `json:"user_id"`
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type Filter struct {
	UserID      *string
	ServiceName *string
	Price       *int
	StartDate   *string
	EndDate     *string
	Limit       *int
	Offset      *int
}

type Repository interface {
	Search(ctx context.Context, filter *Filter) ([]*Subscription, error)
	Create(ctx context.Context, sub *Subscription) error
	Update(ctx context.Context, sub *Subscription) error
	Delete(ctx context.Context, filter *Filter) error
	GetSubscriptionsForPeriod(ctx context.Context, filter *Filter) ([]*Subscription, error)
	CloseDB()
}

type SubscriptionOption func(*Subscription)

func NewSubscription(opts ...SubscriptionOption) *Subscription {
	s := &Subscription{}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SubscriptionInput) SubscriptionToOptions() []SubscriptionOption {
	var opts []SubscriptionOption
	if s.UserID != nil {
		opts = append(opts, WithUserID(*s.UserID))
	}
	if s.ServiceName != nil {
		opts = append(opts, WithServiceName(*s.ServiceName))
	}
	if s.Price != nil {
		opts = append(opts, WithPrice(*s.Price))
	}
	if s.StartDate != nil {
		opts = append(opts, WithStartDate(*s.StartDate))
	}
	if s.EndDate != nil {
		opts = append(opts, WithEndDate(*s.EndDate))
	}
	return opts
}

func WithUserID(id string) SubscriptionOption {
	return func(s *Subscription) {
		s.UserID = id
	}
}

func WithServiceName(name string) SubscriptionOption {
	return func(s *Subscription) {
		s.ServiceName = name
	}
}

func WithPrice(price int) SubscriptionOption {
	return func(s *Subscription) {
		s.Price = price
	}
}

func WithStartDate(date string) SubscriptionOption {
	return func(s *Subscription) {
		s.StartDate = date
	}
}

func WithEndDate(date string) SubscriptionOption {
	return func(s *Subscription) {
		s.EndDate = &date
	}
}
