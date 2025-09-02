package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	dateForm   = "01-2006"
	tOutnormal = 3 * time.Second
	tOutlong   = 10 * time.Second
)

type Handler struct {
	service SubService
}

type SubService interface {
	CloseDB()
	Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error)
	CreateSubscription(ctx context.Context, input *domain.Subscription) error
	UpdateSubscription(ctx context.Context, input *domain.Subscription) error
	DeleteSubscription(ctx context.Context, filter *domain.Filter) error
	GetSubscriptionsSummary(ctx context.Context, filter *domain.Filter) (int, error)
}

func NewHandler(s SubService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes(r chi.Router) {
	// Роуты для управления подписками
	r.Get("/api/subscriptions", h.SearchSubscriptions)   // список подписок
	r.Post("/api/subscriptions", h.CreateSubscription)   // создание новой подписки
	r.Put("/api/subscriptions", h.UpdateSubscription)    // обновление подписки по ID
	r.Delete("/api/subscriptions", h.DeleteSubscription) // удаление подписки по ID

	r.Post("/api/subscriptions/summary", h.GetSubscriptionsSummary) // сводная информация по подпискам
}

func (h *Handler) SearchSubscriptions(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter

	userID := r.URL.Query().Get("user_id")
	if userID != "" || len(userID) == 36 {
		filter.UserID = &userID
	}
	serviceName := r.URL.Query().Get("service_name")
	if serviceName != "" || len(serviceName) <= 255 {
		filter.ServiceName = &serviceName
	}
	priceStr := r.URL.Query().Get("price")
	if priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			filter.Price = &price
		}
	}
	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr != "" {
		if t, err := time.Parse(dateForm, startDateStr); err == nil {
			filter.StartDate = &t
		} else {
			http.Error(w, "invalid start_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
	}
	endDateStr := r.URL.Query().Get("end_date")
	if endDateStr != "" {
		if t, err := time.Parse(dateForm, startDateStr); err == nil {
			filter.EndDate = &t
		} else {
			http.Error(w, "invalid start_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = &limit
		}
	}
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = &offset
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutlong)
	defer cancel()
	subs, err := h.service.Search(ctx, &filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Failed to encode to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input domain.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}
	err := validateSubscriptionInput(&input)
	if err != nil {
		slog.Error("Invalid subscription input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opts := input.SubscriptionToOptions()
	sub := domain.NewSubscription(opts...)

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	if err := h.service.CreateSubscription(ctx, sub); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var input domain.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}
	err := validateSubscriptionInput(&input)
	if err != nil {
		slog.Error("Invalid subscription input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opts := input.SubscriptionToOptions()
	sub := domain.NewSubscription(opts...)

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	if err := h.service.UpdateSubscription(ctx, sub); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter

	userID := r.URL.Query().Get("user_id")
	if userID != "" || len(userID) == 36 {
		filter.UserID = &userID
	}
	serviceName := r.URL.Query().Get("service_name")
	if serviceName != "" || len(serviceName) <= 255 {
		filter.ServiceName = &serviceName
	}

	if filter.UserID == nil || filter.ServiceName == nil {
		http.Error(w, "user_id and service_name are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	err := h.service.DeleteSubscription(ctx, &filter)
	if err != nil {
		if err.Error() == "subscription not found" {
			http.Error(w, "subscription not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to delete subscription", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSubscriptionsSummary(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "invalid filter", http.StatusBadRequest)
		return
	}

	if filter.StartDateStr != nil && *filter.StartDateStr != "" {
		if t, err := time.Parse(dateForm, *filter.StartDateStr); err == nil {
			filter.StartDate = &t
		} else {
			http.Error(w, "invalid start_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
	}
	if filter.EndDateStr != nil && *filter.EndDateStr != "" {
		if t, err := time.Parse(dateForm, *filter.EndDateStr); err == nil {
			filter.EndDate = &t
		} else {
			http.Error(w, "invalid end_date format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
	}
	if filter.StartDate == nil || filter.EndDate == nil {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutlong)
	defer cancel()
	total, err := h.service.GetSubscriptionsSummary(ctx, &filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	resp := map[string]int{"total_price": total}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func validateSubscriptionInput(input *domain.SubscriptionInput) error {
	if input.UserID == nil || *input.UserID == "" || len(*input.UserID) != 36 {
		return fmt.Errorf("user_id is required correct format UUID")
	}
	if input.ServiceName == nil || *input.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if input.ServiceName != nil && len(*input.ServiceName) > 255 {
		return fmt.Errorf("service_name must not exceed 255 characters")
	}
	if input.Price == nil || *input.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if input.StartDate == nil || *input.StartDate == "" {
		return fmt.Errorf("start_date is required")
	}
	if _, err := time.Parse(dateForm, *input.StartDate); err != nil {
		return fmt.Errorf("invalid start_date format, expected MM-YYYY")
	}
	if input.EndDate != nil && *input.EndDate != "" {
		if _, err := time.Parse(dateForm, *input.EndDate); err != nil {
			return fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}
	}
	return nil
}
