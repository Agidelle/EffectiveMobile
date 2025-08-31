package api

import (
	"context"
	"encoding/json"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	service SubService
}

type SubService interface {
	CloseDB()
	Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error)
	CreateSubscription(ctx context.Context, input *domain.Subscription) error
	UpdateSubscription(ctx context.Context, id int, input *domain.Subscription) error
	DeleteSubscription(ctx context.Context, filter *domain.Filter) error
	GetSubscriptionsSummary(ctx context.Context, filter *domain.Filter) (int, error)
}

func NewHandler(s SubService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes(r chi.Router) {
	// Роуты для управления подписками
	r.Get("/api/subscriptions", h.ListSubscriptions)          // список подписок
	r.Get("/api/subscriptions/{id}", h.GetSubscriptionByID)   // получение подписки по ID
	r.Post("/api/subscriptions", h.CreateSubscription)        // создание новой подписки
	r.Put("/api/subscriptions/{id}", h.UpdateSubscription)    // обновление подписки по ID
	r.Delete("/api/subscriptions/{id}", h.DeleteSubscription) // удаление подписки по ID

	r.Get("/api/subscriptions/summary", h.GetSubscriptionsSummary) // сводная информация по подпискам
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	var input domain.Filter
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	subs, err := h.service.Search(ctx, &input)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Failed to encode to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	filter := &domain.Filter{UserID: &userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	subs, err := h.service.Search(ctx, filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if len(subs) == 0 {
		http.Error(w, "subscriptions not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Failed to encode subscriptions to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input domain.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}

	opts := input.SubscriptionToOptions()
	sub := domain.NewSubscription(opts...)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.CreateSubscription(ctx, sub); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetSubscriptionsSummary(w http.ResponseWriter, r *http.Request) {
}
