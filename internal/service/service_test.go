package service

import (
	"context"
	"errors"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"reflect"
	"testing"
	"time"
)

type mockRepo struct {
	searchFunc                    func(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error)
	createFunc                    func(ctx context.Context, input *domain.Subscription) error
	updateFunc                    func(ctx context.Context, input *domain.Subscription) error
	deleteFunc                    func(ctx context.Context, filter *domain.Filter) error
	getSubscriptionsForPeriodFunc func(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error)
	closeDBFunc                   func()
}

func (m *mockRepo) Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, filter)
	}
	return nil, nil
}
func (m *mockRepo) Create(ctx context.Context, input *domain.Subscription) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil
}
func (m *mockRepo) Update(ctx context.Context, input *domain.Subscription) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, input)
	}
	return nil
}
func (m *mockRepo) Delete(ctx context.Context, filter *domain.Filter) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, filter)
	}
	return nil
}
func (m *mockRepo) GetSubscriptionsForPeriod(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
	if m.getSubscriptionsForPeriodFunc != nil {
		return m.getSubscriptionsForPeriodFunc(ctx, filter)
	}
	return nil, nil
}
func (m *mockRepo) CloseDB() {
	if m.closeDBFunc != nil {
		m.closeDBFunc()
	}
}

func TestSubServiceImpl_Search(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name       string
		filter     *domain.Filter
		mockResult []*domain.Subscription
		mockErr    error
		want       []*domain.Subscription
		wantErr    bool
	}{
		{
			name: "success",
			filter: func() *domain.Filter {
				id := validUUID
				return &domain.Filter{UserID: &id}
			}(),
			mockResult: []*domain.Subscription{{UserID: validUUID, ServiceName: "Netflix"}},
			mockErr:    nil,
			want:       []*domain.Subscription{{UserID: validUUID, ServiceName: "Netflix"}},
			wantErr:    false,
		},
		{
			name: "no results",
			filter: func() *domain.Filter {
				id := validUUID
				return &domain.Filter{UserID: &id}
			}(),
			mockResult: []*domain.Subscription{},
			mockErr:    nil,
			want:       nil,
			wantErr:    false,
		},
		{
			name: "repo error",
			filter: func() *domain.Filter {
				id := validUUID
				return &domain.Filter{UserID: &id}
			}(),
			mockResult: nil,
			mockErr:    errors.New("db error"),
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				searchFunc: func(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
					return tt.mockResult, tt.mockErr
				},
			}
			service := NewService(repo)
			got, err := service.Search(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubServiceImpl_CreateSubscription(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	validService := "Netflix"
	validPrice := 100

	tests := []struct {
		name    string
		input   *domain.Subscription
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			input: &domain.Subscription{
				UserID:      validUUID,
				ServiceName: validService,
				Price:       validPrice,
				StartDate:   time.Now(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "repo error",
			input: &domain.Subscription{
				UserID:      validUUID,
				ServiceName: validService,
				Price:       validPrice,
				StartDate:   time.Now(),
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				createFunc: func(ctx context.Context, input *domain.Subscription) error {
					return tt.mockErr
				},
			}
			service := NewService(repo)
			err := service.CreateSubscription(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubServiceImpl_UpdateSubscription(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	validService := "Netflix"
	validPrice := 150

	tests := []struct {
		name    string
		input   *domain.Subscription
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			input: &domain.Subscription{
				UserID:      validUUID,
				ServiceName: validService,
				Price:       validPrice,
				StartDate:   time.Now(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "repo error",
			input: &domain.Subscription{
				UserID:      validUUID,
				ServiceName: validService,
				Price:       validPrice,
				StartDate:   time.Now(),
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				updateFunc: func(ctx context.Context, input *domain.Subscription) error {
					return tt.mockErr
				},
			}
			service := NewService(repo)
			err := service.UpdateSubscription(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubServiceImpl_DeleteSubscription(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	validService := "Netflix"

	tests := []struct {
		name    string
		filter  *domain.Filter
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			filter: &domain.Filter{
				UserID:      &validUUID,
				ServiceName: &validService,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "repo error",
			filter: &domain.Filter{
				UserID:      &validUUID,
				ServiceName: &validService,
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				deleteFunc: func(ctx context.Context, filter *domain.Filter) error {
					return tt.mockErr
				},
			}
			service := NewService(repo)
			err := service.DeleteSubscription(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubServiceImpl_GetSubscriptionsForPeriod(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	validService := "Netflix"
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		filter     *domain.Filter
		mockResult []*domain.Subscription
		mockErr    error
		want       []*domain.Subscription
		wantErr    bool
	}{
		{
			name: "success",
			filter: &domain.Filter{
				UserID:      &validUUID,
				ServiceName: &validService,
				StartDate:   &start,
				EndDate:     &end,
			},
			mockResult: []*domain.Subscription{
				{
					UserID:      validUUID,
					ServiceName: validService,
					Price:       100,
					StartDate:   start,
					EndDate:     &end,
				},
			},
			mockErr: nil,
			want: []*domain.Subscription{
				{
					UserID:      validUUID,
					ServiceName: validService,
					Price:       100,
					StartDate:   start,
					EndDate:     &end,
				},
			},
			wantErr: false,
		},
		{
			name: "repo error",
			filter: &domain.Filter{
				UserID:      &validUUID,
				ServiceName: &validService,
				StartDate:   &start,
				EndDate:     &end,
			},
			mockResult: nil,
			mockErr:    errors.New("db error"),
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				getSubscriptionsForPeriodFunc: func(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error) {
					return tt.mockResult, tt.mockErr
				},
			}
			got, err := repo.GetSubscriptionsForPeriod(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscriptionsForPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSubscriptionsForPeriod() got = %v, want %v", got, tt.want)
			}
		})
	}
}
